package auth

import (
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"siransbach/taskmanagementapi/fiberx"
)

func NewMiddleware(pg *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get(fiber.HeaderAuthorization)
		if len(auth) == 0 {
			log.Error().Msg("error authenticating user: no auth header")
			return fiberx.Err(c, fiber.StatusUnauthorized)
		}
		user, err := authenticate(c.Context(), auth, pg)
		if err != nil {
			log.Err(err).Msg("error authenticating user")
			return fiberx.Err(c, fiber.StatusUnauthorized)
		}
		c.Locals("user", user)
		return c.Next()
	}
}

func CurrentUser(c *fiber.Ctx) (*User, error) {
	if v := c.Locals("user"); v != nil {
		return v.(*User), nil
	}
	return nil, errors.New("no user found")
}

func authenticate(ctx context.Context, auth string, pg *sql.DB) (*User, error) {
	if len(auth) > 6 && strings.ToLower(auth[:5]) == "basic" {
		raw, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return nil, err
		}
		// convert to string
		cred := string(raw)
		// find semicolon
		for i := 0; i < len(cred); i++ {
			if cred[i] == ':' {
				// split into user & pass
				username := cred[:i]
				password := cred[i+1:]
				// if exist & match in users, we let him pass
				user, err := NewDB(pg).FindOne(ctx, FindOptions{
					Username: username,
				})
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
					return nil, err
				}
				if err := bcrypt.CompareHashAndPassword(
					[]byte(user.EncryptedPassword), []byte(password)); err != nil {
					return nil, err
				}
				return user, nil
			}
		}
	}
	return nil, errors.New("invalid credentials")
}
