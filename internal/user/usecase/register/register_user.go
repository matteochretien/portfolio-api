package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/Niromash/niromash-api/internal/user/repository"
	"github.com/Niromash/niromash-api/internal/user/usecase"
	"github.com/Niromash/niromash-api/internal/user/usecase/crypt"
)

var (
	UserAlreadyExistsErr = errors.New("an account with this email already exists!")
)

type RegisterUserUseCase struct {
	existsUserUseCase *usecase.ExistsUserUseCase
	cryptStrategy     crypt.CryptStrategy
	repository        repository.UserRepository
}

func NewRegisterUserUseCase(existsUserUseCase *usecase.ExistsUserUseCase, cryptStrategy crypt.CryptStrategy, repository repository.UserRepository) *RegisterUserUseCase {
	return &RegisterUserUseCase{existsUserUseCase: existsUserUseCase, cryptStrategy: cryptStrategy, repository: repository}
}

func (r *RegisterUserUseCase) Execute(ctx context.Context, request RegisterUserRequest) error {
	exists, err := r.existsUserUseCase.Execute(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if exists {
		return UserAlreadyExistsErr
	}

	hash, err := r.cryptStrategy.Hash(request.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	request.Password = hash

	err = r.repository.AddUser(ctx, r.EntityUser(request))
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	return nil
}

func (r *RegisterUserUseCase) EntityUser(user RegisterUserRequest) *entity.User {
	return &entity.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
}

type RegisterUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}
