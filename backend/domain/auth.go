package domain

import (
	"context"
	gql_model "server/graph/model"
	"server/types"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	AuthTokenService types.AuthTokenService
	UserRepo         types.UserRepo
}

func NewAuthService(ur types.UserRepo, ats types.AuthTokenService) *AuthService {
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
		return nil, types.ErrBadCredentials
	}

	token, err := as.AuthTokenService.CreateAccessToken(ctx, user)
	if err != nil {
		return nil, types.ErrGenAccessToken
	}

	return &gql_model.AuthPayload{
		Token: &token,
		User:  user,
	}, nil
}
