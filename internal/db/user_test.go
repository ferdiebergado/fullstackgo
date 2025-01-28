package db_test

import (
	"context"
	"slices"
	"testing"
	"time"

	"github.com/ferdiebergado/fullstackgo/internal/db"
	"github.com/ferdiebergado/fullstackgo/internal/model"
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

const createUserQuery = `
INSERT into users (email, password_hash, auth_method)
VALUES $1, $2, $3 
RETURNING id, email, auth_method, created_at, updated_at
`

func TestCreateUserRepoSuccess(t *testing.T) {
	now := time.Now().UTC()
	mockDB := &db.MockDB{
		QueryRowContextFn: func(tx context.Context, query string, args ...any) db.Row {
			if query != createUserQuery {
				t.Errorf("want: %s, got: %s", createUserQuery, query)
			}

			wantedArgs := []any{testEmail, testPassword, authMethod}
			if !slices.Equal(args, wantedArgs) {
				t.Errorf("want: %s, got: %s", wantedArgs, args)
			}

			return &db.MockRow{
				ScanFn: func(dest ...any) error {
					if len(dest) != 5 {
						t.Fatalf("expected 5 destinations, got %d", len(dest))
					}
					*(dest[0].(*string)) = testID
					*(dest[1].(*string)) = testEmail
					*(dest[2].(*model.AuthMethod)) = authMethod
					*(dest[3].(*time.Time)) = now
					*(dest[4].(*time.Time)) = now
					return nil
				},
			}
		},
	}

	repo := db.NewUserRepo(mockDB)

	params := model.UserCreateParams{
		Email:      testEmail,
		Password:   testPassword,
		AuthMethod: authMethod,
	}

	user, err := repo.CreateUser(context.Background(), params)

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

	if user.CreatedAt.UTC() != now {
		t.Errorf("want: %v; got: %v", now, user.CreatedAt)
	}

	if user.UpdatedAt.UTC() != now {
		t.Errorf("want: %v; got: %v", now, user.UpdatedAt)
	}
}

func TestUserRepoDuplicateUser(t *testing.T) {
	// TODO:
}

func TestUserRepoInvalidData(t *testing.T) {
	// TODO:
}
