package handler

import (
	"context"
	"smartDriver/internal/db"
	"smartDriver/pkg/log"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"golang.org/x/crypto/bcrypt"
)

type createUserIn struct {
	Body struct {
		Login          string `json:"login" maxLength:"50" doc:"User login"`
		Password       string `json:"password" maxLength:"50" doc:"User password"`
		Name           string `json:"name" maxLength:"50" doc:"User name"`
		Surname        string `json:"surname" maxLength:"50" doc:"User surname"`
		Patronymic     string `json:"patronymic" maxLength:"50" doc:"User patronymic"`
		OrganizationID int64  `json:"organization_id" example:"1" doc:"Organization ID"`
	}
}

type userOut struct {
	Body struct {
		ID             int64  `json:"id" example:"1" doc:"User ID"`
		Login          string `json:"login" maxLength:"50" doc:"User login"`
		Name           string `json:"name" maxLength:"50" doc:"User name"`
		Surname        string `json:"surname" maxLength:"50" doc:"User surname"`
		Patronymic     string `json:"patronymic" maxLength:"50" doc:"User patronymic"`
		OrganizationID int64  `json:"organization_id" example:"1" doc:"Organization ID"`
	}
}

type userBase struct {
	ID             int64  `json:"id" example:"1" doc:"User ID"`
	Login          string `json:"login" maxLength:"50" doc:"User login"`
	Name           string `json:"name" maxLength:"50" doc:"User name"`
	Surname        string `json:"surname" maxLength:"50" doc:"User surname"`
	Patronymic     string `json:"patronymic" maxLength:"50" doc:"User patronymic"`
	OrganizationID int64  `json:"organization_id" example:"1" doc:"Organization ID"`
}

func CreateUser(ctx context.Context, in *createUserIn) (*userOut, error) {
	hashedPassword, err := hashPassword(in.Body.Password)
	if err != nil {
		log.SugaredLogger.Errorf("failed to hash password: %v", err)
		return nil, huma.Error500InternalServerError("failed to create user", err)
	}

	user, err := db.Pool.CreateUser(ctx, db.CreateUserParams{
		Login:          in.Body.Login,
		Password:       hashedPassword,
		Name:           &in.Body.Name,
		Surname:        &in.Body.Surname,
		Patronymic:     &in.Body.Patronymic,
		OrganizationID: in.Body.OrganizationID,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to create user: %v", err)
		return nil, huma.Error500InternalServerError("failed to create user", err)
	}

	var resp userOut
	resp.Body.ID = user.ID
	resp.Body.Login = user.Login
	resp.Body.Name = *user.Name
	resp.Body.Surname = *user.Surname
	resp.Body.Patronymic = *user.Patronymic
	resp.Body.OrganizationID = user.OrganizationID

	return &resp, nil
}

func GetUser(ctx context.Context, in *idPathIn) (*userOut, error) {
	user, err := db.Pool.GetUser(ctx, in.ID)
	if err != nil {
		log.SugaredLogger.Errorf("failed to get user: %v", err)
		return nil, huma.Error500InternalServerError("failed to get user", err)
	}

	var resp userOut
	resp.Body.ID = user.ID
	resp.Body.Login = user.Login
	resp.Body.Name = *user.Name
	resp.Body.Surname = *user.Surname
	resp.Body.Patronymic = *user.Patronymic
	resp.Body.OrganizationID = user.OrganizationID

	return &resp, nil
}

type updateUserIn struct {
	ID   int64 `path:"id" json:"id" example:"1" doc:"User ID"`
	Body struct {
		Password string `json:"password" maxLength:"50" doc:"User password"`
		Name     string `json:"name" maxLength:"50" doc:"User name"`
		Surname  string `json:"surname" maxLength:"50" doc:"User surname"`
	}
}

type updateUserOut struct {
	Body struct {
		ID       int64  `json:"id" example:"1" doc:"User ID"`
		Password string `json:"password" maxLength:"50" doc:"User password"`
		Name     string `json:"name" maxLength:"50" doc:"User name"`
		Surname  string `json:"surname" maxLength:"50" doc:"User surname"`
	}
}

func UpdateUser(ctx context.Context, in *updateUserIn) (*updateUserOut, error) {
	hashedPassword, err := hashPassword(in.Body.Password)
	if err != nil {
		log.SugaredLogger.Errorf("failed to hash passowrd: %v", err)
		return nil, huma.Error500InternalServerError("failed to hash password", err)
	}
	user, err := db.Pool.UpdateUser(ctx, db.UpdateUserParams{
		ID:       in.ID,
		Name:     &in.Body.Name,
		Surname:  &in.Body.Surname,
		Password: hashedPassword,
	})
	if err != nil {
		log.SugaredLogger.Errorf("failed to update user: %v", err)
		return nil, huma.Error500InternalServerError("failed to update user", err)
	}

	var resp updateUserOut
	resp.Body.ID = user.ID
	resp.Body.Name = *user.Name
	resp.Body.Surname = *user.Surname
	resp.Body.Password = user.Password

	return &resp, nil
}

func DeleteUser(ctx context.Context, in *idPathIn) (*successOut, error) {
	if err := db.Pool.DeleteUser(ctx, in.ID); err != nil {
		log.SugaredLogger.Errorf("failed to delete user: %v", err)
		return nil, huma.Error500InternalServerError("failed to delete user", err)
	}

	var resp successOut
	resp.Body.Success = true

	return &resp, nil
}

type listUsersOut struct {
	Body struct {
		Users []userBase `json:"users" doc:"List of users"`
	}
}

func ListUsers(ctx context.Context, in *listIn) (*listUsersOut, error) {
	users, err := db.Pool.ListUsers(ctx)
	if err != nil {
		log.SugaredLogger.Errorf("failed to list users: %v", err)
		return nil, huma.Error500InternalServerError("failed to list users", err)
	}

	var resp listUsersOut
	for _, user := range users {
		resp.Body.Users = append(resp.Body.Users, userBase{
			ID:             user.ID,
			Login:          user.Login,
			Name:           *user.Name,
			Surname:        *user.Surname,
			Patronymic:     *user.Patronymic,
			OrganizationID: user.OrganizationID,
		})
	}

	return &resp, nil
}

type loginIn struct {
	Body struct {
		Username string `json:"username" example:"admin" doc:"User username"`
		Password string `json:"password" example:"admin" doc:"User password"`
	}
}

type loginOut struct {
	Body struct {
		UserID       int32     `json:"user_id" example:"1" doc:"User ID"`
		Username     string    `json:"username" example:"admin" doc:"User name"`
		SessionToken string    `json:"session_token" example:"some hex-string" doc:"Identification Token"`
		RefreshToken string    `json:"refresh_token" example:"some hex-string" doc:"Refresh Token"`
		CreateDate   time.Time `json:"create_date" doc:"Create Date"`
		ExpiresAt    time.Time `json:"expires_at" doc:"Expires Date"`
	}
}

func Login(ctx context.Context, in *loginIn) (*loginOut, error) {
	resp := &loginOut{}

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

func Logout(ctx context.Context, in *struct{}) (*successOut, error) {
	resp := &successOut{}

	queries := db.Queries{}
	err := queries.DeleteSessionByToken(ctx, ctx.Value("token").(string))
	if err != nil {
		resp.Body.Success = false
		return nil, err
	}

	resp.Body.Success = true
	return resp, nil
}

type refreshTokenIn struct {
	Body struct {
		RefreshToken string `json:"refresh_token" example:"some hex-string" doc:"Refresh Token"`
	}
}

// RefreshTokenHandler refreshes token expiry,
// user needs to be authenticated to be able to refresh session token
func RefreshToken(ctx context.Context, in *refreshTokenIn) (*successOut, error) {
	resp := &successOut{}

	// todo: implement token refreshing, need to update sql code that sets expiry date
	queries := db.Queries{}
	err := queries.UpdateSessionExpiry(ctx, in.Body.RefreshToken)
	if err != nil {
		resp.Body.Success = false
		return nil, err
	}

	return resp, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
