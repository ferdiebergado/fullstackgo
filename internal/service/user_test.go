package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

type mockUserRepo struct {
	CreateUserFn func(ctx context.Context, params model.UserCreateParams) (*model.User, error)
}

func (m *mockUserRepo) CreateUser(ctx context.Context, params model.UserCreateParams) (*model.User, error) {
	if m.CreateUserFn != nil {
		return m.CreateUserFn(ctx, params)
	}

	return nil, nil
}

func TestCreateUserService(t *testing.T) {
	createParams := model.UserCreateParams{
		Email:      testEmail,
		Password:   testPassword,
		AuthMethod: authMethod,
	}

	now := time.Now().UTC()

	service := service.NewUserService(&mockUserRepo{
		CreateUserFn: func(ctx context.Context, params model.UserCreateParams) (*model.User, error) {
			if params != createParams {
				t.Errorf("want: %s; got: %s", createParams, params)
			}

			return &model.User{
				ID:         testID,
				Email:      testEmail,
				AuthMethod: authMethod,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
	})

	user, err := service.CreateUser(context.Background(), createParams)

	if err != nil {
		t.Errorf("wanted no error, but got: %v", err)
	}

	if user.ID != testID {
		t.Errorf("want: %s; got: %s", testID, user.ID)
	}

	if user.Email != testEmail {
		t.Errorf("want: %s; got %s", testEmail, user.Email)
	}

	if user.AuthMethod != authMethod {
		t.Errorf("want: %s; but got: %s", authMethod, user.AuthMethod)
	}

	if user.CreatedAt.UTC() != now {
		t.Errorf("want: %v; got: %v", now, user.CreatedAt)
	}

	if user.UpdatedAt.UTC() != now {
		t.Errorf("want: %v; got: %v", now, user.UpdatedAt)
	}
}
