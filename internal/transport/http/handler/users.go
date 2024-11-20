package handler

import (
	"context"
	"database/sql"
	"errors"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register Input/Output structures
type registerIn struct {
	Body struct {
		Login          string `json:"login" maxLength:"50" doc:"User login" validate:"required,email"`
		Password       string `json:"password" maxLength:"50" doc:"User password" validate:"required,min=8"`
		Name           string `json:"name" maxLength:"50" doc:"User name" validate:"required"`
		Surname        string `json:"surname" maxLength:"50" doc:"User surname" validate:"required"`
		Patronymic     string `json:"patronymic" maxLength:"50" doc:"User patronymic"`
		OrganizationID int64  `json:"organization_id" example:"1" doc:"Organization ID" validate:"required"`
	}
}

type userResponse struct {
	ID             int64     `json:"id" doc:"User ID"`
	Login          string    `json:"login" doc:"User login"`
	Name           string    `json:"name" doc:"User name"`
	Surname        string    `json:"surname" doc:"User surname"`
	Patronymic     string    `json:"patronymic" doc:"User patronymic"`
	OrganizationID int64     `json:"organization_id" doc:"Organization ID"`
	CreatedAt      time.Time `json:"created_at" doc:"Account creation time"`
	UpdatedAt      time.Time `json:"updated_at" doc:"Last update time"`
}

// Update User Input/Output structures
type updateUserIn struct {
	ID   int64 `path:"id" doc:"User ID to update"`
	Body struct {
		Name       string `json:"name" maxLength:"50" doc:"User name"`
		Surname    string `json:"surname" maxLength:"50" doc:"User surname"`
		Patronymic string `json:"patronymic" maxLength:"50" doc:"User patronymic"`
		Password   string `json:"password,omitempty" maxLength:"50" doc:"New password (optional)"`
	}
}

// List Users Input/Output structures
type listUsersIn struct {
	Query struct {
		OrganizationID *int64 `query:"organization_id" doc:"Filter users by organization"`
		Search         string `query:"search" doc:"Search in name, surname, or login"`
		Limit          int32  `query:"limit" default:"50" doc:"Maximum number of users to return"`
		Offset         int32  `query:"offset" default:"0" doc:"Number of users to skip"`
	}
}

type listUsersOut struct {
	Body struct {
		Users  []userResponse `json:"users" doc:"List of users"`
		Total  int64          `json:"total" doc:"Total number of users matching criteria"`
		Limit  int32          `json:"limit" doc:"Current page limit"`
		Offset int32          `json:"offset" doc:"Current page offset"`
	}
}

// Register creates a new user account
func Register(ctx context.Context, in *registerIn) (*userResponse, error) {
	// Check if login is already taken
	exists, err := db.Repository.CheckUserLoginExists(ctx, in.Body.Login)
	if err != nil {
		log.SugaredLogger.Errorf("failed to check login existence: %v", err)
		return nil, huma.Error500InternalServerError("registration failed", err)
	}
	if exists {
		return nil, huma.Error400BadRequest("login already taken")
	}

	// Check if organization exists
	org, err := db.Repository.GetOrganization(ctx, in.Body.OrganizationID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error400BadRequest("invalid organization id")
		}
		log.SugaredLogger.Errorf("failed to check organization: %v", err)
		return nil, huma.Error500InternalServerError("registration failed", err)
	}

	// Start transaction
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("registration failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, huma.Error500InternalServerError("registration failed", err)
	}

	// Create user
	user, err := qtx.CreateUser(ctx, db.CreateUserParams{
		Login:          in.Body.Login,
		Password:       string(hashedPassword),
		Name:           &in.Body.Name,
		Surname:        &in.Body.Surname,
		Patronymic:     &in.Body.Patronymic,
		OrganizationID: org.ID,
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create user", err)
	}

	// Assign default role
	if err := qtx.AssignUserRole(ctx, db.AssignUserRoleParams{
		UserID: user.ID,
		RoleID: 1, // Assuming 1 is your default role ID
	}); err != nil {
		return nil, huma.Error500InternalServerError("failed to assign role", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("registration failed", err)
	}

	return &userResponse{
		ID:             user.ID,
		Login:          user.Login,
		Name:           *user.Name,
		Surname:        *user.Surname,
		Patronymic:     *user.Patronymic,
		OrganizationID: user.OrganizationID,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}, nil
}

// GetUser retrieves user details
func GetUser(ctx context.Context, in *idPathIn) (*userResponse, error) {
	user, err := db.Repository.GetUser(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("user not found")
		}
		log.SugaredLogger.Errorf("failed to get user: %v", err)
		return nil, huma.Error500InternalServerError("failed to get user", err)
	}

	// Check if requester has permission to view this user
	requesterID := ctx.Value("user_id").(int64)
	requesterOrg := ctx.Value("organization_id").(int64)
	if requesterID != user.ID && requesterOrg != user.OrganizationID {
		return nil, huma.Error403Forbidden("not authorized to view this user")
	}

	return &userResponse{
		ID:             user.ID,
		Login:          user.Login,
		Name:           *user.Name,
		Surname:        *user.Surname,
		Patronymic:     *user.Patronymic,
		OrganizationID: user.OrganizationID,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}, nil
}

// UpdateUser updates user details
func UpdateUser(ctx context.Context, in *updateUserIn) (*userResponse, error) {
	// Get current user
	currentUser, err := db.Repository.GetUser(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, huma.Error500InternalServerError("failed to get user", err)
	}

	// Check permissions
	requesterID := ctx.Value("user_id").(int64)
	requesterOrg := ctx.Value("organization_id").(int64)
	if requesterID != currentUser.ID && requesterOrg != currentUser.OrganizationID {
		return nil, huma.Error403Forbidden("not authorized to update this user")
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("update failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	params := db.UpdateUserParams{
		ID:         in.ID,
		Name:       &in.Body.Name,
		Surname:    &in.Body.Surname,
		Patronymic: &in.Body.Patronymic,
	}

	// Update password if provided
	if in.Body.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(in.Body.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, huma.Error500InternalServerError("failed to hash password", err)
		}
		params.Password = string(hashedPassword)
	}

	user, err := qtx.UpdateUser(ctx, params)
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to update user", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("update failed", err)
	}

	return &userResponse{
		ID:             user.ID,
		Login:          user.Login,
		Name:           *user.Name,
		Surname:        *user.Surname,
		Patronymic:     *user.Patronymic,
		OrganizationID: user.OrganizationID,
		CreatedAt:      user.CreatedAt.Time,
		UpdatedAt:      user.UpdatedAt.Time,
	}, nil
}

// DeleteUser deletes a user account
func DeleteUser(ctx context.Context, in *idPathIn) (*successOut, error) {
	// Get user to check permissions
	user, err := db.Repository.GetUser(ctx, in.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, huma.Error404NotFound("user not found")
		}
		return nil, huma.Error500InternalServerError("failed to get user", err)
	}

	// Check permissions
	requesterOrg := ctx.Value("organization_id").(int64)
	if requesterOrg != user.OrganizationID {
		return nil, huma.Error403Forbidden("not authorized to delete this user")
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("deletion failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Delete all sessions
	if err := qtx.DeleteUserSessions(ctx, in.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete sessions", err)
	}

	// Delete user roles
	if err := qtx.DeleteUserRoles(ctx, in.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete roles", err)
	}

	// Delete user
	if err := qtx.DeleteUser(ctx, in.ID); err != nil {
		return nil, huma.Error500InternalServerError("failed to delete user", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("deletion failed", err)
	}

	return &successOut{Body: struct {
		Success bool `json:"success" example:"true" doc:"Status of succession"`
	}(struct {
		Success bool "json:\"success\""
	}{Success: true})}, nil
}

// ListUsers retrieves a paginated list of users
//func ListUsers(ctx context.Context, in *listUsersIn) (*listUsersOut, error) {
//	// Check if requester can list users
//	requesterOrg := ctx.Value("organization_id").(int64)
//	if in.Query.OrganizationID != nil && *in.Query.OrganizationID != requesterOrg {
//		return nil, huma.Error403Forbidden("can only list users from your organization")
//	}
//
//	params := db.ListUsersParams{
//		OrganizationID: sql.NullInt64{Int64: requesterOrg, Valid: true},
//		SearchQuery:    fmt.Sprintf("%%%s%%", in.Query.Search),
//		Limit:          in.Query.Limit,
//		Offset:         in.Query.Offset,
//	}
//
//	users, err := db.Repository.ListUsers(ctx, params)
//	if err != nil {
//		return nil, huma.Error500InternalServerError("failed to list users", err)
//	}
//
//	total, err := db.Repository.CountUsers(ctx, db.CountUsersParams{
//		OrganizationID: params.OrganizationID,
//		SearchQuery:    params.SearchQuery,
//	})
//	if err != nil {
//		return nil, huma.Error500InternalServerError("failed to count users", err)
//	}
//
//	var resp listUsersOut
//	resp.Body.Limit = in.Query.Limit
//	resp.Body.Offset = in.Query.Offset
//	resp.Body.Total = total
//
//	for _, user := range users {
//		resp.Body.Users = append(resp.Body.Users, userResponse{
//			ID:             user.ID,
//			Login:          user.Login,
//			Name:           user.Name,
//			Surname:        user.Surname,
//			Patronymic:     user.Patronymic.String,
//			OrganizationID: user.OrganizationID,
//			CreatedAt:      user.CreatedAt,
//			UpdatedAt:      user.UpdatedAt,
//		})
//	}
//
//	return &resp, nil
//}
