package postgres

import (
	"context"
	"database/sql"
	"errors"
	"github.com/ernur-eskermes/crud-app/internal/domain"
	"github.com/jmoiron/sqlx"
)

type Books struct {
	db *sqlx.DB
}

func NewBooks(db *sqlx.DB) *Books {
	return &Books{db}
}

func (b *Books) GetByID(ctx context.Context, id int) (domain.Book, error) {
	var book domain.Book

	if err := b.db.Get(&book, "SELECT * FROM books WHERE id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Book{}, domain.ErrBookNotFound
		}
		return domain.Book{}, err
	}
	return book, nil
}

func (b *Books) Update(ctx context.Context, book domain.Book) error {
	_, err := b.db.NamedExec("UPDATE books SET title = :title, author = :author,publish_date = :publish_date, rating = :rating WHERE id=:id;", &book)
	return err
}

func (b *Books) GetAll(ctx context.Context) ([]domain.Book, error) {
	var books []domain.Book

	if err := b.db.Select(&books, "SELECT * FROM books"); err != nil {
		return nil, err
	}
	return books, nil
}

func (b *Books) Create(ctx context.Context, book domain.Book) error {
	_, err := b.db.NamedExec("INSERT INTO books (title, author, publish_date, rating) VALUES (:title, :author, :publish_date, :rating)", &book)
	if err != nil {
		return err
	}
	return nil
}

func (b *Books) Delete(ctx context.Context, id int) error {
	if _, err := b.db.Exec("DELETE FROM books WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}
