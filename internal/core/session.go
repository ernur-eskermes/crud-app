package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrTokenNotFound = errors.New("token not found")

	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

type RefreshSession struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"-"`
}
