package repository

import (
	"github.com/ernur-eskermes/crud-app/domain"
	"github.com/jmoiron/sqlx"
)

func GetBookById(db sqlx.DB, id int) (domain.Book, error) {
	var book domain.Book

	if err := db.Get(&book, "SELECT * FROM books WHERE id=$1", id); err != nil {
		return domain.Book{}, err
	}
	return book, nil
}

func UpdateBook(db sqlx.DB, id int) {
	return
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