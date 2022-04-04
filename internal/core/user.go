package core

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with such username already exists")

	ErrUserCodeExpired     = errors.New("code has expired. repeat again")
	ErrUserCodeUnknownType = errors.New("code of unknown type")
	ErrUserCodeIncorrect   = errors.New("code is incorrect")
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}
