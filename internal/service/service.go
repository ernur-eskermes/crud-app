package service

import (
	"context"
	"time"

	"github.com/ernur-eskermes/crud-app/pkg/otp"
	cache "github.com/ernur-eskermes/go-homeworks/2-cache-ttl"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"

	"github.com/ernur-eskermes/crud-app/internal/repository"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/hash"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type Books interface {
	Create(ctx context.Context, book core.CreateBookInput, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (core.Book, error)
	GetAll(ctx context.Context) ([]core.Book, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, id, userID uuid.UUID, inp core.UpdateBookInput) error
}

type Users interface {
	SignUp(ctx context.Context, input UserSignUpInput) error
	SignIn(ctx context.Context, input UserSignInInput) (Tokens, error)
	GetByID(ctx context.Context, id uuid.UUID) (core.User, error)
	Verify(ctx context.Context, username, code string) error
}

type Services struct {
	Users Users
	Books Books
}

type Deps struct {
	Repos          *repository.Repositories
	OtpGenerator   otp.Generator
	Hasher         hash.PasswordHasher
	TokenManager   auth.TokenManager
	AccessTokenTTL time.Duration
	Cache          cache.Cache
	Environment    string
	Domain         string
}

func NewServices(deps Deps) *Services {
	usersService := NewUsersService(deps.Repos.Users, deps.Hasher, deps.TokenManager,
		deps.AccessTokenTTL, deps.Domain, deps.Cache, deps.OtpGenerator)
	booksService := NewBooksService(deps.Repos.Books, deps.TokenManager)

	return &Services{
		Users: usersService,
		Books: booksService,
	}
}

type UserSignUpInput struct {
	Username string
	Password string
}

type UserSignInInput struct {
	Username string
	Password string
}

type Tokens struct {
	AccessToken string
}
