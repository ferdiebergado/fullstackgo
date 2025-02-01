package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
	"github.com/stretchr/testify/assert"

	"go.uber.org/mock/gomock"

	secMocks "github.com/ferdiebergado/fullstackgo/internal/pkg/security/mocks"
	repoMocks "github.com/ferdiebergado/fullstackgo/internal/repo/mocks"
)

const (
	testID         = "1"
	testEmail      = "abc@example.com"
	testPassword   = "test"
	hashedPassword = "hashed"
)

func TestAuthService_SignUpUser_Success(t *testing.T) {
	mockRepo, mockHasher, authService := setupMocks(t)
	ctx := context.Background()
	now := time.Now().UTC()
	signUpParams := newSignUpParams()
	createParams := &model.User{
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
	}

	mockRepo.EXPECT().FindUserByEmail(ctx, signUpParams.Email).Return(nil, nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)
	mockRepo.EXPECT().CreateUser(ctx, *createParams).Return(&model.User{
		ID:           testID,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)

	user, err := authService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, testID, user.ID, "ID should match")
	assert.Equal(t, signUpParams.Email, user.Email, "Emails should match")
	assert.NotZero(t, user.CreatedAt, "CreatedAt should not be zero")
	assert.NotZero(t, user.UpdatedAt, "UpdatedAt should not be zero")
}

func TestAuthService_SignUpUser_PasswordHashed(t *testing.T) {
	mockRepo, mockHasher, authService := setupMocks(t)
	ctx := context.Background()
	now := time.Now().UTC()

	signUpParams := newSignUpParams()
	createParams := model.User{
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
	}
	mockRepo.EXPECT().FindUserByEmail(ctx, signUpParams.Email).Return(nil, nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)
	mockRepo.EXPECT().CreateUser(ctx, createParams).Return(&model.User{
		ID:           testID,
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	user, err := authService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, hashedPassword, user.PasswordHash, "password hash must match")
}

func TestAuthService_SignUpUser_Duplicate(t *testing.T) {
	mockRepo, mockHasher, authService := setupMocks(t)
	ctx := context.Background()
	signUpParams := newSignUpParams()

	mockRepo.EXPECT().FindUserByEmail(ctx, signUpParams.Email).Return(&model.User{}, nil)
	mockHasher.EXPECT().Hash(gomock.Any()).Times(0)
	mockRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Times(0)

	user, err := authService.SignUpUser(ctx, signUpParams)
	assert.Error(t, err, "signup should return an error")
	assert.ErrorIs(t, err, service.ErrEmailTaken, "errors should match")
	assert.Nil(t, user, "user should be nil")
}

func TestAuthService_SignInUser_Success(t *testing.T) {
	ctx := context.Background()
	input := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}
	mockRepo, mockHasher, authService := setupMocks(t)
	mockRepo.EXPECT().FindUserByEmail(ctx, input.Email).Return(&model.User{
		ID:           testID,
		PasswordHash: hashedPassword,
	}, nil)
	mockHasher.EXPECT().Verify(input.Password, hashedPassword).Return(true, nil)

	id, err := authService.SignInUser(ctx, input)
	assert.NoError(t, err, "signin should not return an error")
	assert.Equal(t, testID, id, "ID must match")
}

func TestAuthService_SignInUser_NotFound(t *testing.T) {
	ctx := context.Background()
	mockRepo, mockHasher, authService := setupMocks(t)

	signInParams := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}

	mockRepo.EXPECT().FindUserByEmail(ctx, signInParams.Email).Return(nil, service.ErrUserNotFound)
	mockHasher.EXPECT().Verify(gomock.Any(), gomock.Any()).Times(0)

	id, err := authService.SignInUser(ctx, signInParams)
	assert.Error(t, err, "signin should return an error")
	assert.ErrorIs(t, err, service.ErrUserNotFound, "errors should match")
	assert.Zero(t, id, "ID should be empty")
}

func setupMocks(t *testing.T) (*repoMocks.MockUserRepo, *secMocks.MockHasher, service.AuthService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockRepo := repoMocks.NewMockUserRepo(ctrl)
	mockHasher := secMocks.NewMockHasher(ctrl)
	userService := service.NewAuthService(mockRepo, mockHasher)

	return mockRepo, mockHasher, userService
}

func newSignUpParams() model.UserSignUpParams {
	return model.UserSignUpParams{
		Email:           testEmail,
		Password:        testPassword,
		PasswordConfirm: testPassword,
	}
}
