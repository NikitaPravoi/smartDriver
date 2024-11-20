package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

// Input/Output structures
type createRideIn struct {
	Body struct {
		BranchID int64   `json:"branch_id" doc:"Branch ID"`
		OrderIDs []int64 `json:"order_ids" doc:"List of order IDs to attach to the ride"`
	}
}

type updateRideIn struct {
	ID   int64 `path:"id" doc:"Ride ID"`
	Body struct {
		OrderIDs []int64 `json:"order_ids" doc:"List of order IDs to attach to the ride"`
	}
}

type deleteRideIn struct {
	ID int64 `path:"id" doc:"Ride ID"`
}

type rideOut struct {
	Body struct {
		ID        int64       `json:"id" doc:"Ride ID"`
		BranchID  int64       `json:"branch_id" doc:"Branch ID"`
		CreatedAt time.Time   `json:"created_at" doc:"Ride created timestamp"`
		EndedAt   time.Time   `json:"ended_at,omitempty" doc:"Ride ended timestamp"`
		Orders    []orderInfo `json:"orders,omitempty" doc:"List of orders attached to the ride"`
	}
}

type point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// CreateRide creates a new ride and optionally attaches orders to it
func CreateRide(ctx context.Context, in *createRideIn) (*rideOut, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to begin transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to create ride", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Create the ride
	ride, err := qtx.CreateRide(ctx, in.Body.BranchID)
	if err != nil {
		log.SugaredLogger.Errorf("failed to create ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to create ride", err)
	}

	// Attach orders if provided
	if len(in.Body.OrderIDs) > 0 {
		for _, orderID := range in.Body.OrderIDs {
			err := qtx.AttachOrderToRide(ctx, db.AttachOrderToRideParams{
				RideID:  ride.ID,
				OrderID: orderID,
			})
			if err != nil {
				log.SugaredLogger.Errorf("failed to attach order %d to ride: %v", orderID, err)
				return nil, huma.Error500InternalServerError("failed to attach orders to ride", err)
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.SugaredLogger.Errorf("failed to commit transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to create ride", err)
	}

	return buildRideResponse(ctx, *qtx, ride)
}

// GetRide retrieves a ride by ID including its attached orders
func GetRide(ctx context.Context, in *idPathIn) (*rideOut, error) {
	ride, err := db.Repository.GetRide(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("ride not found")
		}
		log.SugaredLogger.Errorf("failed to get ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to get ride", err)
	}

	return buildRideResponse(ctx, *db.Repository, ride)
}

// UpdateRide updates the orders attached to a ride
func UpdateRide(ctx context.Context, in *updateRideIn) (*rideOut, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to begin transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to update ride", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Get the ride to ensure it exists and isn't completed
	ride, err := qtx.GetRide(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("ride not found")
		}
		log.SugaredLogger.Errorf("failed to get ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to update ride", err)
	}

	if !ride.EndedAt.Time.IsZero() {
		return nil, huma.Error400BadRequest("cannot update completed ride")
	}

	// Remove all existing order associations
	if err := qtx.DetachAllOrdersFromRide(ctx, ride.ID); err != nil {
		log.SugaredLogger.Errorf("failed to detach orders from ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to update ride", err)
	}

	// Attach new orders
	for _, orderID := range in.Body.OrderIDs {
		err := qtx.AttachOrderToRide(ctx, db.AttachOrderToRideParams{
			RideID:  ride.ID,
			OrderID: orderID,
		})
		if err != nil {
			log.SugaredLogger.Errorf("failed to attach order %d to ride: %v", orderID, err)
			return nil, huma.Error500InternalServerError("failed to update ride", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		log.SugaredLogger.Errorf("failed to commit transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to update ride", err)
	}

	return buildRideResponse(ctx, *db.Repository, ride)
}

// DeleteRide marks a ride as completed
func DeleteRide(ctx context.Context, in *deleteRideIn) (*successOut, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to begin transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete ride", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Get the ride to ensure it exists
	ride, err := qtx.GetRide(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("ride not found")
		}
		log.SugaredLogger.Errorf("failed to get ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete ride", err)
	}

	// Mark ride as completed
	if err := qtx.CompleteRide(ctx, in.ID); err != nil {
		log.SugaredLogger.Errorf("failed to complete ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete ride", err)
	}

	// Detach all orders
	if err := qtx.DetachAllOrdersFromRide(ctx, ride.ID); err != nil {
		log.SugaredLogger.Errorf("failed to detach orders from ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete ride", err)
	}

	if err := tx.Commit(ctx); err != nil {
		log.SugaredLogger.Errorf("failed to commit transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete ride", err)
	}

	return &successOut{Body: struct {
		Success bool `json:"success" example:"true" doc:"Status of succession"`
	}(struct {
		Success bool "json:\"success\""
	}{Success: true})}, nil
}

// Helper function to build ride response with orders
func buildRideResponse(ctx context.Context, q db.Queries, ride db.Ride) (*rideOut, error) {
	orders, err := q.GetOrdersByRideID(ctx, ride.ID)
	if err != nil {
		log.SugaredLogger.Errorf("failed to get orders for ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to get ride details", err)
	}

	var resp rideOut
	resp.Body.ID = ride.ID
	resp.Body.BranchID = ride.BranchID
	resp.Body.CreatedAt = ride.CreatedAt.Time
	resp.Body.EndedAt = ride.EndedAt.Time

	for _, order := range orders {
		resp.Body.Orders = append(resp.Body.Orders, orderInfo{
			ID:           order.ID,
			ExternalID:   order.ExternalID,
			Status:       *order.Status,
			Address:      formatAddress(order),
			Location:     point{Lat: order.Location.P.Y, Lng: order.Location.P.X},
			CustomerName: order.CustomerName,
		})
	}

	return &resp, nil
}

// Helper function to format address
func formatAddress(order db.Order) string {
	return fmt.Sprintf("%s, %s, д. %s, кв. %s",
		order.City,
		order.Street,
		order.Building,
		order.Apartment,
	)
}
