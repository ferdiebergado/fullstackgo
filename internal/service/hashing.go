package service

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/argon2"
)

//go:generate mockgen -destination=mocks/hasher_mock.go -package=mocks . Hasher
type Hasher interface {
	Hash(plain string) (string, error)
	Verify(plain, hashed string) (bool, error)
}

type Argon2Hasher struct{}

var _ Hasher = (*Argon2Hasher)(nil)

// Parameters for the Argon2ID algorithm
const (
	Memory      = 64 * 1024 // 64 MB
	Iterations  = 3
	Parallelism = 2
	SaltLength  = 16 // 16 bytes
	KeyLength   = 32 // 32 bytes
)

// Hash implements Hasher.
func (h *Argon2Hasher) Hash(plain string) (string, error) {
	// Generate a random salt
	salt, err := GenerateRandomBytes(SaltLength)
	if err != nil {
		return "", fmt.Errorf("generate salt: %w", err)
	}

	// Hash the password
	hash := argon2.IDKey([]byte(plain), salt, Iterations, Memory, Parallelism, KeyLength)

	// Encode the salt and hash for storage
	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hash)

	// Return the formatted password hash
	return fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		Memory, Iterations, Parallelism, saltBase64, hashBase64), nil
}

// Verify implements Hasher.
func (h *Argon2Hasher) Verify(plain string, hashed string) (bool, error) {
	panic("unimplemented")
}
