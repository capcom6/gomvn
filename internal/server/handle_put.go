package server

import (
	"log" //nolint:depguard // TODO

	"github.com/gofiber/fiber/v2"
)

func (s *Server) handlePut(c *fiber.Ctx) error {
	path, err := s.ps.ParsePath(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if wrErr := s.storage.Write(path, c.Context().RequestBodyStream()); wrErr != nil {
		log.Printf("cannot put data: %v", wrErr)
		return c.Status(fiber.StatusBadRequest).SendString(wrErr.Error())
	}

	return c.SendStatus(fiber.StatusOK)
}
