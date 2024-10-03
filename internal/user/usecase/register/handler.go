package register

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterHandler struct {
	registerUserUseCase *RegisterUserUseCase
}

func NewRegisterHandler(useCase *RegisterUserUseCase) *RegisterHandler {
	return &RegisterHandler{registerUserUseCase: useCase}
}

type RegisterRequest struct {
	Email                string `json:"email" binding:"required"`
	FirstName            string `json:"firstName" binding:"required"`
	LastName             string `json:"lastName" binding:"required"`
	Password             string `json:"password" binding:"required,min=3"`
	PasswordConfirmation string `json:"passwordConfirmation" binding:"required,min=3"`
}

func (l RegisterHandler) Handle(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if request.Password != request.PasswordConfirmation {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password and password confirmation are not the same.",
		})
		return
	}

	err := l.registerUserUseCase.Execute(c, l.RegisterUserRequest(request))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.Status(http.StatusNoContent)
}

func (l RegisterHandler) RegisterUserRequest(request RegisterRequest) RegisterUserRequest {
	return RegisterUserRequest{
		Email:     request.Email,
		FirstName: request.FirstName,
		LastName:  request.LastName,
		Password:  request.Password,
	}
}
