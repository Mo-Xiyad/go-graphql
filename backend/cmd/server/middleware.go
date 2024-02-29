package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"server"
	services "server/cmd/services/auth"

	"github.com/gin-gonic/gin"
)

// Add cookies to the Request
func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cookieA := server.CookieAccess{
			Writer:     ctx.Writer,
			IsLoggedIn: false,
			UserId:     0,
		}

		server.SetCookyInCtx(ctx, &cookieA)

		c, err := ctx.Request.Cookie(string(server.CookieAccessTokenKey))
		if err != nil {
			// If there's an error fetching the cookie, log it and proceed
			log.Printf("Error fetching 'CookieAccessTokenKey' cookie: %v", err)
			ctx.Next()
			return
		}

		// Proceed with token parsing only if the cookie is successfully fetched
		rawToken := c.Value
		userId, err := services.ValidateToken(ctx, rawToken)

		if err != nil {
			// If there's an error parsing the token, log it and ensure user is not logged in
			log.Printf("Error parsing token: %v", err)
		} else {
			// If token is successfully parsed, mark user as logged in and set UserId
			cookieA.IsLoggedIn = true
			cookieA.UserId = userId.UserID
		}

		ctx.Next()
	}
}

func authMiddleware(authTokenService services.IAuthTokenService, allowedList map[string]bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		ctx = server.SetIsLoggedIn(ctx, false)

		operationName := getOperationName(c)
		if allowedList[operationName] {
			c.Next()
			return
		}

		token, err := authTokenService.ParseTokenFromRequest(ctx, c.Request)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		ctx = server.SetIsLoggedIn(ctx, true)
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
