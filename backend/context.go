// server/context.go

package server

import (
	"context"
	"fmt"
	"server/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Context struct {
	DB *gorm.DB
}

// NewContext initializes a new application context.
func NewContext() (*Context, error) {
	db, err := db.InitializeDB()
	if err != nil {
		return nil, err
	}

	return &Context{
		DB: db,
	}, nil
}

func WithDB(ctx context.Context, db *gorm.DB) context.Context {
	return context.WithValue(ctx, "DB", db)
}
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("ServerContextKey")
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
