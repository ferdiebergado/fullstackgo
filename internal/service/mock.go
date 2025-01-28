package service

import (
	"context"

	"github.com/ferdiebergado/fullstackgo/internal/model"
)

type MockUserService struct {
	CreateUserFn func(ctx context.Context, params model.UserCreateParams) (*model.User, error)
}

var _ UserService = (*MockUserService)(nil)

func (m *MockUserService) CreateUser(ctx context.Context, params model.UserCreateParams) (*model.User, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, params)
	}

	return nil, nil
}
