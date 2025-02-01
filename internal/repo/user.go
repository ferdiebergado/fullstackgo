//go:generate mockgen -destination=mocks/userrepo_mock.go -package=mocks . UserRepo
package repo

import (
	"context"
	"database/sql"

	"github.com/ferdiebergado/fullstackgo/internal/model"
)

type UserRepo interface {
	CreateUser(ctx context.Context, params model.User) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

const CreateUserQuery = `
INSERT into users (email, password_hash)
VALUES $1, $2
RETURNING id, email, created_at, updated_at
`

func (r *userRepo) CreateUser(ctx context.Context, params model.User) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRowContext(ctx, CreateUserQuery, params.Email, params.PasswordHash).
		Scan(&user.ID, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}

	return &user, nil
}

const FindUserByEmailQuery = `
SELECT id, password_hash
FROM users
WHERE email = $1
`

func (r *userRepo) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	if err := r.db.QueryRowContext(ctx, FindUserByEmailQuery, email).Scan(&user.ID, &user.PasswordHash); err != nil {
		return nil, err
	}

	return &user, nil
}
