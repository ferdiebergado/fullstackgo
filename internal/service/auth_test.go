package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"

	dbmocks "github.com/ferdiebergado/fullstackgo/internal/db/mocks"
	"github.com/ferdiebergado/fullstackgo/internal/service/mocks"
)

const (
	testID         = "1"
	testEmail      = "abc@example.com"
	testPassword   = "test"
	hashedPassword = "hashed"
)

var signUpParams = model.UserSignUpParams{
	Email:           testEmail,
	Password:        testPassword,
	PasswordConfirm: testPassword,
}

func TestSignUpUserSuccess(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	ctrl := gomock.NewController(t)
	signUpParams.Password = hashedPassword

	mockRepo := dbmocks.NewMockAuthenticator(ctrl)
	mockRepo.EXPECT().SignUpUser(ctx, signUpParams).Return(&model.User{
		ID:           testID,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockValidator := mocks.NewMockValidator(ctrl)
	mockValidator.EXPECT().Struct(signUpParams).Return(nil)

	mockHasher := mocks.NewMockHasher(ctrl)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)

	service := service.NewAuthService(mockRepo, mockValidator, mockHasher)
	user, err := service.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, testID, user.ID, "ID should match")
	assert.Equal(t, signUpParams.Email, user.Email, "Emails should match")
	assert.Equal(t, now, user.CreatedAt.UTC(), "CreatedAt should match now")
	assert.Equal(t, now, user.UpdatedAt.UTC(), "UpdatedAt should match now")
}

var ErrInvalidInput = errors.New("invalid input")

func TestSignUpUserInvalidInput(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockRepo := dbmocks.NewMockAuthenticator(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	userService := service.NewAuthService(mockRepo, mockValidator, mockHasher)

	tests := []struct {
		name     string
		expected error
		given    model.UserSignUpParams
	}{
		{"should fail when email is empty", ErrInvalidInput, model.UserSignUpParams{
			Email:           "",
			Password:        testPassword,
			PasswordConfirm: testPassword,
		}},
		{"should fail when email is invalid", ErrInvalidInput, model.UserSignUpParams{
			Email:           "abcd",
			Password:        testPassword,
			PasswordConfirm: testPassword,
		}},
		{"should fail when password is empty", ErrInvalidInput, model.UserSignUpParams{
			Email:           testEmail,
			Password:        "",
			PasswordConfirm: "otherpassword",
		}},
		{"should fail when passwords do not match", ErrInvalidInput, model.UserSignUpParams{
			Email:           testEmail,
			Password:        testPassword,
			PasswordConfirm: "otherpassword",
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.EXPECT().SignUpUser(gomock.Any(), gomock.Any()).Times(0)
			mockValidator.EXPECT().Struct(tt.given).Return(ErrInvalidInput)
			mockHasher.EXPECT().Hash(gomock.Any()).Times(0)

			_, err := userService.SignUpUser(ctx, tt.given)
			assert.Equal(t, err, tt.expected, "signup should return an error")
		})
	}
}

func TestSignUpUserPasswordHashed(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()
	ctrl := gomock.NewController(t)

	mockRepo := dbmocks.NewMockAuthenticator(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	mockRepo.EXPECT().SignUpUser(ctx, signUpParams).Return(&model.User{
		ID:           testID,
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockValidator.EXPECT().Struct(signUpParams).Return(nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)
	userService := service.NewAuthService(mockRepo, mockValidator, mockHasher)
	user, err := userService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, hashedPassword, user.PasswordHash, "password hash must match")
}

func TestSignInUserSuccess(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	mockRepo := dbmocks.NewMockAuthenticator(ctrl)

	signInParams := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}

	mockRepo.EXPECT().SignInUser(ctx, signInParams).Return(hashedPassword, nil)
	mockValidator := mocks.NewMockValidator(ctrl)
	mockValidator.EXPECT().Struct(signInParams).Return(nil)

	mockHasher := mocks.NewMockHasher(ctrl)
	mockHasher.EXPECT().Verify(signInParams.Password, hashedPassword).Return(true, nil)

	service := service.NewAuthService(mockRepo, mockValidator, mockHasher)
	err := service.SignInUser(ctx, signInParams)

	assert.NoError(t, err, "signin should not return an error")
}
