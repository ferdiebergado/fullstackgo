//go:generate mockgen -destination=mocks/auth_service_mock.go -package=mocks . AuthService
package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ferdiebergado/fullstackgo/internal/db"
	"github.com/ferdiebergado/fullstackgo/internal/model"
)

var ErrDuplicateUser = errors.New("user already exists")

type AuthService interface {
	SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error)
	SignInUser(ctx context.Context, params model.UserSignInParams) (string, error)
}

type authService struct {
	repo      db.Authenticator
	validator Validator
	hasher    Hasher
}

func NewAuthService(repo db.Authenticator, validator Validator, hasher Hasher) AuthService {
	return &authService{
		repo:      repo,
		validator: validator,
		hasher:    hasher,
	}
}

func (s *authService) SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error) {
	if err := s.validator.Struct(params); err != nil {
		return nil, err
	}

	hash, err := s.hasher.Hash(params.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	params.Password = hash

	return s.repo.SignUpUser(ctx, params)
}

func (s *authService) SignInUser(ctx context.Context, params model.UserSignInParams) (string, error) {
	if err := s.validator.Struct(params); err != nil {
		return "", err
	}

	user, err := s.repo.SignInUser(ctx, params)
	if err != nil {
		return "", err
	}

	ok, err := s.hasher.Verify(params.Password, user.PasswordHash)

	if err != nil {
		return "", err
	}

	if !ok {
		return "", errors.New("passwords do not match")
	}

	return user.ID, nil
}
