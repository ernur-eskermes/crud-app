package v1

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/ernur-eskermes/crud-app/internal/core"
	"github.com/ernur-eskermes/crud-app/internal/service"
)

func (h *Handler) initAuthRoutes(api fiber.Router) {
	users := api.Group("/auth")
	{
		users.Post("/sign-up", h.userSignUp)
		users.Post("/sign-in", h.userSignIn)
		users.Post("/verify", h.verify)
	}
}

type userSignUpInput struct {
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type signInInput struct {
	Username string `json:"username" validate:"required,max=64"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type tokenResponse struct {
	AccessToken string `json:"accessToken"`
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body userSignUpInput true "sign up info"
// @Success 201 {string} string "Created"
// @Failure 400 {object} response
// @Router /auth/sign-up [post]
func (h *Handler) userSignUp(c *fiber.Ctx) error {
	var inp userSignUpInput
	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := h.validateStruct(inp); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	if err := h.services.Users.SignUp(c.Context(), service.UserSignUpInput{
		Username: inp.Username,
		Password: inp.Password,
	}); err != nil {
		if errors.Is(err, core.ErrUserAlreadyExists) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusCreated)
}

// @Summary User SignIn
// @Tags users-auth
// @Description user sign in
// @ModuleID userSignIn
// @Accept  json
// @Produce  json
// @Param input body signInInput true "sign up info"
// @Success 200 {object} tokenResponse
// @Failure 400 {object} response
// @Router /auth/sign-in [post]
func (h *Handler) userSignIn(c *fiber.Ctx) error {
	var inp signInInput

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := h.validateStruct(inp); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	res, err := h.services.Users.SignIn(c.Context(), service.UserSignInInput{
		Username: inp.Username,
		Password: inp.Password,
	})
	if err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(tokenResponse{
		AccessToken: res.AccessToken,
	})
}

type verifyInput struct {
	Username string `json:"username"`
	Code     string `json:"code"`
}

// @Summary User Verify
// @Tags users-auth
// @Description user verify
// @ModuleID verify
// @Accept  json
// @Produce  json
// @Param input body verifyInput true "verify"
// @Success 200
// @Failure 400 {object} response
// @Router /auth/verify [post]
func (h *Handler) verify(c *fiber.Ctx) error {
	var inp verifyInput

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := h.validateStruct(inp); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	if err := h.services.Users.Verify(c.Context(), inp.Username, inp.Code); err != nil {
		if errors.Is(err, core.ErrUserNotFound) ||
			errors.Is(err, core.ErrUserCodeExpired) ||
			errors.Is(err, core.ErrUserCodeUnknownType) ||
			errors.Is(err, core.ErrUserCodeIncorrect) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
