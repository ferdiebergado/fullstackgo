package model

import "time"

type AuthMethod string

const (
	BasicAuth AuthMethod = "email/password"
	Oath      AuthMethod = "oauth"
)

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserSignUpParams struct {
	Email           string `json:"email,omitempty" validate:"required"`
	Password        string `json:"password,omitempty" validate:"required"`
	PasswordConfirm string `json:"password_confirm,omitempty" validate:"required,eqfield=password"`
}
