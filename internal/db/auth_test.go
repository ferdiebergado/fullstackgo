package db_test

import (
	"context"
	"database/sql/driver"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ferdiebergado/fullstackgo/internal/db"
	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/stretchr/testify/assert"
)

const (
	testID       = "1"
	testEmail    = "abc@example.com"
	testPassword = "hashed"
	authMethod   = model.BasicAuth
)

var sqlmockOpts = sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual)

func TestSignUpUserSuccess(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmockOpts)

	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	defer mockDB.Close()

	repo := db.NewUserRepo(mockDB)

	params := model.UserSignUpParams{
		Email:    testEmail,
		Password: testPassword,
	}

	now := time.Now().UTC()
	cols := []string{"id", "email", "created_at", "updated_at"}
	row := []driver.Value{testID, testEmail, now, now}

	mock.ExpectQuery(db.SignUpUserQuery).WithArgs(params.Email, params.Password).WillReturnRows(sqlmock.NewRows(cols).AddRow(row...))
	user, err := repo.SignUpUser(context.Background(), params)

	assert.NoError(t, err, "signup should not return an error")
	assert.Equal(t, testID, user.ID, "ID must match")
	assert.Equal(t, testEmail, user.Email, "email must match")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func TestSignUpUserDuplicateUser(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmockOpts)

	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	defer mockDB.Close()

	repo := db.NewUserRepo(mockDB)

	params := model.UserSignUpParams{
		Email:    testEmail,
		Password: testPassword,
	}

	mock.ExpectQuery(db.SignUpUserQuery).WithArgs(params.Email, params.Password).WillReturnError(db.ErrDuplicateUser)

	_, err = repo.SignUpUser(context.Background(), params)

	assert.Error(t, err, "signup should return an error")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}

func TestSignUpUserInvalidData(t *testing.T) {
	mockDB, mock, err := sqlmock.New(sqlmockOpts)

	if err != nil {
		t.Fatalf("failed to create mock db: %v", err)
	}

	defer mockDB.Close()

	repo := db.NewUserRepo(mockDB)

	params := model.UserSignUpParams{
		Email:    testEmail,
		Password: testPassword,
	}

	mock.ExpectQuery(db.SignUpUserQuery).WithArgs(testEmail, testPassword).WillReturnError(db.ErrNullValue)

	_, err = repo.SignUpUser(context.Background(), params)

	assert.Error(t, err, "signup should return an error")
	assert.NoError(t, mock.ExpectationsWereMet(), "some expectations were not met")
}
