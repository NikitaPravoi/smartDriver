package main

import (
	"context"
	"fmt"
	"log"
	"smartDriver/internal/config"
	"smartDriver/internal/db"
	"smartDriver/pkg/iiko"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	if err := db.InitConnection(cfg); err != nil {
		fmt.Errorf("failed to init database connection: %v", err)
	}

	service := iiko.NewOrderPollingService(
		db.Pool,
		db.Repository,
		"http://localhost:8000", // Centrifugo URL
		"your-api-key",          // Centrifugo API key
		time.Second*60,          // Poll interval
	)

	if err := service.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Keep the service running
	select {}
}
