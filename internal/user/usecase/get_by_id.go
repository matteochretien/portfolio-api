package usecase

import (
	"context"
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/Niromash/niromash-api/internal/user/repository"
)

type GetUserByIdUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUserByIdUseCase(userRepository repository.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{userRepository: userRepository}
}

func (u *GetUserByIdUseCase) Execute(ctx context.Context, userId entity.UserId) (*entity.User, error) {
	return u.userRepository.GetUserById(ctx, userId)
}
