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

func TestAuthService_SignUpUser_Success(t *testing.T) {
	mockRepo, mockValidator, mockHasher, authService := setupMocks(t)
	ctx := context.Background()
	now := time.Now().UTC()
	signUpParams := newSignUpParams()
	signUpParamsHashed := newSignUpParams()
	signUpParamsHashed.Password = hashedPassword

	mockRepo.EXPECT().SignUpUser(ctx, signUpParamsHashed).Return(&model.User{
		ID:           testID,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockValidator.EXPECT().Struct(signUpParams).Return(nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)

	user, err := authService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, testID, user.ID, "ID should match")
	assert.Equal(t, signUpParams.Email, user.Email, "Emails should match")
	assert.Equal(t, now, user.CreatedAt.UTC(), "CreatedAt should match now")
	assert.Equal(t, now, user.UpdatedAt.UTC(), "UpdatedAt should match now")
}

var ErrInvalidInput = errors.New("invalid input")

func TestAuthService_SignUpUser_InvalidInput(t *testing.T) {
	mockRepo, mockValidator, mockHasher, authService := setupMocks(t)
	ctx := context.Background()

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

			_, err := authService.SignUpUser(ctx, tt.given)
			assert.Equal(t, err, tt.expected, "signup should return an error")
		})
	}
}

func TestAuthService_SignUpUser_PasswordHashed(t *testing.T) {
	mockRepo, mockValidator, mockHasher, authService := setupMocks(t)
	ctx := context.Background()
	now := time.Now().UTC()

	signUpParams := newSignUpParams()
	signUpParamsHashed := newSignUpParams()

	signUpParamsHashed.Password = hashedPassword
	mockRepo.EXPECT().SignUpUser(ctx, signUpParamsHashed).Return(&model.User{
		ID:           testID,
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockValidator.EXPECT().Struct(signUpParams).Return(nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)
	user, err := authService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, hashedPassword, user.PasswordHash, "password hash must match")
}

func TestAuthService_SignInUser_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo, mockValidator, mockHasher, authService := setupMocks(t)

	signInParams := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}

	mockRepo.EXPECT().SignInUser(ctx, signInParams).Return(&model.User{
		ID:           testID,
		PasswordHash: hashedPassword,
	}, nil)
	mockValidator.EXPECT().Struct(signInParams).Return(nil)
	mockHasher.EXPECT().Verify(signInParams.Password, hashedPassword).Return(true, nil)

	id, err := authService.SignInUser(ctx, signInParams)

	assert.NoError(t, err, "signin should not return an error")
	assert.Equal(t, testID, id, "ID must match")
}

func setupMocks(t *testing.T) (*dbmocks.MockAuthenticator, *mocks.MockValidator, *mocks.MockHasher, service.AuthService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := dbmocks.NewMockAuthenticator(ctrl)
	mockValidator := mocks.NewMockValidator(ctrl)
	mockHasher := mocks.NewMockHasher(ctrl)
	userService := service.NewAuthService(mockRepo, mockValidator, mockHasher)

	return mockRepo, mockValidator, mockHasher, userService
}

func newSignUpParams() model.UserSignUpParams {
	return model.UserSignUpParams{
		Email:           testEmail,
		Password:        testPassword,
		PasswordConfirm: testPassword,
	}
}
