package service

import (
	"context"
	"time"

	"github.com/ernur-eskermes/crud-app/pkg/logging"
	audit "github.com/ernur-eskermes/crud-audit-log/pkg/domain"

	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
)

type BooksRepository interface {
	Create(ctx context.Context, book *core.Book) error
	GetByID(ctx context.Context, id uuid.UUID) (core.Book, error)
	GetAll(ctx context.Context) ([]core.Book, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, inp core.Book) error
}

type BooksService struct {
	repo         BooksRepository
	tokenManager auth.TokenManager
	auditClient  AuditClient
	logger       *logging.Logger
}

func NewBooksService(repo BooksRepository, auditClient AuditClient, tokenManager auth.TokenManager, logger *logging.Logger) *BooksService {
	return &BooksService{
		repo:         repo,
		tokenManager: tokenManager,
		auditClient:  auditClient,
		logger:       logger,
	}
}

func (s *BooksService) Create(ctx context.Context, inp core.CreateBookInput, userID uuid.UUID) (core.Book, error) {
	if inp.PublishDate.IsZero() {
		inp.PublishDate = time.Now()
	}

	book := core.Book{
		Title:       inp.Title,
		Author:      userID,
		PublishDate: inp.PublishDate,
		Rating:      inp.Rating,
	}

	err := s.repo.Create(ctx, &book)
	if err != nil {
		return core.Book{}, err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionCreate,
		Entity:    audit.EntityBook,
		EntityID:  book.ID.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return book, nil
}

func (s *BooksService) GetByID(ctx context.Context, id uuid.UUID) (core.Book, error) {
	book, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return book, err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionGet,
		Entity:    audit.EntityBook,
		EntityID:  book.ID.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return book, nil
}

func (s *BooksService) GetAll(ctx context.Context) ([]core.Book, error) {
	return s.repo.GetAll(ctx)
}

func (s *BooksService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		return err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionDelete,
		Entity:    audit.EntityBook,
		EntityID:  id.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return nil
}

func (s *BooksService) Update(ctx context.Context, id, userID uuid.UUID, inp core.UpdateBookInput) error {
	err := s.repo.Update(ctx, core.Book{
		ID:          id,
		Title:       inp.Title,
		Author:      userID,
		PublishDate: inp.PublishDate,
		Rating:      inp.Rating,
	})
	if err != nil {
		return err
	}

	if err = s.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    audit.ActionUpdate,
		Entity:    audit.EntityBook,
		EntityID:  id.String(),
		Timestamp: time.Now(),
	}); err != nil {
		s.logger.Error("failed to send log request: ", err)
	}

	return nil
}
