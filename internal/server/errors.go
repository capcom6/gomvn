package server

import "github.com/gofiber/fiber/v2"

var (
	ErrValidationFailed = fiber.NewError(
		fiber.StatusBadRequest,
		"validation failed",
	)
)
