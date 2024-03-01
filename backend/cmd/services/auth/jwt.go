package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"server"
	"server/config"
	"server/pkg/model"
	"strings"
	"time"

	jwtGo "github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var signatureType = jwtGo.SigningMethodHS256

var now = time.Now

type TokenService struct {
	Conf *config.Config
}

func NewTokenService(conf *config.Config) *TokenService {
	return &TokenService{
		Conf: conf,
	}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (*CustomClaims, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return nil, server.ErrInvalidAccessToken
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return nil, server.ErrInvalidAccessToken
	}

	tokenString := tokenParts[1]

	claims, err := ValidateToken(ctx, tokenString)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, server.ErrInvalidAccessToken
	}

	return claims, nil
}

func (ts *TokenService) CreateRefreshToken(ctx context.Context, user *model.User) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwtGo.RegisteredClaims{
			ExpiresAt: jwtGo.NewNumericDate(time.Now().Add(RefreshTokenLifetime)),
			Issuer:    "Feastly.com",
			IssuedAt:  jwtGo.NewNumericDate(time.Now()),
			ID:        uuid.New().String(), // TODO: Add this to the db
		},
		UserID: user.ID,
	}

	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign jwt: %v", err)
	}

	return tokenString, nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user *model.User) (string, error) {
	claims := CustomClaims{
		RegisteredClaims: jwtGo.RegisteredClaims{
			ExpiresAt: jwtGo.NewNumericDate(time.Now().Add(AccessTokenLifetime)),
			Issuer:    "Feastly.com",
			IssuedAt:  jwtGo.NewNumericDate(time.Now()),
			ID:        uuid.New().String(), // TODO: Add this to the db
		},
		UserID: user.ID,
	}

	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(ctx context.Context, tokenString string) (*CustomClaims, error) {
	claims := &CustomClaims{}

	secret := []byte(os.Getenv("JWT_SECRET"))

	token, err := jwtGo.ParseWithClaims(tokenString, claims, func(token *jwtGo.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwtGo.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
