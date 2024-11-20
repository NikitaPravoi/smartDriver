package handler

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// Input structure for filtering unbound orders
type getUnboundOrdersIn struct {
	Query struct {
		Status   *string   `query:"status" doc:"Filter by order status"`
		FromDate time.Time `query:"from_date" doc:"Filter orders from this date (RFC3339 format)"`
		ToDate   time.Time `query:"to_date" doc:"Filter orders until this date (RFC3339 format)"`
		Limit    int32     `query:"limit" default:"50" doc:"Maximum number of orders to return"`
		Offset   int32     `query:"offset" default:"0" doc:"Number of orders to skip"`
	}
}

// Output structure for the unbound orders list
type unboundOrdersOut struct {
	Body struct {
		Orders []orderInfo `json:"orders" doc:"List of unbound orders"`
		Total  int64       `json:"total" doc:"Total number of unbound orders"`
		Limit  int32       `json:"limit" doc:"Current page limit"`
		Offset int32       `json:"offset" doc:"Current page offset"`
	}
}

// GetUnboundOrders retrieves orders that aren't attached to any rides
func GetUnboundOrders(ctx context.Context, in *getUnboundOrdersIn) (*unboundOrdersOut, error) {
	params := db.GetUnboundOrdersParams{
		Limit:  in.Query.Limit,
		Offset: in.Query.Offset,
	}

	if !in.Query.FromDate.IsZero() {
		params.CreatedAt = pgtype.Timestamp{
			Time:  in.Query.FromDate,
			Valid: true,
		}
	}

	if !in.Query.ToDate.IsZero() {
		params.CreatedAt_2 = pgtype.Timestamp{
			Time:  in.Query.ToDate,
			Valid: true,
		}
	}

	if in.Query.Status != nil {
		params.Status = in.Query.Status

	}

	// Get unbound orders
	orders, err := db.Repository.GetUnboundOrders(ctx, params)
	if err != nil {
		log.SugaredLogger.Errorf("failed to get unbound orders: %v", err)
		return nil, huma.Error500InternalServerError("failed to get unbound orders", err)
	}

	// Get total count
	total, err := db.Repository.CountUnboundOrders(ctx, db.CountUnboundOrdersParams{
		Status:      params.Status,
		CreatedAt:   params.CreatedAt,
		CreatedAt_2: params.CreatedAt_2,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to count unbound orders: %v", err)
		return nil, huma.Error500InternalServerError("failed to get unbound orders count", err)
	}

	var resp unboundOrdersOut
	resp.Body.Limit = in.Query.Limit
	resp.Body.Offset = in.Query.Offset
	resp.Body.Total = total

	for _, order := range orders {
		resp.Body.Orders = append(resp.Body.Orders, orderInfo{
			ID:           order.ID,
			ExternalID:   order.ExternalID,
			Status:       *order.Status,
			Address:      formatAddress(order),
			Location:     point{Lat: order.Location.P.Y, Lng: order.Location.P.X},
			CustomerName: order.CustomerName,
			CreatedAt:    order.CreatedAt.Time,
			Cost:         order.Cost.Int.Int64(),
		})
	}

	return &resp, nil
}

// Add this to the existing orderInfo struct
type orderInfo struct {
	ID           int64     `json:"id"`
	ExternalID   string    `json:"external_id"`
	Status       string    `json:"status"`
	Address      string    `json:"address"`
	Location     point     `json:"location"`
	CustomerName string    `json:"customer_name"`
	CreatedAt    time.Time `json:"created_at"`
	Cost         int64     `json:"cost"`
}
