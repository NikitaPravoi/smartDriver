package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

const (
	sessionDuration = 24 * time.Hour
	refreshDuration = 7 * 24 * time.Hour
	tokenLength     = 32
)

type loginIn struct {
	Body struct {
		Username string `json:"username" example:"admin" doc:"User username"`
		Password string `json:"password" example:"admin" doc:"User password"`
	}
}

type loginOut struct {
	Body struct {
		User struct {
			ID             int64  `json:"id" doc:"User ID"`
			Login          string `json:"login" doc:"Username"`
			Name           string `json:"name" doc:"User's name"`
			Surname        string `json:"surname" doc:"User's surname"`
			Patronymic     string `json:"patronymic" doc:"User's patronymic"`
			OrganizationID int64  `json:"organization_id" doc:"Organization ID"`
		} `json:"user" doc:"User information"`
		Session struct {
			SessionToken string    `json:"session_token" doc:"Authentication token"`
			RefreshToken string    `json:"refresh_token" doc:"Token for refreshing session"`
			CreatedAt    time.Time `json:"created_at" doc:"Session creation time"`
			ExpiresAt    time.Time `json:"expires_at" doc:"Session expiration time"`
		} `json:"session" doc:"Session information"`
	}
}

type refreshTokenIn struct {
	Body struct {
		RefreshToken string `json:"refresh_token" doc:"Refresh token for extending session"`
	}
}

type refreshTokenOut struct {
	Body struct {
		SessionToken string    `json:"session_token" doc:"New session token"`
		RefreshToken string    `json:"refresh_token" doc:"New refresh token"`
		ExpiresAt    time.Time `json:"expires_at" doc:"New expiration time"`
	}
}

// Login authenticates a user and creates a new session
func Login(ctx context.Context, in *loginIn) (*loginOut, error) {
	// Get user by login
	user, err := db.Repository.GetUserByLogin(ctx, in.Body.Username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, huma.Error401Unauthorized("invalid credentials")
		}
		log.SugaredLogger.Errorf("failed to get user: %v", err)
		return nil, huma.Error500InternalServerError("login failed", err)
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Body.Password)); err != nil {
		return nil, huma.Error401Unauthorized("invalid credentials")
	}

	// Start transaction for session creation
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("login failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Generate tokens
	sessionToken, err := generateToken()
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to generate session token", err)
	}

	refreshToken, err := generateToken()
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to generate refresh token", err)
	}

	now := time.Now()
	expiresAt := now.Add(sessionDuration)

	// Create new session
	session, err := qtx.CreateSession(ctx, db.CreateSessionParams{
		UserID:       user.ID,
		SessionToken: sessionToken,
		RefreshToken: refreshToken,
		CreatedAt: pgtype.Timestamp{
			Time:             now,
			InfinityModifier: 0,
			Valid:            true,
		},
		ExpiresAt: pgtype.Timestamp{
			Time:             expiresAt,
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to create session", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("failed to commit session", err)
	}

	var resp loginOut
	resp.Body.User.ID = user.ID
	resp.Body.User.Login = user.Login
	resp.Body.User.Name = *user.Name
	resp.Body.User.Surname = *user.Surname
	resp.Body.User.Patronymic = *user.Patronymic
	resp.Body.User.OrganizationID = user.OrganizationID

	resp.Body.Session.SessionToken = session.SessionToken
	resp.Body.Session.RefreshToken = session.RefreshToken
	resp.Body.Session.CreatedAt = session.CreatedAt.Time
	resp.Body.Session.ExpiresAt = session.ExpiresAt.Time

	return &resp, nil
}

// Logout invalidates the current session
func Logout(ctx context.Context, _ *struct{}) (*successOut, error) {
	token := ctx.Value("session_token")
	if token == nil {
		return nil, huma.Error401Unauthorized("no session token provided")
	}

	if err := db.Repository.DeleteSessionByToken(ctx, token.(string)); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, huma.Error401Unauthorized("invalid session")
		}
		return nil, huma.Error500InternalServerError("logout failed", err)
	}

	return &successOut{Body: struct {
		Success bool `json:"success" example:"true" doc:"Status of succession"`
	}(struct {
		Success bool "json:\"success\""
	}{Success: true})}, nil
}

// RefreshToken extends the session using a refresh token
func RefreshToken(ctx context.Context, in *refreshTokenIn) (*refreshTokenOut, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, huma.Error500InternalServerError("refresh failed", err)
	}
	defer tx.Rollback(ctx)

	qtx := db.Repository.WithTx(tx)

	// Get current session
	session, err := qtx.GetSessionByRefreshToken(ctx, in.Body.RefreshToken)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, huma.Error401Unauthorized("invalid refresh token")
		}
		return nil, huma.Error500InternalServerError("refresh failed", err)
	}

	// Check if refresh token is expired
	if time.Now().After(session.ExpiresAt.Time.Add(refreshDuration)) {
		if err := qtx.DeleteSession(ctx, session.ID); err != nil {
			log.SugaredLogger.Errorf("failed to delete expired session: %v", err)
		}
		return nil, huma.Error401Unauthorized("refresh token expired")
	}

	// Generate new tokens
	newSessionToken, err := generateToken()
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to generate session token", err)
	}

	newRefreshToken, err := generateToken()
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to generate refresh token", err)
	}

	// Update session
	newExpiresAt := time.Now().Add(sessionDuration)
	updatedSession, err := qtx.UpdateSession(ctx, db.UpdateSessionParams{
		ID:           session.ID,
		SessionToken: newSessionToken,
		RefreshToken: newRefreshToken,
		ExpiresAt: pgtype.Timestamp{
			Time:             newExpiresAt,
			InfinityModifier: 0,
			Valid:            true,
		},
	})
	if err != nil {
		return nil, huma.Error500InternalServerError("failed to update session", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, huma.Error500InternalServerError("failed to commit session update", err)
	}

	return &refreshTokenOut{
		Body: struct {
			SessionToken string    `json:"session_token" doc:"New session token"`
			RefreshToken string    `json:"refresh_token" doc:"New refresh token"`
			ExpiresAt    time.Time `json:"expires_at" doc:"New expiration time"`
		}(struct {
			SessionToken string    `json:"session_token"`
			RefreshToken string    `json:"refresh_token"`
			ExpiresAt    time.Time `json:"expires_at"`
		}{
			SessionToken: updatedSession.SessionToken,
			RefreshToken: updatedSession.RefreshToken,
			ExpiresAt:    updatedSession.ExpiresAt.Time,
		}),
	}, nil
}

// Helper function to generate secure random tokens
func generateToken() (string, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
