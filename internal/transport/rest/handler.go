package rest

import (
	"context"
	"time"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/google/uuid"

	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/logging"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

type UsersService interface {
	SignUp(ctx context.Context, input core.AuthInput) error
	Verify(ctx context.Context, username, code string) error
	SignIn(ctx context.Context, input core.AuthInput) (core.Tokens, error)
	RefreshTokens(ctx context.Context, refreshToken string) (core.Tokens, error)
	GetByID(ctx context.Context, id uuid.UUID) (core.User, error)
}

type BooksService interface {
	Create(ctx context.Context, book core.CreateBookInput, userID uuid.UUID) (core.Book, error)
	GetByID(ctx context.Context, id uuid.UUID) (core.Book, error)
	GetAll(ctx context.Context) ([]core.Book, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Update(ctx context.Context, id, userID uuid.UUID, inp core.UpdateBookInput) error
}

type Handler struct {
	usersService UsersService
	booksService BooksService
	tokenManager auth.TokenManager
	logger       *logging.Logger
}

func NewHandler(usersService UsersService, booksService BooksService, tokenManager auth.TokenManager, logger *logging.Logger) *Handler {
	return &Handler{
		usersService: usersService,
		booksService: booksService,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (h *Handler) InitRouter(app *fiber.App) {
	app.Use(cors.New())
	app.Use(logger.New(logger.Config{
		TimeFormat: time.RFC3339,
		TimeZone:   "Asia/Almaty",
	}))
	app.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return c.IP() == "127.0.0.1"
		},
		Max:        20,
		Expiration: 30 * time.Second,
	}))

	app.Get("/dashboard", monitor.New())
	app.Get("/swagger/*", swagger.HandlerDefault)
	h.initAPI(app)
}

func (h *Handler) initAPI(app fiber.Router) {
	api := app.Group("/api")
	{
		h.initAuthRoutes(api)
		h.initBooksRoutes(api)
	}
}

type response struct {
	Message string `json:"message"`
}
