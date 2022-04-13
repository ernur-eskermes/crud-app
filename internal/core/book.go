package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var ErrBookNotFound = errors.New("book not found")

type Book struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      uuid.UUID `json:"author"`
	PublishDate time.Time `json:"publish_date"`
	Rating      int       `json:"rating"`
}

type CreateBookInput struct {
	Title       string    `json:"title" validate:"required,max=64"`
	PublishDate time.Time `json:"publish_date" validate:"required"`
	Rating      int       `json:"rating" validate:"required,number,max=5,min=0"`
}

type UpdateBookInput struct {
	Title       string    `json:"title" validate:"required,max=64"`
	PublishDate time.Time `json:"publish_date" validate:"required"`
	Rating      int       `json:"rating" validate:"required,number,max=5,min=0"`
}
