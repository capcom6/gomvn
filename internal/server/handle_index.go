package server

import (
	"github.com/gofiber/fiber/v2"
)

func (s *Server) handleIndex(c *fiber.Ctx) error {
	return c.Render("index", fiber.Map{
		"Name":         s.name,
		"Repositories": s.rs.GetRepositories(),
	})
}
