package server

import (
	"context"
	"fmt"
	"server/types"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}
type contextKey string

var (
	contextAuthIDKey contextKey = "currentUserId"
	ServerContextKey contextKey = "ServerContextKey"
)

// NewContext initializes a new application context.
func NewContext(db *gorm.DB) (*Context, error) {
	return &Context{
		DB: db,
	}, nil
}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "DB", db)
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
	if ctx.Value(contextAuthIDKey) == nil {
		return "", types.ErrNoUserIDInContext
	}

	userID, ok := ctx.Value(contextAuthIDKey).(string)
	if !ok {
		return "", types.ErrNoUserIDInContext
	}

	return userID, nil
}

func PutUserIDIntoContext(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, contextAuthIDKey, id)
}
