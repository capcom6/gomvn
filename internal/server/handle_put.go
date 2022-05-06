package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handlePut(c *fiber.Ctx) error {
	path, err := s.ps.ParsePath(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if err := s.storage.WriteFromRequest(c, path); err != nil {
		log.Printf("cannot put data: %v", err)
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
