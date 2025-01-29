package db

import (
	"context"

	"github.com/ferdiebergado/fullstackgo/internal/model"
)

//go:generate mockgen -destination=mocks/authenticator_mock.go -package=mocks . Authenticator
type Authenticator interface {
	SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error)
}

type authRepo struct {
	db Querier
}

func NewUserRepo(db Querier) Authenticator {
	return &authRepo{
		db: db,
	}
}

const SignUpUserQuery = `
INSERT into users (email, password_hash)
VALUES $1, $2
RETURNING id, email, created_at, updated_at
`

func (r *authRepo) SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRowContext(ctx, SignUpUserQuery, params.Email, params.Password).Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}
