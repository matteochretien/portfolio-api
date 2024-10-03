package usecase

import (
	"context"
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/Niromash/niromash-api/internal/user/repository"
)

type GetUserByEmailUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUserByEmailUseCase(userRepository repository.UserRepository) *GetUserByEmailUseCase {
	return &GetUserByEmailUseCase{userRepository: userRepository}
}

type GetUserByEmailOptions struct {
	withUserPassword bool
}

type GetUserByEmailOptionFunc func(options *GetUserByEmailOptions)

func (u *GetUserByEmailUseCase) Execute(ctx context.Context, userEmail string, options ...GetUserByEmailOptionFunc) (*entity.User, error) {
	getUserByEmailOptions := &GetUserByEmailOptions{}
	for _, option := range options {
		option(getUserByEmailOptions)
	}

	if getUserByEmailOptions.withUserPassword {
		return u.userRepository.GetUserByEmail(ctx, userEmail, repository.WithUserPassword())
	}

	return u.userRepository.GetUserByEmail(ctx, userEmail)
}

func WithUserPassword() GetUserByEmailOptionFunc {
	return func(options *GetUserByEmailOptions) {
		options.withUserPassword = true
	}
}
