//go:generate mockgen -destination=mocks/authenticator_mock.go -package=mocks . Authenticator
package db

import (
	"context"
	"database/sql"

	"github.com/ferdiebergado/fullstackgo/internal/model"
)

type Authenticator interface {
	SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error)
	SignInUser(ctx context.Context, params model.UserSignInParams) (*model.User, error)
}

type authRepo struct {
	db *sql.DB
}

func NewAuthRepo(db *sql.DB) Authenticator {
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

const SignInUserQuery = `
SELECT id, password_hash
FROM users
WHERE email = $1
`

func (r *authRepo) SignInUser(ctx context.Context, params model.UserSignInParams) (*model.User, error) {
	var id string
	var hash string
	if err := r.db.QueryRowContext(ctx, SignInUserQuery, params.Email).Scan(&id, &hash); err != nil {
		return nil, err
	}

	return &model.User{ID: id, PasswordHash: hash}, nil
}
