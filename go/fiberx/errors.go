package fiberx

import "github.com/gofiber/fiber/v2"

const (
	ErrMsg400 = "the server could not understand the request due to invalid syntax"
	ErrMsg401 = "the server could not verify that you are authorized to access the requested resource"
	ErrMsg403 = "you do not have to sufficient role to access the requested resource"
	ErrMsg404 = "the server could not find the requested resource"
	ErrMsg500 = "the server encountered an error and could not complete your request"
)

var lookup = map[int]string{
	fiber.StatusInternalServerError: ErrMsg500,
	fiber.StatusBadRequest:          ErrMsg400,
	fiber.StatusUnauthorized:        ErrMsg401,
	fiber.StatusForbidden:           ErrMsg403,
	fiber.StatusNotFound:            ErrMsg404,
}

func Err(c *fiber.Ctx, code int, customMsg ...string) error {
	msg := lookup[code]
	if len(customMsg) > 0 {
		msg = customMsg[0]
	}
	return c.Status(code).JSON(fiber.Map{
		"error": msg,
		"code":  code,
	})
}
