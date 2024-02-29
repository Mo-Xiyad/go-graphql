package services

import (
	"context"
	"net/http"
	"server"
	user "server/cmd/services/user"
	gql_model "server/graph/model"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthTokenService IAuthTokenService
	UserRepo         user.IUserRepo
}

func NewAuthService(ur user.IUserRepo, ats IAuthTokenService) *AuthService {
	return &AuthService{
		AuthTokenService: ats,
		UserRepo:         ur,
	}
}

func (as *AuthService) Login(ctx context.Context, input gql_model.LoginInput) (*gql_model.AuthPayload, error) {
	user, err := as.UserRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, server.ErrBadCredentials
	}

	token, err := as.AuthTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, server.ErrGenAccessToken
	}

	ca, err := server.GetCookieAccessFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	http.SetCookie(ca.Writer, &http.Cookie{
		Name:     string(server.CookieAccessTokenKey),
		Value:    token,
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Now().Add(AccessTokenLifetime),
	})

	return &gql_model.AuthPayload{
		Token: &token,
		User:  user,
	}, nil
}
