package v1

import (
	"errors"
	"strings"

	"github.com/google/uuid"

	"github.com/gofiber/fiber/v2"
)

const (
	authorizationHeader = "Authorization"

	userCtx = "userID"
)

func (h *Handler) userIdentity(c *fiber.Ctx) error {
	id, err := h.parseAuthHeader(c.Get(authorizationHeader))
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	userID, err := uuid.Parse(id)
	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Locals(userCtx, userID)

	return c.Next()
}

func (h *Handler) parseAuthHeader(header string) (string, error) {
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return h.tokenManager.Parse(headerParts[1])
}

func getUserID(c *fiber.Ctx) (uuid.UUID, error) {
	idFromCtx := c.Locals(userCtx)
	if idFromCtx == "" {
		return uuid.UUID{}, errors.New("userCtx not found")
	}

	idStr, ok := idFromCtx.(uuid.UUID)
	if !ok {
		return uuid.UUID{}, errors.New("userCtx is of invalid type")
	}

	return idStr, nil
}
