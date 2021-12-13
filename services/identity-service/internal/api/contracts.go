package api

import (
	"identity-service/internal/users"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) (bool, error)
}

type TokenIssuer interface {
	Issue(user *users.User) (string, error)
}
