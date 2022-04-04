package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/internal/repository/postgres"
	"github.com/ernur-eskermes/crud-app/pkg/database/postgresql"
)

//go:generate mockgen -source=repository.go -destination=mocks/mock.go

type Users interface {
	Create(ctx context.Context, user *core.User) error
	GetByCredentials(ctx context.Context, email, password string) (core.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (core.User, error)
	Verify(ctx context.Context, username string) error
}

type Books interface {
	Create(ctx context.Context, book core.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (core.Book, error)
	GetAll(ctx context.Context) ([]core.Book, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, inp core.Book) error
}

type Repositories struct {
	Users Users
	Books Books
}

func NewRepositories(db postgresql.Client) *Repositories {
	return &Repositories{
		Users: postgres.NewUsersRepo(db),
		Books: postgres.NewBooksRepo(db),
	}
}
