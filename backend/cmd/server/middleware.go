package main

import (
	"server"
	services "server/cmd/services/auth"

	"github.com/gin-gonic/gin"
)

// func authMiddleware(authTokenService types.AuthTokenService) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx := c.Request.Context()

// 		log.Printf("token: %v", c.Request)
// 		token, err := authTokenService.ParseTokenFromRequest(ctx, c.Request)
// 		log.Printf("token: %v", token)
// 		if err != nil {
// 			c.Next()
// 			return
// 		}

// 		ctx = server.PutUserIDIntoContext(ctx, token.Sub)
// 		c.Request = c.Request.WithContext(ctx)

//			c.Next()
//		}
//	}
func authMiddleware(authTokenService services.IAuthTokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		token, err := authTokenService.ParseTokenFromRequest(ctx, c.Request)
		if err != nil {
			c.Next()
			return
		}

		// Log the parsed token
		ctx = server.PutUserIDIntoContext(ctx, token.Sub)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
