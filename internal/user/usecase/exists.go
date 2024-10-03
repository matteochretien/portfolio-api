package usecase

import (
	"context"
	"errors"
	"github.com/Niromash/niromash-api/internal/user/repository"
)

type ExistsUserUseCase struct {
	userRepository repository.UserRepository
}

func NewExistsUserUseCase(userRepository repository.UserRepository) *ExistsUserUseCase {
	return &ExistsUserUseCase{userRepository: userRepository}
}

func (u *ExistsUserUseCase) Execute(ctx context.Context, userEmail string) (bool, error) {
	_, err := u.userRepository.GetUserByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, repository.UserNotFoundErr) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
