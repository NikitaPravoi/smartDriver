package db

import (
	"context"
	"fmt"
	"smartDriver/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Repository *Queries
var Pool *pgxpool.Pool

// InitConnection create new database connection pool and assign it to pointer on Queries
// struct that was generated by SQLC. after call this function you will be able
// to interact with database through Repository global variable.
//
// for example:
//
//	database.Repository.CreateGeolocation(...)
func InitConnection(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
	))
	if err != nil {
		return err
	}

	Repository = New(pool)
	Pool = pool
	return nil
}
