package repository

import (
	"context"
	"errors"
	"github.com/Niromash/niromash-api/internal/user/entity"
)

var UserNotFoundErr = errors.New("user not found")

type UserRepository interface {
	AddUser(ctx context.Context, user *entity.User) error
	GetUserById(ctx context.Context, userId entity.UserId) (user *entity.User, err error)
	GetUserByEmail(ctx context.Context, userEmail string, options ...GetUserOptionFunc) (user *entity.User, err error)
	ListUsers(ctx context.Context) (users []*entity.User, err error)
	RemoveUser(ctx context.Context, userId entity.UserId) error
	UpdateUser(ctx context.Context, user *entity.User) error
}

type getUserOptions struct {
	withPassword bool
}

type GetUserOptionFunc func(options *getUserOptions)

func WithUserPassword() GetUserOptionFunc {
	return func(options *getUserOptions) {
		options.withPassword = true
	}
}
