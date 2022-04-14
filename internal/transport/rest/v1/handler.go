package v1

import (
	"errors"

	"github.com/ernur-eskermes/crud-app/internal/service"
	"github.com/ernur-eskermes/crud-app/pkg/auth"
	"github.com/ernur-eskermes/crud-app/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
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

func (h *Handler) Init(api fiber.Router) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		h.initBooksRoutes(v1)
	}
}

type response struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       string
}

func (h *Handler) validateStruct(st interface{}) []*ErrorResponse {
	var res []*ErrorResponse

	if err := h.validate.Struct(st); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, err := range validationErrors {
				var element ErrorResponse
				element.FailedField = err.StructNamespace()
				element.Tag = err.Tag()
				element.Value = err.Param()
				res = append(res, &element)
			}
		}
	}

	return res
}
