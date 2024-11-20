package main

import (
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

func main() {
	cfg := config.MustLoad()

	log.MustInit(cfg)

	if err := db.InitConnection(cfg); err != nil {
		log.SugaredLogger.Errorf("failed to init database connection: %v", err)
	}

	router := chi.NewMux()
	api := humachi.New(router, huma.DefaultConfig("SmartDriver", "0.0.1"))
	httptransport.Register(api)

	log.SugaredLogger.Infof("HTTP server is listening on %s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.HTTP.Port), router); err != nil {
		log.SugaredLogger.Fatalf("failed to listen given address: %v", err)
	}
}
