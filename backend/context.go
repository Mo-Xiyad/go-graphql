package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}
type contextKey string

var (
	CurrentAuthUserId    contextKey = "currentUserId"
	ServerContextKey     contextKey = "ServerContextKey"
	dbContextKey         contextKey = "DB"
	IsLoggedIn           contextKey = "isLoggedIn"
	CookieAccessTokenKey contextKey = "cookieAccessToken"
)

// NewContext initializes a new application context.
func NewContext(db *gorm.DB) (*Context, error) {
	return &Context{
		DB: db,
	}, nil
}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, dbContextKey, db)
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ServerContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ServerContextKey)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetUserIDFromContext(ctx context.Context) (string, error) {
	if ctx.Value(CurrentAuthUserId) == nil {
		return "", ErrNoUserIDInContext
	}

	userID, ok := ctx.Value(CurrentAuthUserId).(string)
	if !ok {
		return "", ErrNoUserIDInContext
	}

	return userID, nil
}

func PutUserIDIntoContext(ctx context.Context, id uint64) context.Context {
	return context.WithValue(ctx, CurrentAuthUserId, id)
}

func SetIsLoggedIn(ctx context.Context, isAuth bool) context.Context {
	return context.WithValue(context.Background(), IsLoggedIn, isAuth)
}

func CheckIsLoggedIn(ctx context.Context) bool {
	if ctx.Value(IsLoggedIn) == nil {
		return false
	}

	auth, ok := ctx.Value(IsLoggedIn).(bool)
	if !ok {
		return false
	}

	return auth
}

type CookieAccess struct {
	Writer     http.ResponseWriter
	IsLoggedIn bool
	UserId     uint64
}

func SetCookyInCtx(ctx *gin.Context, val *CookieAccess) {
	newCtx := context.WithValue(ctx.Request.Context(), CookieAccessTokenKey, val)
	ctx.Request = ctx.Request.WithContext(newCtx)
}

func GetCookieAccessFromCtx(ctx context.Context) (*CookieAccess, error) {
	if ctx.Value(CookieAccessTokenKey) == nil {
		return nil, errors.New("no cookie access in context")
	}

	ca, ok := ctx.Value(CookieAccessTokenKey).(*CookieAccess)
	if !ok {
		return nil, errors.New("no cookie access in context")
	}

	return ca, nil
}
