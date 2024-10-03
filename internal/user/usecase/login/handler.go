package login

import (
	"errors"
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type LoginHandler struct {
	loginUserUseCase *LoginUserUseCase
}

func NewLoginHandler(loginUserUseCase *LoginUserUseCase) *LoginHandler {
	return &LoginHandler{loginUserUseCase: loginUserUseCase}
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=3"`
}

type LoginUserResponse struct {
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	LastName    string   `json:"lastName"`
	Permissions []string `json:"permissions"`
}

func (l LoginHandler) Handle(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		if errors.Is(err, io.EOF) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": "request body is empty",
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response, err := l.loginUserUseCase.Execute(c, LoginUserRequest{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  response.AccessToken,
		"refreshToken": response.RefreshToken,
		"user":         l.LoginUserResponse(response.User),
	})
}

func (l LoginHandler) LoginUserResponse(user *entity.User) LoginUserResponse {
	return LoginUserResponse{
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Permissions: user.Permissions,
	}
}
