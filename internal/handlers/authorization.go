package handlers

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"smartDriver/internal/db"
	"time"
)

type CreateUserInput struct {
	Body struct {
		Login          string `json:"login" maxLength:"50" doc:"User login"`
		Password       string `json:"password" maxLength:"50" doc:"User password"`
		Name           string `json:"name" maxLength:"50" doc:"User name"`
		Surname        string `json:"surname" maxLength:"50" doc:"User surname"`
		Patronymic     string `json:"patronymic" maxLength:"50" doc:"User patronymic"`
		OrganizationID int32  `json:"organization_id" example:"1" doc:"Organization ID"`
	}
}

func CreateUserHandler(ctx context.Context, input *CreateUserInput) (*SuccessOutput, error) {
	resp := &SuccessOutput{}

	params := db.CreateUserParams{
		Login: input.Body.Login,
		// TODO: implement hashing
		Password:       input.Body.Password,
		Name:           pgtype.Text{input.Body.Name, true},
		Surname:        pgtype.Text{input.Body.Surname, true},
		Patronymic:     pgtype.Text{input.Body.Patronymic, true},
		OrganizationID: input.Body.OrganizationID,
	}

	queries := db.Queries{}
	_, err := queries.CreateUser(ctx, params)
	if err != nil {
		resp.Body.Success = false
		return nil, err
	}

	resp.Body.Success = true
	return resp, nil
}

type LoginInput struct {
	Body struct {
		Username string `json:"username" example:"admin" doc:"User name"`
		Password string `json:"password" example:"admin" doc:"User password"`
	}
}

type LoginOutput struct {
	Body struct {
		UserID       int32     `json:"user_id" example:"1" doc:"User ID"`
		Username     string    `json:"username" example:"admin" doc:"User name"`
		SessionToken string    `json:"session_token" example:"some hex-string" doc:"Identification Token"`
		RefreshToken string    `json:"refresh_token" example:"some hex-string" doc:"Refresh Token"`
		CreateDate   time.Time `json:"create_date" doc:"Create Date"`
		ExpiresAt    time.Time `json:"expires_at" doc:"Expires Date"`
	}
}

func LoginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	resp := &LoginOutput{}

	// todo: implement validating user

	// todo: implement creating session
	// todo: sub: implement creating session and refresh tokens

	//params :=
	//
	//queries := db.Queries{}
	//_, err := queries.CreateUser(ctx, params)
	//if err != nil {
	//	resp.Body.Success = false
	//	return nil, err
	//}

	return resp, nil
}

func LogoutHandler(ctx context.Context, input *struct{}) (*SuccessOutput, error) {
	resp := &SuccessOutput{}

	queries := db.Queries{}
	err := queries.DeleteSessionByToken(ctx, ctx.Value("token").(string))
	if err != nil {
		resp.Body.Success = false
		return nil, err
	}

	resp.Body.Success = true
	return resp, nil
}

type RefreshTokenInput struct {
	Body struct {
		RefreshToken string `json:"refresh_token" example:"some hex-string" doc:"Refresh Token"`
	}
}

// RefreshTokenHandler refreshes token expiry,
// user needs to be authenticated to be able to refresh session token
func RefreshTokenHandler(ctx context.Context, input *RefreshTokenInput) (*SuccessOutput, error) {
	resp := &SuccessOutput{}

	// todo: implement token refreshing, need to update sql code that sets expiry date
	queries := db.Queries{}
	err := queries.UpdateSessionExpiry(ctx, input.Body.RefreshToken)
	if err != nil {
		resp.Body.Success = false
		return nil, err
	}

	return resp, nil
}
