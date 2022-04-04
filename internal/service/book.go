package service

import (
	"context"
	"github.com/ernur-eskermes/crud-app/internal/domain"
	"time"
)

type BooksRepository interface {
	Create(ctx context.Context, book domain.Book) error
	GetByID(ctx context.Context, id int) (domain.Book, error)
	GetAll(ctx context.Context) ([]domain.Book, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, inp domain.Book) error
}

type Books struct {
	repo BooksRepository
}

func NewBooks(repo BooksRepository) *Books {
	return &Books{
		repo: repo,
	}
}

func (b *Books) Create(ctx context.Context, book domain.Book) error {
	if book.PublishDate.IsZero() {
		book.PublishDate = time.Now()
	}

	return b.repo.Create(ctx, book)
}

func (b *Books) GetByID(ctx context.Context, id int) (domain.Book, error) {
	return b.repo.GetByID(ctx, id)
}

func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
	return b.repo.GetAll(ctx)
}

func (b *Books) Delete(ctx context.Context, id int) error {
	return b.repo.Delete(ctx, id)
}

func (b *Books) Update(ctx context.Context, inp domain.Book) error {
	return b.repo.Update(ctx, inp)
}
