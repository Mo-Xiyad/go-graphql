package types

import (
	"context"
	"net/http"
	gql_model "server/graph/model"
	"server/pkg/model"
)

type AuthService interface {
	Login(ctx context.Context, input gql_model.LoginInput) (*gql_model.AuthPayload, error)
}

type AuthTokenService interface {
	CreateAccessToken(ctx context.Context, user *model.User) (string, error)
	// CreateRefreshToken(ctx context.Context, user User, tokenID string) (string, error)
	// ParseToken(ctx context.Context, payload string) (AuthToken, error)
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error)
}

type AuthToken struct {
	ID  string
	Sub string
}
