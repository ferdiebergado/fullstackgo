package api_test

import (
	"bytes"
	"encoding/json"
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
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

func TestHandleUserSignUpSuccess(t *testing.T) {
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
	req.Header.Set("content-type", contentType)
	rr := httptest.NewRecorder()

	now := time.Now().UTC()
	ctrl := gomock.NewController(t)
	mockService := mocks.NewMockAuthService(ctrl)
	mockService.EXPECT().SignUpUser(req.Context(), newUser).Return(&model.User{
		ID:        testID,
		Email:     newUser.Email,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	handler := api.NewUserHandler(mockService)
	handler.HandleUserSignUp(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code, "Response status code should match")

	actualContentType := rr.Header().Get("content-type")
	assert.Equal(t, contentType, actualContentType, "Content-Type header should match")

	var user model.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	assert.Equal(t, testID, user.ID, "ID should match")
	assert.Equal(t, newUser.Email, user.Email, "Emails should match")
	assert.Equal(t, now, user.CreatedAt.UTC(), "CreatedAt should match now")
	assert.Equal(t, now, user.UpdatedAt.UTC(), "UpdatedAt should match now")
}
