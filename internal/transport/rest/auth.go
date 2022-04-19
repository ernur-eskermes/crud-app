package rest

import (
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/ernur-eskermes/crud-app/internal/core"
)

func (h *Handler) initAuthRoutes(api fiber.Router) {
	users := api.Group("/auth")
	{
		users.Post("/sign-up", h.signUp)
		users.Post("/sign-in", h.signIn)
		users.Post("/refresh", h.refresh)
		users.Post("/verify", h.verify)
	}
}

// @Summary User SignUp
// @Tags users-auth
// @Description create user account
// @ModuleID signUp
// @Accept  json
// @Produce  json
// @Param input body core.AuthInput true "sign up info"
// @Success 201 {string} string "Created"
// @Failure 400 {object} response
// @Router /auth/sign-up [post]
func (h *Handler) signUp(c *fiber.Ctx) error {
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
// @ModuleID signIn
// @Accept  json
// @Produce  json
// @Param input body core.AuthInput true "sign up info"
// @Success 200 {object} core.Tokens
// @Failure 400 {object} response
// @Router /auth/sign-in [post]
func (h *Handler) signIn(c *fiber.Ctx) error {
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

	c.Set("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", res.RefreshToken))

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

// @Summary Refresh Token
// @Tags users-auth
// @Description refresh token
// @ModuleID refresh
// @Accept  json
// @Produce  json
// @Success 200 {object} core.Tokens
// @Failure 400 {object} response
// @Router /auth/refresh [post]
func (h *Handler) refresh(c *fiber.Ctx) error {
	cookie := c.Cookies("refresh-token", "")
	if cookie == "" {
		return c.Status(fiber.StatusBadRequest).JSON(response{"refresh token is empty"})
	}

	res, err := h.usersService.RefreshTokens(c.Context(), cookie)
	if err != nil {
		if errors.Is(err, core.ErrTokenNotFound) || errors.Is(err, core.ErrRefreshTokenExpired) {
			return c.Status(fiber.StatusBadRequest).JSON(response{err.Error()})
		}

		h.logger.Error(err)

		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Set("Set-Cookie", fmt.Sprintf("refresh-token=%s; HttpOnly", res.RefreshToken))

	return c.JSON(res)
}
