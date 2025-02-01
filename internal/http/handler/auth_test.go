package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/http/handler"
	"github.com/ferdiebergado/fullstackgo/internal/model"
	validationMocks "github.com/ferdiebergado/fullstackgo/internal/pkg/validation/mocks"
	"github.com/ferdiebergado/fullstackgo/internal/service"
	"github.com/ferdiebergado/fullstackgo/internal/service/mocks"
	"github.com/go-playground/validator/v10"

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
	params := model.UserSignUpParams{
		Email:           testEmail,
		Password:        testPassword,
		PasswordConfirm: testPassword,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("json.Marshal: %v, err: %v", params, err)
	}

	req := httptest.NewRequest(http.MethodPost, signUpURL, bytes.NewBuffer(jsonParams))
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	mockService, mockValidator, authHandler := setupMockService(t)
	mockValidator.EXPECT().Struct(params).Return(nil)
	mockService.EXPECT().SignUpUser(req.Context(), params).DoAndReturn(
		func(_ context.Context, params model.UserSignUpParams) (*model.User, error) {
			tNow := time.Now().UTC()
			return &model.User{
				ID:        testID,
				Email:     params.Email,
				CreatedAt: tNow,
				UpdatedAt: tNow,
			}, nil
		},
	)

	authHandler.HandleUserSignUp(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Response status code should match")

	actualContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, contentType, actualContentType, "Content-Type header should match")

	var newUser model.User
	if err = json.Unmarshal(rr.Body.Bytes(), &newUser); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	assert.Equal(t, testID, newUser.ID, "ID should match")
	assert.Equal(t, params.Email, newUser.Email, "email should match")
	assert.NotZero(t, newUser.CreatedAt, "CreatedAt should not be zero")
	assert.NotZero(t, newUser.UpdatedAt, "UpdatedAt should not be zero")
}

func TestAuthHandler_HandleUserSignUp_DuplicateUser(t *testing.T) {
	params := model.UserSignUpParams{
		Email:           testEmail,
		Password:        testPassword,
		PasswordConfirm: testPassword,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("json.Marshal: %v, err: %v", params, err)
	}

	req := httptest.NewRequest(http.MethodPost, signUpURL, bytes.NewBuffer(jsonParams))
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	mockService, mockValidator, authHandler := setupMockService(t)
	mockValidator.EXPECT().Struct(params).Return(nil)
	mockService.EXPECT().SignUpUser(req.Context(), params).Return(nil, service.ErrEmailTaken)

	authHandler.HandleUserSignUp(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code, "Response status code should match")

	actualContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, contentType, actualContentType, "Content-Type header should match")

	var res handler.APIResponse
	if err = json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	assert.Equal(t, "Invalid input!", res.Message, "Message should match")
	assert.NotEmpty(t, res.Errors, "Errors should not be empty")
	assert.Equal(t, service.ErrEmailTaken.Error(), res.Errors[0]["email"], "Validation errors should match")
}

func TestAuthHandler_HandleUserSignUp_InvalidInput(t *testing.T) {
	mockService, mockValidator, authHandler := setupMockService(t)

	tests := []struct {
		name  string
		given model.UserSignUpParams
		field string
		tag   string
	}{
		{"should fail when email is empty", model.UserSignUpParams{
			Email:           "",
			Password:        testPassword,
			PasswordConfirm: testPassword,
		},
			"email",
			"required",
		},
		{"should fail when email is invalid", model.UserSignUpParams{
			Email:           "abcd",
			Password:        testPassword,
			PasswordConfirm: testPassword,
		}, "Email", "email"},
		{"should fail when password is empty", model.UserSignUpParams{
			Email:           testEmail,
			Password:        "",
			PasswordConfirm: "otherpassword",
		}, "Password", "required",
		},
		{"should fail when passwords do not match", model.UserSignUpParams{
			Email:           testEmail,
			Password:        testPassword,
			PasswordConfirm: "otherpassword",
		}, "PasswordConfirm", "eqfield",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonParams, err := json.Marshal(tt.given)
			if err != nil {
				t.Fatalf("json.Marshal: %v, err: %v", tt.given, err)
			}

			mockValidator.EXPECT().Struct(tt.given).Return(&validator.ValidationErrors{
				validationMocks.MockFieldError{
					FieldName: tt.field,
					TagName:   tt.tag,
				},
			})
			mockService.EXPECT().SignUpUser(gomock.Any(), gomock.Any()).Times(0)
			req := httptest.NewRequest(http.MethodPost, signUpURL, bytes.NewBuffer(jsonParams))
			req.Header.Set("Content-Type", contentType)
			rr := httptest.NewRecorder()

			authHandler.HandleUserSignUp(rr, req)
			assert.Equal(t, rr.Code, http.StatusUnprocessableEntity, "signup should return http error 422")

			var res handler.APIResponse
			if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
				t.Fatalf("decode json: %v", err)
			}

			assert.Equal(t, "Invalid input!", res.Message, "Message should match")
			assert.NotEmpty(t, res.Errors, "Errors must not be empty")
			assert.Equal(t, fmt.Sprintf("validation failed on field %s with tag %s", tt.field, tt.tag), res.Errors[0][tt.field], "validation error must match")
		})
	}
}

func TestAuthHandler_HandleUserSignIn_Success(t *testing.T) {
	params := model.UserSignInParams{
		Email:    testEmail,
		Password: testPassword,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		t.Fatalf("json.Marshal: %v, err: %v", params, err)
	}

	req := httptest.NewRequest(http.MethodPost, signInURL, bytes.NewBuffer(jsonParams))
	req.Header.Set("Content-Type", contentType)
	rr := httptest.NewRecorder()

	mockService, mockValidator, authHandler := setupMockService(t)
	mockService.EXPECT().SignInUser(req.Context(), params).Return(testID, nil)
	mockValidator.EXPECT().Struct(params).Return(nil)

	authHandler.HandleUserSignIn(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Response status code should match")

	actualContentType := rr.Header().Get("Content-Type")
	assert.Equal(t, contentType, actualContentType, "Content-Type header should match")

	var res handler.APIResponse
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	assert.Equal(t, "Signin successful.", res.Message, "Message should match")
}

func setupMockService(t *testing.T) (*mocks.MockAuthService, *validationMocks.MockValidator, handler.AuthHandler) {
	t.Helper()
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockAuthService(ctrl)
	mockValidator := validationMocks.NewMockValidator(ctrl)
	authHandler := handler.NewAuthHandler(mockService, mockValidator)

	return mockService, mockValidator, authHandler
}
