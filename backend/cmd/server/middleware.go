package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"server"
	services "server/cmd/services/auth"

	"github.com/gin-gonic/gin"
)

func authMiddleware(authTokenService services.IAuthTokenService, allowedList map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		operationName := getOperationName(c)
		if allowedList[operationName] {
			c.Next()
			return
		}

		token, err := authTokenService.ParseTokenFromRequest(ctx, c.Request)

		if err != nil {
			if os.Getenv("ENV") == "production" {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			}
			return
		}

		ctx = server.PutUserIDIntoContext(ctx, token.UserID)

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
func getOperationName(c *gin.Context) string {
	var requestBody map[string]interface{}

	if c.Request.ContentLength == 0 {
		return ""
	}

	body, err := c.GetRawData()
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		return ""
	}

	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Printf("Error decoding request body: %v", err)
		return ""
	}

	if operationName, ok := requestBody["operationName"].(string); ok {
		return operationName
	}

	return ""
}
