package model

import "time"

type AuthMethod string

const (
	BasicAuth AuthMethod = "email/password"
	Oath      AuthMethod = "oauth"
)

type User struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	AuthMethod   AuthMethod `json:"auth_method"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type UserCreateParams struct {
	Email           string     `json:"email,omitempty"`
	Password        string     `json:"password,omitempty"`
	PasswordConfirm string     `json:"password_confirm,omitempty"`
	AuthMethod      AuthMethod `json:"auth_method"`
}
