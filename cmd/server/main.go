package main

import (
	"net/http"
	"smartDriver/internal/handlers"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

var api huma.API

func main() {
	router := chi.NewMux()
	api = humachi.New(router, huma.DefaultConfig("SmartDriver", "0.0.1"))

	huma.Register(api, huma.Operation{
		OperationID:   "create-user",
		Method:        http.MethodPost,
		Path:          "/user",
		Summary:       "Create a user",
		Description:   "Creating user and returning state of creation",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusCreated,
	}, handlers.CreateUserHandler)
	huma.Register(api, huma.Operation{
		OperationID:   "login",
		Method:        http.MethodPost,
		Path:          "/login",
		Summary:       "Log in",
		Description:   "Logging in",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handlers.LoginHandler)
	huma.Register(api, huma.Operation{
		OperationID:   "logout",
		Method:        http.MethodGet,
		Path:          "/logout",
		Summary:       "Log out",
		Description:   "Logging out",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusOK,
	}, handlers.LogoutHandler)
	huma.Register(api, huma.Operation{
		OperationID:   "refresh-token",
		Method:        http.MethodPost,
		Path:          "/refresh-token",
		Summary:       "Refresh token",
		Description:   "Refreshing session token expiry",
		Tags:          []string{"Authorization"},
		DefaultStatus: http.StatusCreated,
	}, handlers.RefreshTokenHandler)

	// Start the server!
	http.ListenAndServe(":8888", router)
}
