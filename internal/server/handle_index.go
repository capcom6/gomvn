package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleIndex(c *fiber.Ctx) error {
	err := c.Render("index", fiber.Map{
		"Name":         s.name,
		"Repositories": s.rs.GetRepositories(),
	})
	if err != nil {
		// TODO: log err for debugging
		return c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
	}
	return nil
}
