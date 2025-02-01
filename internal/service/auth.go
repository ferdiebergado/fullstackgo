//go:generate mockgen -destination=mocks/auth_service_mock.go -package=mocks . AuthService
package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ferdiebergado/fullstackgo/internal/model"
	"github.com/ferdiebergado/fullstackgo/internal/pkg/security"
	"github.com/ferdiebergado/fullstackgo/internal/repo"
)

type AuthService interface {
	SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error)
	SignInUser(ctx context.Context, params model.UserSignInParams) (string, error)
}

type authService struct {
	repo   repo.UserRepo
	hasher security.Hasher
}

func NewAuthService(repo repo.UserRepo, hasher security.Hasher) AuthService {
	return &authService{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *authService) SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error) {
	existing, err := s.repo.FindUserByEmail(ctx, params.Email)

	if err != nil {
		return nil, err
	}

	if existing != nil {
		return nil, ErrModelExists
	}

	hash, err := s.hasher.Hash(params.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	createParams := model.User{
		Email:        params.Email,
		PasswordHash: hash,
	}

	return s.repo.CreateUser(ctx, createParams)
}

func (s *authService) SignInUser(ctx context.Context, params model.UserSignInParams) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrModelNotFound
		}

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
