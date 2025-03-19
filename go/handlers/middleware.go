package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"siransbach/taskmanagementapi/auth"
	"siransbach/taskmanagementapi/fiberx"
)

func userMustHaveRole(role auth.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := auth.CurrentUser(c)
		if err != nil {
			log.Err(err).Msg("error getting current user")
			return fiberx.Err(c, fiber.StatusUnauthorized)
		}
		if user.Role != role {
			log.Error().Msg("user role mismatch")
			return fiberx.Err(c, fiber.StatusForbidden)
		}
		return c.Next()
	}
}
