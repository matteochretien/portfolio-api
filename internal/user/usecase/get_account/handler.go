package get_account

import (
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/Niromash/niromash-api/internal/user/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type GetUserHandler struct {
	useCase *usecase.GetUserByIdUseCase
}

func NewGetUserHandler(useCase *usecase.GetUserByIdUseCase) *GetUserHandler {
	return &GetUserHandler{useCase: useCase}
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	LastName    string    `json:"lastName"`
	Email       string    `json:"email"`
	Permissions []string  `json:"permissions"`
	Allowed     bool      `json:"allowed"`
}

func (g *GetUserHandler) Handle(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	c.JSON(http.StatusOK, g.Response(user.(*entity.User)))
}

func (g *GetUserHandler) Response(user *entity.User) *UserResponse {
	return &UserResponse{
		ID:          uuid.UUID(user.Id),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Permissions: user.Permissions,
	}
}
