package http

import (
	"net/http"
	"smartDriver/internal/transport/http/handler"

	"github.com/danielgtaylor/huma/v2"
)

func Register(api huma.API) {
	// Authorization endpoints
	huma.Register(api, huma.Operation{
		OperationID:   "register",
		Method:        http.MethodPost,
		Path:          "/register",
		Summary:       "Register new user",
		Description:   "Create a new user account",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusCreated,
	}, handler.Register)

	//huma.Register(api, huma.Operation{
	//	OperationID:   "create-user",
	//	Method:        http.MethodPost,
	//	Path:          "/users",
	//	Summary:       "Create user",
	//	Description:   "Create a user and return it",
	//	Tags:          []string{"Authorization"},
	//	DefaultStatus: http.StatusCreated,
	//}, handler.Register)

	huma.Register(api, huma.Operation{
		OperationID:   "get-user",
		Method:        http.MethodGet,
		Path:          "/users/{id}",
		Summary:       "Get user",
		Description:   "Get a user by ID",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.GetUser)

	huma.Register(api, huma.Operation{
		OperationID:   "update-user",
		Method:        http.MethodPut,
		Path:          "/users/{id}",
		Summary:       "Update user",
		Description:   "Update a user by ID",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.UpdateUser)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-user",
		Method:        http.MethodDelete,
		Path:          "/users/{id}",
		Summary:       "Delete user",
		Description:   "Delete a user by ID",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.DeleteUser)

	//huma.Register(api, huma.Operation{
	//	OperationID:   "list-users",
	//	Method:        http.MethodGet,
	//	Path:          "/users",
	//	Summary:       "List users",
	//	Description:   "List users with pagination and filtering",
	//	Tags:          []string{"Authorization"},
	//	DefaultStatus: http.StatusOK,
	//}, handler.ListUsers)

	huma.Register(api, huma.Operation{
		OperationID:   "login",
		Method:        http.MethodPost,
		Path:          "/login",
		Summary:       "Log in",
		Description:   "Authenticate user and create session",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.Login)

	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodPost, // Changed to POST as it modifies state
		Path:          "/logout",
		Summary:       "Log out",
		Description:   "End current session",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.Logout)

	huma.Register(api, huma.Operation{
		OperationID:   "refresh-token",
		Method:        http.MethodPost,
		Path:          "/refresh-token",
		Summary:       "Refresh token",
		Description:   "Refresh session token using refresh token",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.RefreshToken)

	// Organizations endpoints
	huma.Register(api, huma.Operation{
		OperationID:   "create-organization",
		Method:        http.MethodPost,
		Path:          "/organizations",
		Summary:       "Create organization",
		Description:   "Create a new organization",
		Tags:          []string{"Organizations"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateOrganization)

	huma.Register(api, huma.Operation{
		OperationID:   "get-organization",
		Method:        http.MethodGet,
		Path:          "/organizations/{id}",
		Summary:       "Get organization",
		Description:   "Get organization details by ID",
		Tags:          []string{"Organizations"},
		DefaultStatus: http.StatusOK,
	}, handler.GetOrganization)

	huma.Register(api, huma.Operation{
		OperationID:   "update-organization",
		Method:        http.MethodPut,
		Path:          "/organizations/{id}",
		Summary:       "Update organization",
		Description:   "Update organization details",
		Tags:          []string{"Organizations"},
		DefaultStatus: http.StatusOK,
	}, handler.UpdateOrganization)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-organization",
		Method:        http.MethodDelete,
		Path:          "/organizations/{id}",
		Summary:       "Delete organization",
		Description:   "Delete organization and related data",
		Tags:          []string{"Organizations"},
		DefaultStatus: http.StatusOK,
	}, handler.DeleteOrganization)

	huma.Register(api, huma.Operation{
		OperationID:   "list-organizations",
		Method:        http.MethodGet,
		Path:          "/organizations",
		Summary:       "List organizations",
		Description:   "List all organizations (admin only)",
		Tags:          []string{"Organizations"},
		DefaultStatus: http.StatusOK,
	}, handler.ListOrganizations)

	// plans
	huma.Register(api, huma.Operation{
		OperationID:   "create-plan",
		Method:        http.MethodPost,
		Path:          "/plans",
		Summary:       "Create plan",
		Description:   "Create a plan and return it",
		Tags:          []string{"Plans"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreatePlan)
	huma.Register(api, huma.Operation{
		OperationID:   "get-plan",
		Method:        http.MethodGet,
		Path:          "/plans/{id}",
		Summary:       "Get plan",
		Description:   "Get a plan by ID",
		Tags:          []string{"Plans"},
		DefaultStatus: http.StatusOK,
	}, handler.GetPlan)
	huma.Register(api, huma.Operation{
		OperationID:   "update-plan",
		Method:        http.MethodPut,
		Path:          "/plans/{id}",
		Summary:       "Update plan",
		Description:   "Update a plan by ID",
		Tags:          []string{"Plans"},
		DefaultStatus: http.StatusOK,
	}, handler.UpdatePlan)
	huma.Register(api, huma.Operation{
		OperationID:   "delete-plan",
		Method:        http.MethodDelete,
		Path:          "/plans/{id}",
		Summary:       "Delete plan",
		Description:   "Delete a plan by ID",
		Tags:          []string{"Plans"},
		DefaultStatus: http.StatusOK,
	}, handler.DeletePlan)
	huma.Register(api, huma.Operation{
		OperationID:   "list-plans",
		Method:        http.MethodGet,
		Path:          "/plans",
		Summary:       "List plans",
		Description:   "List a plans with pagination",
		Tags:          []string{"Plans"},
		DefaultStatus: http.StatusOK,
	}, handler.ListPlans)

	// Rides endpoints
	huma.Register(api, huma.Operation{
		OperationID:   "create-ride",
		Method:        http.MethodPost,
		Path:          "/rides",
		Summary:       "Create ride",
		Description:   "Create a new ride with optional orders",
		Tags:          []string{"Rides"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateRide)

	huma.Register(api, huma.Operation{
		OperationID:   "get-ride",
		Method:        http.MethodGet,
		Path:          "/rides/{id}",
		Summary:       "Get ride",
		Description:   "Get ride details with attached orders",
		Tags:          []string{"Rides"},
		DefaultStatus: http.StatusOK,
	}, handler.GetRide)

	huma.Register(api, huma.Operation{
		OperationID:   "update-ride",
		Method:        http.MethodPut,
		Path:          "/rides/{id}",
		Summary:       "Update ride",
		Description:   "Update ride orders",
		Tags:          []string{"Rides"},
		DefaultStatus: http.StatusOK,
	}, handler.UpdateRide)

	huma.Register(api, huma.Operation{
		OperationID:   "delete-ride",
		Method:        http.MethodDelete,
		Path:          "/rides/{id}",
		Summary:       "Delete ride",
		Description:   "Complete ride and detach orders",
		Tags:          []string{"Rides"},
		DefaultStatus: http.StatusOK,
	}, handler.DeleteRide)

	// Orders endpoints
	huma.Register(api, huma.Operation{
		OperationID:   "get-unbound-orders",
		Method:        http.MethodGet,
		Path:          "/orders/unbound",
		Summary:       "Get unbound orders",
		Description:   "Get orders not attached to any ride",
		Tags:          []string{"Orders"},
		DefaultStatus: http.StatusOK,
	}, handler.GetUnboundOrders)

	//huma.Register(api, huma.Operation{
	//	OperationID:   "get-order-statuses",
	//	Method:        http.MethodGet,
	//	Path:          "/orders/statuses",
	//	Summary:       "Get order statuses",
	//	Description:   "Get list of available order statuses",
	//	Tags:          []string{"Orders"},
	//	DefaultStatus: http.StatusOK,
	//}, handler.GetOrderStatuses)

	// Internal endpoints (requires system authentication)
	//huma.Register(api, huma.Operation{
	//	OperationID:   "get-organization-tokens",
	//	Method:        http.MethodGet,
	//	Path:          "/internal/organization-tokens",
	//	Summary:       "Get organization tokens",
	//	Description:   "Get API tokens for organizations with positive balance",
	//	Tags:          []string{"Internal"},
	//	DefaultStatus: http.StatusOK,
	//}, handler.GetOrganizationApiTokens)
}
