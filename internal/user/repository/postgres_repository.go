package repository

import (
	"context"
	"errors"
	"github.com/Niromash/niromash-api/internal/postgres"
	"github.com/Niromash/niromash-api/internal/user/entity"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type PostgresUserRepository struct {
	client *postgres.Client
}

func NewPostgresUserRepository(client *postgres.Client) UserRepository {
	return &PostgresUserRepository{client: client}
}

type postgresUser struct {
	Id          uuid.UUID
	FirstName   string
	LastName    string
	Email       string
	Password    string
	Permissions pgtype.Array[string]
}

func (p PostgresUserRepository) AddUser(ctx context.Context, user *entity.User) error {
	query := "INSERT INTO users (first_name, last_name, email, password, permissions) VALUES (@first_name, @last_name, @email, @password, @permissions)"
	args := pgx.NamedArgs{
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"email":       user.Email,
		"password":    user.Password,
		"permissions": user.Permissions,
	}

	_, err := p.client.Client.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func (p PostgresUserRepository) GetUserById(ctx context.Context, userId entity.UserId) (user *entity.User, err error) {
	query := "SELECT id, first_name, last_name, email, password, permissions FROM users WHERE id = $1"
	row := p.client.Client.QueryRow(ctx, query, userId)

	var pgUser postgresUser
	err = row.Scan(&pgUser.Id, &pgUser.FirstName, &pgUser.LastName, &pgUser.Email, &pgUser.Password, &pgUser.Permissions)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotFoundErr
		}
		return nil, err
	}

	return p.Entity(&pgUser), nil
}

func (p PostgresUserRepository) GetUserByEmail(ctx context.Context, email string, options ...GetUserOptionFunc) (user *entity.User, err error) {
	opts := getUserOptions{}
	for _, option := range options {
		option(&opts)
	}

	var query string
	if opts.withPassword {
		query = "SELECT id, first_name, last_name, email, password, permissions FROM users WHERE lower(email) = lower($1)"
	} else {
		query = "SELECT id, first_name, last_name, email, permissions FROM users WHERE lower(email) = lower($1)"
	}

	row := p.client.Client.QueryRow(ctx, query, email)

	var pgUser postgresUser
	if opts.withPassword {
		err = row.Scan(&pgUser.Id, &pgUser.FirstName, &pgUser.LastName, &pgUser.Email, &pgUser.Password, &pgUser.Permissions)
	} else {
		err = row.Scan(&pgUser.Id, &pgUser.FirstName, &pgUser.LastName, &pgUser.Email, &pgUser.Permissions)
	}
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, UserNotFoundErr
		}
		return nil, err
	}

	return p.Entity(&pgUser), nil
}

func (p PostgresUserRepository) ListUsers(ctx context.Context) (users []*entity.User, err error) {
	query := "SELECT id, first_name, last_name, email, password, permissions FROM users"
	rows, err := p.client.Client.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pgUser postgresUser
		err = rows.Scan(&pgUser.Id, &pgUser.FirstName, &pgUser.LastName, &pgUser.Email, &pgUser.Password, &pgUser.Permissions)
		if err != nil {
			return nil, err
		}

		user := &entity.User{
			Id:          entity.UserId(pgUser.Id),
			FirstName:   pgUser.FirstName,
			LastName:    pgUser.LastName,
			Email:       pgUser.Email,
			Password:    pgUser.Password,
			Permissions: pgUser.Permissions.Elements,
		}

		users = append(users, user)
	}

	return users, nil
}

func (p PostgresUserRepository) RemoveUser(ctx context.Context, userId entity.UserId) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := p.client.Client.Exec(ctx, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserNotFoundErr
		}
		return err
	}

	return nil
}

func (p PostgresUserRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	query := "UPDATE users SET first_name = @first_name, last_name = @last_name, email = @email, password = @password, permissions = @permissions WHERE id = @id"
	args := pgx.NamedArgs{
		"id":          user.Id,
		"first_name":  user.FirstName,
		"last_name":   user.LastName,
		"email":       user.Email,
		"password":    user.Password,
		"permissions": user.Permissions,
	}

	_, err := p.client.Client.Exec(ctx, query, args)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return UserNotFoundErr
		}
		return err
	}

	return nil
}

func (p PostgresUserRepository) Entity(postgresUser *postgresUser) *entity.User {
	return &entity.User{
		Id:          entity.UserId(postgresUser.Id),
		FirstName:   postgresUser.FirstName,
		LastName:    postgresUser.LastName,
		Email:       postgresUser.Email,
		Password:    postgresUser.Password,
		Permissions: postgresUser.Permissions.Elements,
	}
}
