package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleApiDeleteUsers(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := s.us.Delete(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}
