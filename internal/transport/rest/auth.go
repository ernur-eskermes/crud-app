package rest

import (
	"errors"

	"github.com/gofiber/fiber/v2"

	"github.com/ernur-eskermes/crud-app/internal/core"
)

func (h *Handler) initAuthRoutes(api fiber.Router) {
	users := api.Group("/auth")
	{
		users.Post("/sign-up", h.userSignUp)
		users.Post("/sign-in", h.userSignIn)
		users.Post("/verify", h.verify)
	}
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID userSignUp
// @Accept  json
// @Produce  json
// @Param input body core.AuthInput true "sign up info"
// @Success 201 {string} string "Created"
// @Failure 400 {object} response
// @Router /auth/sign-up [post]
func (h *Handler) userSignUp(c *fiber.Ctx) error {
	var inp core.AuthInput
	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := inp.Validate(); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	if err := h.usersService.SignUp(c.Context(), inp); err != nil {
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
// @Param input body core.AuthInput true "sign up info"
// @Success 200 {object} core.Tokens
// @Failure 400 {object} response
// @Router /auth/sign-in [post]
func (h *Handler) userSignIn(c *fiber.Ctx) error {
	var inp core.AuthInput

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := inp.Validate(); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	res, err := h.usersService.SignIn(c.Context(), inp)
	if err != nil {
		if errors.Is(err, core.ErrUserNotFound) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(res)
}

// @Summary User Verify
// @Tags users-auth
// @Description user verify
// @ModuleID verify
// @Accept  json
// @Produce  json
// @Param input body core.VerifyInput true "verify"
// @Success 200
// @Failure 400 {object} response
// @Router /auth/verify [post]
func (h *Handler) verify(c *fiber.Ctx) error {
	var inp core.VerifyInput

	if err := c.BodyParser(&inp); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
	}

	if validationError := inp.Validate(); validationError != nil {
		return c.Status(fiber.StatusBadRequest).JSON(validationError)
	}

	if err := h.usersService.Verify(c.Context(), inp.Username, inp.Code); err != nil {
		if errors.Is(err, core.ErrUserNotFound) || errors.Is(err, core.ErrUserCodeIncorrect) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
