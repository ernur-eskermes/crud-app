package rest

import (
	"time"

	"github.com/ernur-eskermes/crud-app/internal/config"
	"github.com/ernur-eskermes/crud-app/internal/service"
	v1 "github.com/ernur-eskermes/crud-app/internal/transport/rest/v1"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/monitor"

	swagger "github.com/arsmn/fiber-swagger/v2"
)

type Handler struct {
	services     *service.Services
	tokenManager auth.TokenManager
	validate     *validator.Validate
	logger       *logging.Logger
}

func NewHandler(services *service.Services, tokenManager auth.TokenManager, validate *validator.Validate, logger *logging.Logger) *Handler {
	return &Handler{
		services:     services,
		validate:     validate,
		tokenManager: tokenManager,
		logger:       logger,
	}
}

func (h *Handler) InitRouter(app *fiber.App, cfg *config.Config) {
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
	handlerV1 := v1.NewHandler(h.services, h.tokenManager, h.validate, h.logger)
	api := app.Group("/api")
	{
		handlerV1.Init(api)
	}
}
