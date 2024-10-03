package middlewares

import (
	"github.com/Niromash/niromash-api/internal/user/usecase"
	"github.com/Niromash/niromash-api/internal/user/usecase/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

type ApiKeyMiddleware struct {
	getUserByEmail *usecase.GetUserByEmailUseCase
	tokenStrategy  token.TokenStrategy
	apiKey         string
}

func NewApiKeyMiddleware(getUserByEmail *usecase.GetUserByEmailUseCase, tokenStrategy token.TokenStrategy) *ApiKeyMiddleware {
	return &ApiKeyMiddleware{getUserByEmail: getUserByEmail, tokenStrategy: tokenStrategy, apiKey: os.Getenv("API_KEY")}
}

func (a *ApiKeyMiddleware) Execute(c *gin.Context) {
	extractedToken := a.extractToken(c)
	if extractedToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token required"})
		return
	}

	if extractedToken != a.apiKey {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token"})
		return
	}

	c.Next()
}

func (a *ApiKeyMiddleware) extractToken(c *gin.Context) string {
	queryToken := c.Query("apiKey")
	if queryToken != "" {
		return queryToken
	}
	return c.Request.Header.Get("ApiKey")
}
