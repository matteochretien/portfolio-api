package login

import (
	"context"
	"errors"
	"fmt"
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/Niromash/niromash-api/internal/user/usecase"
	"github.com/Niromash/niromash-api/internal/user/usecase/token"
	"golang.org/x/crypto/bcrypt"
)

type LoginUserUseCase struct {
	getUserByEmailUseCase *usecase.GetUserByEmailUseCase
	tokenStrategy         token.TokenStrategy
}

var InvalidUsernameOrPasswordErr = errors.New("invalid username or password")
var UserNotAllowedErr = errors.New("user is not allowed to login")

func NewLoginUserUseCase(getUserByEmailUseCase *usecase.GetUserByEmailUseCase, tokenStrategy token.TokenStrategy) *LoginUserUseCase {
	return &LoginUserUseCase{getUserByEmailUseCase: getUserByEmailUseCase, tokenStrategy: tokenStrategy}
}

func (l *LoginUserUseCase) Execute(ctx context.Context, request LoginUserRequest) (*LoginUserResult, error) {
	user, err := l.getUserByEmailUseCase.Execute(ctx, request.Email, usecase.WithUserPassword())
	if err != nil {
		return nil, err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)) != nil {
		return nil, InvalidUsernameOrPasswordErr
	}

	accessToken, err := l.tokenStrategy.GenerateAccessToken(map[string]any{
		"email": request.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	refreshToken, err := l.tokenStrategy.GenerateRefreshToken(map[string]any{
		"email": request.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &LoginUserResult{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

type LoginUserRequest struct {
	Email    string
	Password string
}

type LoginUserResult struct {
	AccessToken  string
	RefreshToken string
	User         *entity.User
}
