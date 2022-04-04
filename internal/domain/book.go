package domain

import (
	"errors"
	"time"
)

var (
	ErrBookNotFound = errors.New("book not found")
)

type Book struct {
	Id          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Author      string    `json:"author" db:"author"`
	PublishDate time.Time `json:"publish_date" db:"publish_date"`
	Rating      int       `json:"rating" db:"rating"`
}
