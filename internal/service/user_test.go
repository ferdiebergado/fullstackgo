package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
	"go.uber.org/mock/gomock"

	dbmocks "github.com/ferdiebergado/fullstackgo/internal/db/mocks"
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

func TestCreateUserService(t *testing.T) {
	createParams := model.UserCreateParams{
		Email:      testEmail,
		Password:   testPassword,
		AuthMethod: authMethod,
	}

	ctx := context.Background()
	now := time.Now().UTC()

	ctrl := gomock.NewController(t)
	mockRepo := dbmocks.NewMockUserRepo(ctrl)
	mockRepo.EXPECT().CreateUser(ctx, createParams).Return(&model.User{
		ID:         testID,
		Email:      testEmail,
		AuthMethod: authMethod,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil)
	service := service.NewUserService(mockRepo)

	user, err := service.CreateUser(ctx, createParams)

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
