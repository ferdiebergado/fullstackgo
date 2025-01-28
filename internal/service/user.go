package service

import (
	"context"

	"github.com/ferdiebergado/fullstackgo/internal/db"
	"github.com/ferdiebergado/fullstackgo/internal/model"
)

type UserService interface {
	CreateUser(ctx context.Context, params model.UserCreateParams) (*model.User, error)
}

type userService struct {
	repo db.UserRepo
}

func NewUserService(repo db.UserRepo) UserService {
	return &userService{repo: repo}
}

func (u *userService) CreateUser(ctx context.Context, params model.UserCreateParams) (*model.User, error) {
	user, err := u.repo.CreateUser(ctx, params)

	if err != nil {
		return nil, err
	}

	return user, nil
}
