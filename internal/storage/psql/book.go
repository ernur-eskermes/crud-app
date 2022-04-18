package psql

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/pkg/database/postgresql"
	"github.com/jackc/pgx/v4"
)

type BooksRepo struct {
	db postgresql.Client
}

func NewBooksRepo(db postgresql.Client) *BooksRepo {
	return &BooksRepo{db}
}

func (b *BooksRepo) GetByID(ctx context.Context, id uuid.UUID) (core.Book, error) {
	q := "SELECT id, title, author_id, publish_date, rating FROM book WHERE id=$1"

	var book core.Book

	if err := b.db.QueryRow(ctx, q, id).Scan(
		&book.ID,
		&book.Title,
		&book.Author,
		&book.PublishDate,
		&book.Rating,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core.Book{}, core.ErrBookNotFound
		}

		return core.Book{}, err
	}

	return book, nil
}

func (b *BooksRepo) Update(ctx context.Context, book core.Book) error {
	q := "UPDATE book SET title=$1, publish_date=$2, rating=$3 WHERE id=$4 and author_id=$5"
	res, err := b.db.Exec(ctx, q, book.Title, book.PublishDate, book.Rating, book.ID, book.Author)

	if res.RowsAffected() == 0 {
		return core.ErrBookNotFound
	}

	return err
}

func (b *BooksRepo) GetAll(ctx context.Context) ([]core.Book, error) {
	q := "SELECT id, title, author_id, publish_date, rating FROM book"

	var books []core.Book

	rows, err := b.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var book core.Book

		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.PublishDate, &book.Rating)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

func (b *BooksRepo) Create(ctx context.Context, book core.Book) error {
	q := "INSERT INTO book (title, author_id, publish_date, rating) VALUES ($1, $2, $3, $4) RETURNING id"

	return b.db.QueryRow(ctx, q, book.Title, book.Author, book.PublishDate, book.Rating).Scan(&book.ID)
}

func (b *BooksRepo) Delete(ctx context.Context, id, userID uuid.UUID) error {
	q := "DELETE FROM book WHERE id=$1 and author_id=$2"

	res, err := b.db.Exec(ctx, q, id, userID)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return core.ErrBookNotFound
	}

	return nil
}
