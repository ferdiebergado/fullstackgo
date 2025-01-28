package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/api"
	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/service"
)

const (
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

func TestCreateUserAPI(t *testing.T) {
	newUser := model.UserCreateParams{
		Email:           testEmail,
		Password:        testPassword,
		PasswordConfirm: testPassword,
		AuthMethod:      authMethod,
	}

	userJson, err := json.Marshal(newUser)

	if err != nil {
		t.Errorf("wanted no error, but got: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer(userJson))
	rr := httptest.NewRecorder()

	now := time.Now().UTC()
	service := &service.MockUserService{
		CreateUserFn: func(ctx context.Context, params model.UserCreateParams) (*model.User, error) {
			return &model.User{
				Email:      testEmail,
				AuthMethod: authMethod,
				CreatedAt:  now,
				UpdatedAt:  now,
			}, nil
		},
	}
	handler := api.NewUserHandler(service)

	handler.HandleCreateUser(rr, req)

	wantedStatus := http.StatusCreated
	if rr.Code != wantedStatus {
		t.Errorf("want: %d; got: %d", wantedStatus, rr.Code)
	}

	wantedContentType := "application/json"
	actualContentType := rr.Header().Get("content-type")

	if actualContentType != wantedContentType {
		t.Errorf("want: %s, got: %s", wantedContentType, actualContentType)
	}

	var user model.User
	if err := json.NewDecoder(rr.Body).Decode(&user); err != nil {
		t.Fatalf("decode json: %v", err)
	}

	if user.Email != newUser.Email {
		t.Errorf("want: %s; got: %s", newUser.Email, user.Email)
	}

	if string(user.AuthMethod) != string(newUser.AuthMethod) {
		t.Errorf("want: %s; got: %s", newUser.AuthMethod, user.AuthMethod)
	}

	if user.CreatedAt.UTC() != now {
		t.Errorf("want: %v; got: %v", now, user.CreatedAt)
	}

	if user.UpdatedAt.UTC() != now {
		t.Errorf("want: %v; got: %v", now, user.UpdatedAt)
	}
}
