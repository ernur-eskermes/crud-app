package repository

import (
	"database/sql"
	"errors"
	"github.com/ernur-eskermes/crud-app/domain"
	"github.com/jmoiron/sqlx"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

func GetBookById(db sqlx.DB, id int) (domain.Book, error) {
	var book domain.Book

	if err := db.Get(&book, "SELECT * FROM books WHERE id=$1", id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Book{}, ErrBookNotFound
		}
		return domain.Book{}, err
	}
	return book, nil
}

func UpdateBook(db sqlx.DB, book domain.Book) error {
	_, err := db.NamedExec("UPDATE books SET title = :title, author = :author,publish_date = :publish_date, rating = :rating WHERE id=:id;", &book)
	return err
}

func GetBooks(db sqlx.DB) ([]domain.Book, error) {
	var books []domain.Book

	if err := db.Select(&books, "SELECT * FROM books"); err != nil {
		return nil, err
	}
	return books, nil
}

func CreateBook(db sqlx.DB, book domain.Book) error {
	_, err := db.NamedExec("INSERT INTO books (title, author, publish_date, rating) VALUES (:title, :author, :publish_date, :rating)", &book)
	if err != nil {
		return err
	}
	return nil
}

func DeleteBook(db sqlx.DB, id int) error {
	if _, err := db.Exec("DELETE FROM books WHERE id=$1", id); err != nil {
		return err
	}
	return nil
}
