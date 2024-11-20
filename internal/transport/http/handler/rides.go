package handler

import (
	"context"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type createRideIn struct {
	Body struct {
		BranchID int64 `json:"branch_id" doc:"Branch ID"`
	}
}

type rideOut struct {
	Body struct {
		ID        int64     `json:"id" doc:"Ride ID"`
		BranchID  int64     `json:"branch_id" doc:"Branch ID"`
		CreatedAt time.Time `json:"created_at" doc:"Ride created timestamp"`
		EndedAt   time.Time `json:"ended_at" doc:"Ride ended timestamp"`
	}
}

func CreateRide(ctx context.Context, in *createRideIn) (*rideOut, error) {
	ride, err := db.Pool.CreateRide(ctx, in.Body.BranchID)
	if err != nil {
		log.SugaredLogger.Errorf("failed to create ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to create ride", err)
	}

	var resp rideOut
	resp.Body.ID = ride.ID
	resp.Body.BranchID = ride.BranchID
	resp.Body.CreatedAt = ride.CreatedAt.Time
	resp.Body.EndedAt = ride.EndedAt.Time

	return &resp, nil
}

func GetRide(ctx context.Context, in *idPathIn) (*rideOut, error) {
	ride, err := db.Pool.CreateRide(ctx, in.ID)
	if err != nil {
		log.SugaredLogger.Errorf("failed to get ride: %v", err)
		return nil, huma.Error500InternalServerError("failed to get ride", err)
	}

	var resp rideOut
	resp.Body.ID = ride.ID
	resp.Body.BranchID = ride.BranchID
	resp.Body.CreatedAt = ride.CreatedAt.Time
	resp.Body.EndedAt = ride.EndedAt.Time

	return &resp, nil
}

func UpdateRide() {}

func DeleteRide() {}
