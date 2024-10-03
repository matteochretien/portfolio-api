package middlewares

import (
	"errors"
	user_repository "github.com/Niromash/niromash-api/internal/user/repository"
	"github.com/Niromash/niromash-api/internal/user/usecase"
	"github.com/Niromash/niromash-api/internal/user/usecase/token"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type AuthMiddleware struct {
	getUserByEmail *usecase.GetUserByEmailUseCase
	tokenStrategy  token.TokenStrategy
}

func NewAuthMiddleware(getUserByEmail *usecase.GetUserByEmailUseCase, tokenStrategy token.TokenStrategy) *AuthMiddleware {
	return &AuthMiddleware{getUserByEmail: getUserByEmail, tokenStrategy: tokenStrategy}
}

func (a *AuthMiddleware) Execute(c *gin.Context) {
	extractedToken := a.extractToken(c)
	if extractedToken == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "token required"})
		return
	}

	accessTokenClaims, err := a.tokenStrategy.ValidateAccessToken(extractedToken)
	if err != nil {
		if errors.Is(err, token.UnableToValidateAccessTokenErr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unable to validate access token"})
			return
		}
		if errors.Is(err, token.AccessTokenExpiredErr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "access token expired"})
			return
		}
		if errors.Is(err, token.InvalidTokenTypeErr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid token type"})
			return
		}
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "invalid access token"})
		return
	}

	user, err := a.getUserByEmail.Execute(c, accessTokenClaims["email"].(string))
	if err != nil {
		if errors.Is(err, user_repository.UserNotFoundErr) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Unable to find user relative to this account!"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal server error"})
		return
	}

	c.Set("user", user)
}

func (a *AuthMiddleware) extractToken(c *gin.Context) string {
	queryToken := c.Query("token")
	if queryToken != "" {
		return queryToken
	}
	bearerToken := c.Request.Header.Get("Authorization")
	if len(strings.Split(bearerToken, " ")) == 2 {
		return strings.Split(bearerToken, " ")[1]
	}
	return ""
}
