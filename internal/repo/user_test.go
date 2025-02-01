package repo_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/repo"
	"github.com/ferdiebergado/fullstackgo/internal/service"
	"github.com/stretchr/testify/assert"
)

const (
	testID             = "1"
	testEmail          = "abc@example.com"
	testPassword       = "test"
	testPasswordHashed = "hashed"
)

var sqlmockOpts = sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)

func TestAuthRepo_SignUpUser_Success(t *testing.T) {
	mock, userRepo := setupMockDB(t)

	params := model.User{
		Email:        testEmail,
		PasswordHash: testPasswordHashed,
	}

	now := time.Now().UTC()
	cols := []string{"id", "email", "created_at", "updated_at"}

	mock.ExpectQuery(repo.CreateUserQuery).
		WithArgs(params.Email, params.PasswordHash).
		WillReturnRows(sqlmock.NewRows(cols).
			AddRow(testID, testEmail, now, now))
	user, err := userRepo.CreateUser(context.Background(), params)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, testID, user.ID, "ID must match")
	assert.Equal(t, testEmail, user.Email, "email must match")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func TestAuthRepo_SignUpUser_Duplicate(t *testing.T) {
	mock, userRepo := setupMockDB(t)

	params := model.User{
		Email:        testEmail,
		PasswordHash: testPasswordHashed,
	}

	mock.ExpectQuery(repo.CreateUserQuery).
		WithArgs(params.Email, params.PasswordHash).
		WillReturnError(service.ErrModelExists)

	_, err := userRepo.CreateUser(context.Background(), params)

	assert.Error(t, err, "signup should return an error")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func TestAuthRepo_SignUpUser_InvalidData(t *testing.T) {
	mock, userRepo := setupMockDB(t)

	params := model.User{
		Email:        testEmail,
		PasswordHash: testPasswordHashed,
	}

	mock.ExpectQuery(repo.CreateUserQuery).
		WithArgs(params.Email, params.PasswordHash).
		WillReturnError(repo.ErrNullValue)

	_, err := userRepo.CreateUser(context.Background(), params)

	assert.Error(t, err, "signup should return an error")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func TestAuthRepo_FindUserByEmail_Success(t *testing.T) {
	mock, userRepo := setupMockDB(t)
	mock.ExpectQuery(repo.FindUserByEmailQuery).
		WithArgs(testEmail).
		WillReturnRows(sqlmock.NewRows([]string{"id", "password_hash"}).
			AddRow(testID, testPasswordHashed))

	user, err := userRepo.FindUserByEmail(context.Background(), testEmail)
	assert.NoError(t, err, "signin should not return an error")
	assert.Equal(t, testPasswordHashed, user.PasswordHash, "password hash must match")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func TestAuthRepo_FindUserByEmail_UserNotFound(t *testing.T) {
	mock, userRepo := setupMockDB(t)
	mock.ExpectQuery(repo.FindUserByEmailQuery).WithArgs(testEmail).WillReturnError(sql.ErrNoRows)

	_, err := userRepo.FindUserByEmail(context.Background(), testEmail)
	assert.Error(t, err, "signin should return an error")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, repo.UserRepo) {
	t.Helper()
	mockDB, mock, err := sqlmock.New(sqlmockOpts)
	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}
	repo := repo.NewUserRepo(mockDB)
	t.Cleanup(func() { mockDB.Close() })
	return mock, repo
}
