package services

import (
	"context"
	"fmt"
	"log"
	"net/http"
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

func buildToken(token *jwtGo.Token) AuthToken {
	claims, ok := token.Claims.(jwtGo.MapClaims)
	log.Printf("claims: %v", claims)
	if !ok {
		return AuthToken{}
	}

	if id, ok := claims["jti"].(string); ok {
		return AuthToken{
			ID:  id,
			Sub: fmt.Sprintf("%v", claims["sub"]),
		}
	} else {
		return AuthToken{
			ID:  "",
			Sub: fmt.Sprintf("%v", claims["sub"]),
		}
	}
}

func (ts *TokenService) ParseTokenFromRequest(ctx context.Context, r *http.Request) (AuthToken, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return AuthToken{}, server.ErrInvalidAccessToken
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return AuthToken{}, server.ErrInvalidAccessToken
	}

	tokenString := tokenParts[1]

	secret := []byte(ts.Conf.JWT.Secret)

	token, err := jwtGo.Parse(tokenString, func(token *jwtGo.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		log.Printf("error: %v", err)
		return AuthToken{}, server.ErrInvalidAccessToken
	}

	return buildToken(token), nil
}

func (ts *TokenService) ParseToken(ctx context.Context, payload string) (AuthToken, error) {
	secret := []byte(ts.Conf.JWT.Secret)
	token, err := jwtGo.ParseWithClaims(payload, jwtGo.MapClaims{}, func(token *jwtGo.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return AuthToken{}, server.ErrInvalidAccessToken
	}

	return buildToken(token), nil
}

func (ts *TokenService) CreateRefreshToken(ctx context.Context, user *model.User, tokenID string) (string, error) {
	claims := jwtGo.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(RefreshTokenLifetime).Unix(),
		"iat": time.Now().Unix(),
		"jti": tokenID,
	}

	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign jwt: %v", err)
	}

	return tokenString, nil
}

func (ts *TokenService) CreateAccessToken(ctx context.Context, user *model.User) (string, error) {
	claims := jwtGo.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(AccessTokenLifetime).Unix(),
		"iat": time.Now().Unix(),
		"jti": uuid.New().String(),
	}

	token := jwtGo.NewWithClaims(signatureType, claims)
	tokenString, err := token.SignedString([]byte(ts.Conf.JWT.Secret))
	if err != nil {
		return "", fmt.Errorf("error sign jwt: %v", err)
	}

	return tokenString, nil
}
