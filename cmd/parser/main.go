package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"smartDriver/internal/db"
	"smartDriver/pkg/iiko"
	"time"
)

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:1234qwerASDF@localhost:5432/smart_driver?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to create postgresql pool: %v", err)
	}
	queries := db.New(pool)

	service := iiko.NewOrderPollingService(
		pool,
		queries,
		"http://localhost:8000", // Centrifugo URL
		"your-api-key",          // Centrifugo API key
		time.Second*10,          // Poll interval
	)

	if err := service.Start(ctx); err != nil {
		log.Fatal(err)
	}

	// Keep the service running
	select {}
}
