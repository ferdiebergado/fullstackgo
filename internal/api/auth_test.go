package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/api"
	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

const (
	contentType = "application/json"
	signUpURL   = "/api/signup"
	signInURL   = "/api/signin"
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
)

func TestAuthHandler_HandleUserSignUp_Success(t *testing.T) {
	newUser := model.UserSignUpParams{
		Email:           testEmail,
		Password:        testPassword,
		PasswordConfirm: testPassword,
	}

	userJSON, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("json.Marshal: %v, err: %v", newUser, err)
	}

	req := httptest.NewRequest(http.MethodPost, signUpURL, bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	now := time.Now().UTC().Truncate(time.Millisecond)
	mockService, handler := setupMockService(t)
	mockService.EXPECT().SignUpUser(req.Context(), newUser).DoAndReturn(
		func(_ context.Context, params model.UserSignUpParams) (*model.User, error) {
			tNow := time.Now().UTC().Truncate(time.Millisecond)
			return &model.User{
				ID:        testID,
				Email:     params.Email,
				CreatedAt: tNow,
				UpdatedAt: tNow,
			}, nil
		},
	)

	handler.HandleUserSignUp(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Response status code should match")

	actualContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, contentType, actualContentType, "Content-Type header should match")

	expectedJSON := fmt.Sprintf(`{"id": "%s", "email": "%s", "created_at": "%s", "updated_at": "%s"}`,
		testID, newUser.Email, now.Format(time.RFC3339Nano), now.Format(time.RFC3339Nano))
	assert.JSONEq(t, expectedJSON, rr.Body.String(), "Response body should match expected JSON")
}

func TestAuthHandler_HandleUserSignIn_Success(t *testing.T) {
	signInParams := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}

	signInJSON, err := json.Marshal(signInParams)
	if err != nil {
		t.Fatalf("json.Marshal: %v, err: %v", signInParams, err)
	}

	req := httptest.NewRequest(http.MethodPost, signInURL, bytes.NewBuffer(signInJSON))
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	mockService, handler := setupMockService(t)
	mockService.EXPECT().SignInUser(req.Context(), signInParams).Return(testID, nil)

	handler.HandleUserSignIn(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code, "Response status code should match")

	actualContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, contentType, actualContentType, "Content-Type header should match")
}

func setupMockService(t *testing.T) (*mocks.MockAuthService, api.AuthHandler) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockAuthService(ctrl)
	handler := api.NewAuthHandler(mockService)

	return mockService, handler
}
