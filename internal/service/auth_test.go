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
	createParams := &model.UserCreateParams{
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
	}

	mockRepo.EXPECT().CreateUser(ctx, *createParams).Return(&model.User{
		ID:           testID,
		Email:        testEmail,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)

	user, err := authService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, testID, user.ID, "ID should match")
	assert.Equal(t, signUpParams.Email, user.Email, "Emails should match")
	assert.Equal(t, now, user.CreatedAt.UTC(), "CreatedAt should match now")
	assert.Equal(t, now, user.UpdatedAt.UTC(), "UpdatedAt should match now")
}

func TestAuthService_SignUpUser_PasswordHashed(t *testing.T) {
	mockRepo, mockHasher, authService := setupMocks(t)
	ctx := context.Background()
	now := time.Now().UTC()

	signUpParams := newSignUpParams()
	createParams := model.UserCreateParams{
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
	}
	mockRepo.EXPECT().CreateUser(ctx, createParams).Return(&model.User{
		ID:           testID,
		Email:        signUpParams.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil)
	mockHasher.EXPECT().Hash(signUpParams.Password).Return(hashedPassword, nil)
	user, err := authService.SignUpUser(ctx, signUpParams)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, hashedPassword, user.PasswordHash, "password hash must match")
}

func TestAuthService_SignInUser_Success(t *testing.T) {
	ctx := context.Background()
	mockRepo, mockHasher, authService := setupMocks(t)

	signInParams := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}

	mockRepo.EXPECT().FindUserByEmail(ctx, signInParams.Email).Return(&model.User{
		ID:           testID,
		PasswordHash: hashedPassword,
	}, nil)
	mockHasher.EXPECT().Verify(signInParams.Password, hashedPassword).Return(true, nil)

	id, err := authService.SignInUser(ctx, signInParams)

	assert.NoError(t, err, "signin should not return an error")
	assert.Equal(t, testID, id, "ID must match")
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
