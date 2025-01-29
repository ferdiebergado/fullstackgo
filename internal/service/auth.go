package service

import (
	"context"

	"github.com/ferdiebergado/fullstackgo/internal/db"
	"github.com/ferdiebergado/fullstackgo/internal/model"
)

//go:generate mockgen -destination=mocks/auth_service_mock.go -package=mocks . AuthService
type AuthService interface {
	SignUpUser(ctx context.Context, params model.UserSignUpParams) (*model.User, error)
}

type authService struct {
	repo      db.Authenticator
	validator Validator
	hasher    Hasher
}

func NewUserService(repo db.Authenticator, validator Validator, hasher Hasher) AuthService {
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
		return nil, err
	}

	params.Password = hash

	return s.repo.SignUpUser(ctx, params)
}
