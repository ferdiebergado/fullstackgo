package db

import (
	"context"

	"github.com/ferdiebergado/fullstackgo/internal/model"
)

//go:generate mockgen -destination=mocks/user_repo_mock.go -package=mocks . UserRepo
type UserRepo interface {
	CreateUser(ctx context.Context, params model.UserCreateParams) (*model.User, error)
}

type userRepo struct {
	db Querier
}

func NewUserRepo(db Querier) UserRepo {
	return &userRepo{
		db: db,
	}
}

const CreateUserQuery = `
INSERT into users (email, password_hash, auth_method)
VALUES $1, $2, $3
RETURNING id, email, auth_method, created_at, updated_at
`

func (r *userRepo) CreateUser(ctx context.Context, params model.UserCreateParams) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRowContext(ctx, CreateUserQuery, params.Email, params.Password, params.AuthMethod).Scan(&user.ID, &user.Email, &user.AuthMethod, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
