package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary      Delete user
// @Description  deletes user by id
// @Tags         Users
// @Security     BasicAuth
// @Produce      text/plain
// @Param        id   path      int  true  "User ID"
// @Success      203  {object}  string
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /api/users/{id} [delete]
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
