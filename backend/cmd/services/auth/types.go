package services

import (
	"context"
	"net/http"
	gql_model "server/graph/model"
	"server/pkg/model"
	"time"

	jwtGo "github.com/golang-jwt/jwt/v5"
)

type IAuthService interface {
	Login(ctx context.Context, input gql_model.LoginInput) (*gql_model.AuthPayload, error)
}

// TODO: add refresh token
// used in the middleware right now
type IAuthTokenService interface {
	CreateAccessToken(ctx context.Context, user *model.User) (string, error)
	// CreateRefreshToken(ctx context.Context, user User, tokenID string) (string, error)
	// ParseToken(ctx context.Context, payload string) (AuthToken, error)
	ParseTokenFromRequest(ctx context.Context, r *http.Request) (*CustomClaims, error)
	// ValidateToken(ctx context.Context, token string) (bool, error)
}

type CustomClaims struct {
	jwtGo.RegisteredClaims
	UserID uint64
}

var (
	AccessTokenLifetime  = time.Minute * 15   // 15 minutes
	RefreshTokenLifetime = time.Hour * 24 * 7 // 1 week
)

type RefreshToken struct {
	ID         string
	Name       string
	UserID     string
	LastUsedAt time.Time
	ExpiredAt  time.Time
	CreatedAt  time.Time
}

type CreateRefreshTokenParams struct {
	Sub  string
	Name string
}

type IRefreshTokenRepo interface {
	Create(ctx context.Context, params CreateRefreshTokenParams) (RefreshToken, error)
	GetByID(ctx context.Context, id string) (RefreshToken, error)
}
