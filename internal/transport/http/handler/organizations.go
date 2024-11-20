package handler

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgx/v5/pgtype"
	"math/big"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"

	"github.com/danielgtaylor/huma/v2"
)

// Input/Output structures
type createOrganizationIn struct {
	Body struct {
		Name         string  `json:"name" maxLength:"100" doc:"Organization name" validate:"required"`
		Balance      float64 `json:"balance" doc:"Initial balance" validate:"gte=0"`
		IikoApiToken string  `json:"iiko_api_token" doc:"iiko API token" validate:"required"`
	}
}

type organizationResponse struct {
	ID           int64   `json:"id" doc:"Organization ID"`
	Name         string  `json:"name" doc:"Organization name"`
	Balance      float64 `json:"balance" doc:"Current balance"`
	IikoApiToken string  `json:"iiko_api_token" doc:"iiko API token"`
}

type updateOrganizationIn struct {
	ID   int64 `path:"id" doc:"Organization ID to update"`
	Body struct {
		Name         string  `json:"name" maxLength:"100" doc:"Organization name"`
		Balance      float64 `json:"balance" doc:"New balance"`
		IikoApiToken string  `json:"iiko_api_token" doc:"iiko API token"`
	}
}

type listOrganizationsOut struct {
	Body struct {
		Organizations []organizationResponse `json:"organizations" doc:"List of organizations"`
	}
}

// CreateOrganization creates a new organization
func CreateOrganization(ctx context.Context, in *createOrganizationIn) (*organizationResponse, error) {
	// Start transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to begin transaction: %v", err)
		return nil, huma.Error500InternalServerError("failed to create organization", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Create organization
	org, err := qtx.CreateOrganization(ctx, db.CreateOrganizationParams{
		Name:         in.Body.Name,
		IikoApiToken: in.Body.IikoApiToken,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to create organization: %v", err)
		return nil, huma.Error500InternalServerError("failed to create organization", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("failed to commit transaction", err)
	}

	return &organizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		Balance:      float64(org.Balance.Int.Int64()),
		IikoApiToken: org.IikoApiToken,
	}, nil
}

// GetOrganization retrieves organization details
func GetOrganization(ctx context.Context, in *idPathIn) (*organizationResponse, error) {
	org, err := db.Repository.GetOrganization(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("organization not found")
		}
		log.SugaredLogger.Errorf("failed to get organization: %v", err)
		return nil, huma.Error500InternalServerError("failed to get organization", err)
	}

	// Check if user has permission to view this organization
	userOrgID := ctx.Value("organization_id").(int64)
	if userOrgID != org.ID {
		return nil, huma.Error403Forbidden("not authorized to view this organization")
	}

	return &organizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		Balance:      float64(org.Balance.Int.Int64()),
		IikoApiToken: org.IikoApiToken,
	}, nil
}

// UpdateOrganization updates organization details
func UpdateOrganization(ctx context.Context, in *updateOrganizationIn) (*organizationResponse, error) {
	// Check if organization exists and user has permission
	currentOrg, err := db.Repository.GetOrganization(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("organization not found")
		}
		return nil, huma.Error500InternalServerError("failed to get organization", err)
	}

	userOrgID := ctx.Value("organization_id").(int64)
	if userOrgID != currentOrg.ID {
		return nil, huma.Error403Forbidden("not authorized to update this organization")
	}

	// Start transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("update failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Update organization
	org, err := qtx.UpdateOrganization(ctx, db.UpdateOrganizationParams{
		ID:   in.ID,
		Name: in.Body.Name,
		// FIXME: NUMERIC FUCK YOU FUCKING FUCK
		Balance: pgtype.Numeric{
			Int:              big.NewInt(int64(in.Body.Balance)),
			Exp:              0,
			NaN:              false,
			InfinityModifier: 0,
			Valid:            true,
		},
		IikoApiToken: in.Body.IikoApiToken,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to update organization: %v", err)
		return nil, huma.Error500InternalServerError("failed to update organization", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("failed to commit transaction", err)
	}

	return &organizationResponse{
		ID:           org.ID,
		Name:         org.Name,
		Balance:      float64(org.Balance.Int.Int64()),
		IikoApiToken: org.IikoApiToken,
	}, nil
}

// DeleteOrganization deletes an organization
func DeleteOrganization(ctx context.Context, in *idPathIn) (*successOut, error) {
	// Check if organization exists and user has permission
	org, err := db.Repository.GetOrganization(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("organization not found")
		}
		return nil, huma.Error500InternalServerError("failed to get organization", err)
	}

	userOrgID := ctx.Value("organization_id").(int64)
	if userOrgID != org.ID {
		return nil, huma.Error403Forbidden("not authorized to delete this organization")
	}

	// Start transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("deletion failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Delete all users in the organization
	if err := qtx.DeleteOrganizationUsers(ctx, in.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete organization users", err)
	}

	// Delete all branches
	if err := qtx.DeleteOrganizationBranches(ctx, in.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete organization branches", err)
	}

	// Delete organization
	if err := qtx.DeleteOrganization(ctx, in.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete organization", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("failed to commit transaction", err)
	}

	return &successOut{Body: struct {
		Success bool `json:"success" example:"true" doc:"Status of succession"`
	}(struct {
		Success bool "json:\"success\""
	}{Success: true})}, nil
}

// ListOrganizations retrieves all organizations
func ListOrganizations(ctx context.Context, _ *struct{}) (*listOrganizationsOut, error) {
	// Only admin can list all organizations
	if !isAdmin(ctx) {
		return nil, huma.Error403Forbidden("only admins can list all organizations")
	}

	orgs, err := db.Repository.ListOrganizations(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to list organizations: %v", err)
		return nil, huma.Error500InternalServerError("failed to list organizations", err)
	}

	var resp listOrganizationsOut
	for _, org := range orgs {
		resp.Body.Organizations = append(resp.Body.Organizations, organizationResponse{
			ID:           org.ID,
			Name:         org.Name,
			Balance:      float64(org.Balance.Int.Int64()),
			IikoApiToken: org.IikoApiToken,
		})
	}

	return &resp, nil
}

// GetOrganizationApiTokens retrieves API tokens for organizations with positive balance
func GetOrganizationApiTokens(ctx context.Context) ([]string, error) {
	// Only system processes should call this
	if !isSystemProcess(ctx) {
		return nil, huma.Error403Forbidden("unauthorized")
	}

	tokens, err := db.Repository.GetOrganizationsApiTokens(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to get API tokens: %v", err)
		return nil, huma.Error500InternalServerError("failed to get API tokens", err)
	}

	return tokens, nil
}

// Helper functions
func isAdmin(ctx context.Context) bool {
	roles := ctx.Value("user_roles").([]string)
	for _, role := range roles {
		if role == "admin" {
			return true
		}
	}
	return false
}

func isSystemProcess(ctx context.Context) bool {
	return ctx.Value("is_system_process").(bool)
}
