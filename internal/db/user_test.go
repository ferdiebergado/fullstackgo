package db_test

import (
	"context"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/db"
	"github.com/ferdiebergado/fullstackgo/internal/db/mocks"
	"github.com/ferdiebergado/fullstackgo/internal/model"
	"go.uber.org/mock/gomock"
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

func TestCreateUserRepoSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDB := mocks.NewMockQuerier(ctrl)
	mockRow := mocks.NewMockRow(ctrl)
	ctx := context.Background()
	mockDB.EXPECT().QueryRowContext(ctx, db.CreateUserQuery, testEmail, testPassword, authMethod).Return(mockRow)
	now := time.Now().UTC()
	user := &model.User{}
	mockRow.EXPECT().Scan(&user.ID, &user.Email, &user.AuthMethod, &user.CreatedAt, &user.UpdatedAt).Do(func(dest ...any) {
		if len(dest) > 0 {
			*(dest[0].(*string)) = testID
			*(dest[1].(*string)) = testEmail
			*(dest[2].(*model.AuthMethod)) = authMethod
			*(dest[3].(*time.Time)) = now
			*(dest[4].(*time.Time)) = now
		}
	}).Return(nil)

	repo := db.NewUserRepo(mockDB)

	params := model.UserCreateParams{
		Email:      testEmail,
		Password:   testPassword,
		AuthMethod: authMethod,
	}

	user, err := repo.CreateUser(ctx, params)

	if err != nil {
		t.Errorf("wanted no error, but got: %v", err)
	}

	if user.ID != testID {
		t.Errorf("want: %s; got: %s", testID, user.ID)
	}

	if user.Email != testEmail {
		t.Errorf("want: %s, but got: %s", testEmail, user.Email)
	}

	if user.AuthMethod != authMethod {
		t.Errorf("want: %s; but got: %s", authMethod, user.AuthMethod)
	}
}

func TestUserRepoDuplicateUser(t *testing.T) {
	// TODO:
}

func TestUserRepoInvalidData(t *testing.T) {
	// TODO:
}
