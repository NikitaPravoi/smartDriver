package iiko

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"math/big"
	"net/http"
	"smartDriver/internal/db"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	OldRevision = errors.New("TOO_OLD_REVISION")
)

// IikoClient handles communication with the iiko API
type IikoClient struct {
	httpClient *http.Client
	baseURL    string
}

// OrderPollingService manages polling for orders from multiple organizations
type OrderPollingService struct {
	db         *pgxpool.Pool
	queries    *db.Queries // sqlc generated queries
	iikoClient *IikoClient
	//centrifuge   *gocent.Client
	pollInterval time.Duration
	tokenCache   sync.Map // Cache for organization tokens
	lastRevision sync.Map // Track last revision for each organization
}

type AuthResponse struct {
	CorrelationID string `json:"correlationId"`
	Token         string `json:"token"`
}

type OrdersResponse struct {
	CorrelationID         string `json:"correlationId"`
	MaxRevision           int64  `json:"maxRevision"`
	OrdersByOrganizations []struct {
		OrganizationID string  `json:"organizationId"`
		Orders         []Order `json:"orders"`
	} `json:"ordersByOrganizations"`
}

type ErrorResponse struct {
	CorrelationId    string `json:"correlationId"`
	ErrorDescription string `json:"errorDescription"`
	ErrorType        string `json:"error"`
}

func (e *ErrorResponse) Error() string {
	return e.ErrorDescription
}

type Customer struct {
	Name string `json:"name"`
}

type Address struct {
	Street struct {
		Name string `json:"name"`
		City struct {
			Name string `json:"name"`
		} `json:"city"`
	} `json:"street"`
	Index    string `json:"index"`
	House    string `json:"house"`
	Building string `json:"building"`
	Flat     string `json:"flat"`
	Entrance string `json:"entrance"`
	Floor    string `json:"floor"`
	Comment  string `json:"comment"`
}

type Order struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organizationId"`
	Info           struct {
		Status         string   `json:"status"`
		DeliveryStatus string   `json:"deliveryStatus"`
		CreatedAt      string   `json:"whenCreated"`
		CompleteBefore string   `json:"completeBefore"`
		Customer       Customer `json:"customer"`
		DeliveryPoint  struct {
			Coordinates struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			} `json:"coordinates"`
			Address `json:"address"`
		} `json:"deliveryPoint"`
		Sum float64 `json:"sum"`
	} `json:"order"`
}

// NewOrderPollingService creates a new polling service instance
func NewOrderPollingService(
	db *pgxpool.Pool,
	queries *db.Queries,
	centrifugeURL string,
	centrifugeAPIKey string,
	pollInterval time.Duration,
) *OrderPollingService {
	return &OrderPollingService{
		db:         db,
		queries:    queries,
		iikoClient: NewIikoClient("https://api-ru.iiko.services"),
		//centrifuge:   gocent.New(centrifugeURL, gocent.WithAPIKey(centrifugeAPIKey)),
		pollInterval: pollInterval,
	}
}

// Start begins polling for all organizations
func (s *OrderPollingService) Start(ctx context.Context) error {
	// Load all organizations on startup
	orgs, err := s.queries.ListOrganizations(ctx)
	if err != nil || len(orgs) == 0 {
		return fmt.Errorf("failed to list organizations: %w", err)
	}
	fmt.Println("List of organizations:", orgs)

	// Start polling for each organization
	for _, org := range orgs {
		go s.pollOrganization(ctx, org)
	}

	return nil
}

// pollOrganization handles polling for a single organization
func (s *OrderPollingService) pollOrganization(ctx context.Context, org db.Organization) {
	ticker := time.NewTicker(s.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := s.pollOrganizationOrders(ctx, org); err != nil {
				// Log error but continue polling
				fmt.Printf("Error polling organization %s: %v\n", org.Name, err)
			}
		}
	}
}

// pollOrganizationOrders fetches and processes orders for one organization
func (s *OrderPollingService) pollOrganizationOrders(ctx context.Context, org db.Organization) error {
	// Get or refresh auth token
	token, err := s.getAuthToken(ctx, org)
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}
	fmt.Println("Received auth token:", token)

	orgs, err := s.iikoClient.GetOrganizationIds(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to get organization ids: %w", err)
	}

	// Get last known revision
	var lastRevId int64
	lastRev, ok := s.lastRevision.Load(org.ID)
	if !ok {
		lastRevFromDB, err := s.queries.GetLastRevision(ctx, &org.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("failed to get last revision from db: %w", err)
		}

		if lastRevFromDB != 0 {
			lastRevId = lastRevFromDB
		} else {
			// If no revision in DB, fetch initial revision
			initialRev, err := s.iikoClient.GetInitialRevision(ctx, token, orgs)
			if err != nil {
				return fmt.Errorf("failed to get initial revision: %w", err)
			}

			// Save initial revision to DB
			if err := s.queries.CreateRevision(ctx, db.CreateRevisionParams{
				OrganizationID: &org.ID,
				RevisionID:     initialRev,
			}); err != nil {
				return fmt.Errorf("failed to save initial revision: %w", err)
			}

			lastRevId = initialRev
		}
		s.lastRevision.Store(org.ID, lastRevId)
	} else {
		lastRevId = lastRev.(int64)
	}
	fmt.Println("Last rev value:", lastRevId)

	// Fetch orders since last revision
GettingOrders:
	orders, maxRevision, err := s.iikoClient.GetOrdersByRevision(ctx, token, orgs, lastRevId)
	if err != nil {
		if errors.Is(err, OldRevision) {
			initialRev, err := s.iikoClient.GetInitialRevision(ctx, token, orgs)
			if err != nil {
				return fmt.Errorf("failed to get initial revision: %w", err)
			}
			lastRevId = initialRev
			s.lastRevision.Store(org.ID, initialRev)
			goto GettingOrders
		}
		return fmt.Errorf("failed to get orders: %w", err)
	}
	fmt.Printf("Received orders len=%d by last rev value: %d\n", len(orders), maxRevision)

	// Process orders and update database
	if err := s.processOrders(ctx, orders); err != nil {
		return fmt.Errorf("failed to process orders: %w", err)
	}
	fmt.Println("Processed orders:", orders)

	// Update last known revision in both memory and DB
	s.lastRevision.Store(org.ID, maxRevision)
	if err := s.queries.CreateRevision(ctx, db.CreateRevisionParams{
		OrganizationID: &org.ID,
		RevisionID:     maxRevision,
	}); err != nil {
		return fmt.Errorf("failed to save max revision: %w", err)
	}

	// Publish updates to Centrifugo
	if err := s.publishUpdates(orders); err != nil {
		return fmt.Errorf("failed to publish updates: %w", err)
	}

	return nil
}

// getAuthToken gets a cached token or requests a new one
func (s *OrderPollingService) getAuthToken(ctx context.Context, org db.Organization) (string, error) {
	// Check cache first
	if token, ok := s.tokenCache.Load(org.ID); ok {
		return token.(string), nil
	}

	// Request new token
	token, err := s.iikoClient.Authenticate(ctx, org.IikoApiToken)
	if err != nil {
		return "", err
	}

	// Cache token (consider adding expiration)
	s.tokenCache.Store(org.ID, token)
	return token, nil
}

// processOrders updates the database with new order information
func (s *OrderPollingService) processOrders(ctx context.Context, orders []Order) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := s.queries.WithTx(tx)

	for _, order := range orders {
		cost := pgtype.Numeric{
			Int:   big.NewInt(int64(order.Info.Sum)),
			Exp:   0,
			NaN:   false,
			Valid: true,
		}

		floor64, _ := strconv.ParseInt(order.Info.DeliveryPoint.Floor, 10, 32)
		floor := int32(floor64)

		entrance64, _ := strconv.ParseInt(order.Info.DeliveryPoint.Entrance, 10, 32)
		entrance := int32(entrance64)

		createdAt, _ := time.Parse("2006-01-02 15:04:05.000", order.Info.CreatedAt)

		params := db.CreateOrderParams{
			CustomerName: strings.TrimSpace(order.Info.Customer.Name),
			City:         &order.Info.DeliveryPoint.Address.Street.City.Name,
			Street:       &order.Info.DeliveryPoint.Address.Street.Name,
			Apartment:    &order.Info.DeliveryPoint.Address.Flat,
			Floor:        &floor,
			Entrance:     &entrance,
			Comment:      &order.Info.DeliveryPoint.Comment,
			Cost:         cost,
			Status:       &order.Info.Status,
			Point:        order.Info.DeliveryPoint.Coordinates.Longitude,
			Point_2:      order.Info.DeliveryPoint.Coordinates.Latitude,
			CreatedAt:    pgtype.Timestamp{Time: createdAt, Valid: true},
		}

		if _, err := qtx.CreateOrder(ctx, params); err != nil {
			return fmt.Errorf("failed to create order: %w", err)
		}
		fmt.Println("Created order:", order)
	}

	return tx.Commit(ctx)
}

// publishUpdates sends order updates to Centrifugo
func (s *OrderPollingService) publishUpdates(orders []Order) error {
	//for _, order := range orders {
	//	data, err := json.Marshal(order)
	//	if err != nil {
	//		return fmt.Errorf("failed to marshal order: %w", err)
	//	}
	//
	//	// Publish to organization-specific channel
	//	channel := fmt.Sprintf("orders:%s", order.OrganizationID)
	//	if err := s.centrifuge.Publish(context.Background(), channel, data); err != nil {
	//		return fmt.Errorf("failed to publish to centrifugo: %w", err)
	//	}
	//}
	return nil
}

// Helper function to map iiko status to internal status code
func getStatusCode(status string) int32 {
	statusMap := map[string]int32{
		"Unconfirmed":      0,
		"WaitCooking":      1,
		"ReadyForCooking":  2,
		"CookingStarted":   3,
		"CookingCompleted": 4,
		"Waiting":          5,
		"OnWay":            6,
		"Delivered":        7,
		"Closed":           8,
		"Cancelled":        9,
	}
	return statusMap[status]
}

// IikoClient implementation methods

func NewIikoClient(baseURL string) *IikoClient {
	return &IikoClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (c *IikoClient) Authenticate(ctx context.Context, apiLogin string) (string, error) {
	reqBody := struct {
		ApiLogin string `json:"apiLogin"`
	}{
		ApiLogin: apiLogin,
	}

	var response AuthResponse
	if err := c.doRequest(ctx, "POST", "/api/1/access_token", reqBody, &response, nil); err != nil {
		return "", err
	}

	return response.Token, nil
}

func (c *IikoClient) GetOrdersByRevision(ctx context.Context, token string, organizationIds []string, startRevision int64) ([]Order, int64, error) {
	reqBody := struct {
		OrganizationIds []string `json:"organizationIds"`
		StartRevision   int64    `json:"startRevision"`
	}{
		OrganizationIds: organizationIds,
		StartRevision:   startRevision,
	}

	fmt.Println(reqBody.StartRevision, reqBody.OrganizationIds)

	var response OrdersResponse
	if err := c.doRequest(ctx, "POST", "/api/1/deliveries/by_revision", reqBody, &response, &token); err != nil {
		return nil, 0, err
	}

	var orders []Order
	for _, orgOrders := range response.OrdersByOrganizations {
		orders = append(orders, orgOrders.Orders...)
	}

	return orders, response.MaxRevision, nil
}

func (c *IikoClient) GetOrganizationIds(ctx context.Context, token string) ([]string, error) {
	type request struct {
		OrganizationIds *struct{} `json:"organizationIds"`
	}
	req := request{OrganizationIds: nil}

	type OrganizationResponse struct {
		CorrelationId string `json:"correlationId"`
		Organizations []struct {
			ResponseType string `json:"responseType"`
			Id           string `json:"id"`
			Name         string `json:"name"`
			Code         string `json:"code"`
		} `json:"organizations"`
	}

	var response OrganizationResponse
	if err := c.doRequest(ctx, "POST", "/api/1/organizations", req, &response, &token); err != nil {
		return nil, err
	}

	orgIds := make([]string, 0, len(response.Organizations))
	for _, org := range response.Organizations {
		orgIds = append(orgIds, org.Id)
	}

	return orgIds, nil
}

// GetInitialRevision fetches initial revision by getting orders from last 24 hours
func (c *IikoClient) GetInitialRevision(ctx context.Context, token string, organizationIds []string) (int64, error) {
	reqBody := struct {
		OrganizationIds  []string `json:"organizationIds"`
		DeliveryDateFrom string   `json:"deliveryDateFrom"`
		DeliveryDateTo   string   `json:"deliveryDateTo"`
		Statuses         []string `json:"statuses,omitempty"`
	}{
		OrganizationIds: organizationIds,
		// Format date in iiko's required format: yyyy-MM-dd HH:mm:ss.fff
		DeliveryDateFrom: time.Now().Add(-3 * time.Hour).Format("2006-01-02 15:04:05.000"),
		DeliveryDateTo:   time.Now().Format("2006-01-02 15:04:05.000"),
		// Get orders in all statuses
		Statuses: []string{
			"Unconfirmed", "WaitCooking", "ReadyForCooking",
			"CookingStarted", "CookingCompleted", "Waiting",
			"OnWay", "Delivered", "Closed", "Cancelled",
		},
	}

	var response struct {
		CorrelationId         string `json:"correlationId"`
		OrdersByOrganizations []struct {
			OrganizationId string  `json:"organizationId"`
			Orders         []Order `json:"orders"`
		} `json:"ordersByOrganizations"`
		MaxRevision int64 `json:"maxRevision"`
	}

	if err := c.doRequest(ctx, "POST", "/api/1/deliveries/by_delivery_date_and_status", reqBody, &response, &token); err != nil {
		return 0, fmt.Errorf("failed to get initial revision: %w", err)
	}

	return response.MaxRevision, nil
}

func (c *IikoClient) doRequest(ctx context.Context, method, path string, reqBody, response interface{}, token *string) error {
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bytes.NewBuffer(reqJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if token != nil {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", *token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}

		if errorResponse.ErrorType == "TOO_OLD_REVISION" {
			return OldRevision
		}
		return fmt.Errorf("failed to do request: %w", errorResponse.Error())
	}

	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
