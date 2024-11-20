package http

import (
	"net/http"
	"smartDriver/internal/transport/http/handler"

	"github.com/danielgtaylor/huma/v2"
)

func Register(api huma.API) {
	// authorization
	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/users",
		Summary:       "Create user",
		Description:   "Create a user and return it",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateUser)
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
	huma.Register(api, huma.Operation{
		OperationID:   "list-users",
		Method:        http.MethodGet,
		Path:          "/users",
		Summary:       "List users",
		Description:   "List a users with pagination",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.ListUsers)
	huma.Register(api, huma.Operation{
		OperationID:   "login",
		Method:        http.MethodPost,
		Path:          "/login",
		Summary:       "Log in",
		Description:   "Logging in",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.Login)
	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodGet,
		Path:          "/logout",
		Summary:       "Log out",
		Description:   "Logging out",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handler.Logout)
	huma.Register(api, huma.Operation{
		OperationID:   "refresh-token",
		Method:        http.MethodPost,
		Path:          "/refresh-token",
		Summary:       "Refresh token",
		Description:   "Refreshing session token expiry",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusCreated,
	}, handler.RefreshToken)

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

	// rides
	huma.Register(api, huma.Operation{
		OperationID:   "create-ride",
		Method:        http.MethodPost,
		Path:          "/rides",
		Summary:       "Create ride",
		Description:   "Create a ride and return it",
		Tags:          []string{"Rides"},
		DefaultStatus: http.StatusCreated,
	}, handler.CreateRide)
	huma.Register(api, huma.Operation{
		OperationID:   "get-ride",
		Method:        http.MethodGet,
		Path:          "/rides/{id}",
		Summary:       "Get ride",
		Description:   "Get a ride by ID",
		Tags:          []string{"rides"},
		DefaultStatus: http.StatusOK,
	}, handler.GetRide)
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "update-ride",
	// 	Method:        http.MethodPut,
	// 	Path:          "/rides/{id}",
	// 	Summary:       "Update a ride",
	// 	Description:   "Update a ride by ID",
	// 	Tags:          []string{"rides"},
	// 	DefaultStatus: http.StatusOK,
	// }, handler.UpdateRide)
	// huma.Register(api, huma.Operation{
	// 	OperationID:   "delete-ride",
	// 	Method:        http.MethodDelete,
	// 	Path:          "/users/{id}",
	// 	Summary:       "Delete a ride",
	// 	Description:   "Delete a ride by ID",
	// 	Tags:          []string{"rides"},
	// 	DefaultStatus: http.StatusOK,
	// }, handler.DeleteRide)
}
