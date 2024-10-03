package token

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RefreshTokenHandler struct {
	useCase *RefreshAccessTokenUseCase
}

func NewRefreshTokenHandler(useCase *RefreshAccessTokenUseCase) *RefreshTokenHandler {
	return &RefreshTokenHandler{useCase: useCase}
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func (g *RefreshTokenHandler) Handle(c *gin.Context) {
	var body RefreshTokenRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Unable to parse body.",
			"error":   err.Error(),
		})
		return
	}

	accessToken, err := g.useCase.Execute(body.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "An error occurred while trying to refresh the token.",
			"error":   err.Error(),
		})
		return
	}

	c.String(http.StatusOK, accessToken)
}
