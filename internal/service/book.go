package service

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
)

type BooksRepository interface {
	Create(ctx context.Context, book core.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (core.Book, error)
	GetAll(ctx context.Context) ([]core.Book, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, inp core.Book) error
}

type BooksService struct {
	repo         BooksRepository
	tokenManager auth.TokenManager
}

func NewBooksService(repo BooksRepository, tokenManager auth.TokenManager) *BooksService {
	return &BooksService{
		repo:         repo,
		tokenManager: tokenManager,
	}
}

func (b *BooksService) Create(ctx context.Context, book core.CreateBookInput, userID uuid.UUID) error {
	if book.PublishDate.IsZero() {
		book.PublishDate = time.Now()
	}

	return b.repo.Create(ctx, core.Book{
		Title:       book.Title,
		Author:      userID,
		PublishDate: book.PublishDate,
		Rating:      book.Rating,
	})
}

func (b *BooksService) GetByID(ctx context.Context, id uuid.UUID) (core.Book, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *BooksService) GetAll(ctx context.Context) ([]core.Book, error) {
	return b.repo.GetAll(ctx)
}

func (b *BooksService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return b.repo.Delete(ctx, id, userID)
}

func (b *BooksService) Update(ctx context.Context, id, userID uuid.UUID, inp core.UpdateBookInput) error {
	return b.repo.Update(ctx, core.Book{
		ID:          id,
		Title:       inp.Title,
		Author:      userID,
		PublishDate: inp.PublishDate,
		Rating:      inp.Rating,
	})
}
