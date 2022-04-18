package core

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user with such username already exists")

	ErrUserCodeIncorrect = errors.New("code is incorrect")
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"-"`
}

// DTO

type AuthInput struct {
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

func (i *AuthInput) Validate() []*ErrorResponse {
	return validateStruct(i)
}

type VerifyInput struct {
	Username string `json:"username" validate:"required"`
	Code     string `json:"code" validate:"required,eq=6"`
}

func (i *VerifyInput) Validate() []*ErrorResponse {
	return validateStruct(i)
}
