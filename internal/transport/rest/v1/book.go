package v1

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/gofiber/fiber/v2"
)

/*
type BookService interface {
	Create(ctx context.Context, book core.CreateBookInput, userID uuid.UUID) error
	GetByID(ctx context.Context, id uuid.UUID) (core.Book, error)
	GetAll(ctx context.Context) ([]core.Book, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, id, userID uuid.UUID, inp core.UpdateBookInput) error
}

type BookHandler struct {
	service BookService
}

func NewBookHandler(service BookService) *BookHandler {
	return &BookHandler{service: service}
}
*/

func (h *Handler) initBooksRoutes(api fiber.Router) {
	books := api.Group("/books")
	{
		books.Get("", h.getAllBooks)
		books.Get("/:id", h.getBookByID)

		authenticated := books.Group("", h.userIdentity)
		{
			authenticated.Post("", h.createBook)
			authenticated.Delete("/:id", h.deleteBook)
			authenticated.Put("/:id", h.updateBook)
		}
	}
}

// @Summary Get Book
// @Tags books
// @Description get book by id
// @ModuleID getBookByID
// @Accept  json
// @Produce  json
// @Param id path string true "book id"
// @Success 200 {object} core.Book
// @Failure 404 {object} response
// @Router /books/{id} [get]
func (h *Handler) getBookByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{"id value of incorrect type"})
	}

	book, err := h.services.Books.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, core.ErrBookNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(book)
}

// @Summary Create Book
// @Tags books
// @Description create book
// @ModuleID createBook
// @Security UsersAuth
// @Accept  json
// @Produce  json
// @Param input body core.CreateBookInput true "create book"
// @Success 201 {string} string "Created"
// @Failure 400 {object} response
// @Router /books [post]
func (h *Handler) createBook(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		h.logger.Warning(err)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if _, err = h.services.Users.GetByID(context.TODO(), userID); err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{"user not found"})
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var inp core.CreateBookInput

	if err = c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := h.validateStruct(inp); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	err = h.services.Books.Create(context.TODO(), inp, userID)
	if err != nil {
		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary Delete Book
// @Tags books
// @Description delete book by id
// @ModuleID deleteBook
// @Security UsersAuth
// @Accept  json
// @Produce  json
// @Param id path string true "book id"
// @Success 204 {string} string "No Content"
// @Failure 400,404 {object} response
// @Router /books/{id} [delete]
func (h *Handler) deleteBook(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		h.logger.Warning(err)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{"id value of incorrect type"})
	}

	if _, err = h.services.Users.GetByID(context.TODO(), userID); err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{"user not found"})
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	if err = h.services.Books.Delete(context.TODO(), id, userID); err != nil {
		if errors.Is(err, core.ErrBookNotFound) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary Get Books
// @Tags books
// @Description Get all book
// @ModuleID getAllBooks
// @Accept  json
// @Produce  json
// @Success 200 {object} []core.Book
// @Router /books [get]
func (h *Handler) getAllBooks(c *fiber.Ctx) error {
	books, err := h.services.Books.GetAll(context.TODO())
	if err != nil {
		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(books)
}

// @Summary Update Book
// @Tags books
// @Description update book
// @ModuleID updateBook
// @Security UsersAuth
// @Accept  json
// @Produce  json
// @Param id path string true "book id"
// @Param input body core.UpdateBookInput true "update book"
// @Success 200 {string} string "OK"
// @Failure 400 {object} response
// @Router /books/{id} [put]
func (h *Handler) updateBook(c *fiber.Ctx) error {
	userID, err := getUserID(c)
	if err != nil {
		h.logger.Warning(err)

		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{"id value of incorrect type"})
	}

	if _, err = h.services.Users.GetByID(context.TODO(), userID); err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{"user not found"})
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	var inp core.UpdateBookInput
	if err = c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := h.validateStruct(inp); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	if err = h.services.Books.Update(context.TODO(), id, userID, inp); err != nil {
		if errors.Is(err, core.ErrBookNotFound) {
			return c.SendStatus(fiber.StatusForbidden)
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
