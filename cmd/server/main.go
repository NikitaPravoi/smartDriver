package main

import (
	"context"
	"fmt"
	"net/http"
	"smartDriver/internal/config"
	"smartDriver/internal/db"
	httptransport "smartDriver/internal/transport/http"
	"smartDriver/pkg/log"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		session, err := db.Repository.GetSessionByToken(r.Context(), token)
		if err != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", session.UserID)
		ctx = context.WithValue(ctx, "session_token", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {
	cfg, err := config.Load()
	log.MustInit(cfg)
	if err != nil {
		log.SugaredLogger.Fatalf("Failed to load configuration: %v", err)
	}

	if err := db.InitConnection(cfg); err != nil {
		log.SugaredLogger.Errorf("failed to init database connection: %v", err)
	}

	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("SmartDriver", "0.5.3"))
	httptransport.Register(api)

	log.SugaredLogger.Infof("HTTP server is listening on %s:%s", cfg.Server.Host, cfg.Server.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.Server.Port), router); err != nil {
		log.SugaredLogger.Fatalf("failed to listen given address: %v", err)
	}
}
