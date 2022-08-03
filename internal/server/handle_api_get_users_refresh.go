package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary      Refreshes user token
// @Description  regenerates user access token
// @Tags         Users
// @Security     BasicAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  apiGetUsersTokenResponse
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /api/users/{id}/refresh [get]
func (s *Server) handleApiGetUsersRefresh(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, token, err := s.us.UpdateToken(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(&apiGetUsersTokenResponse{
		Id:    user.ID,
		Name:  user.Name,
		Token: token,
	})
}

type apiGetUsersTokenResponse struct {
	Id    uint   `json:"id"`    // User ID
	Name  string `json:"name"`  // User name
	Token string `json:"token"` // User new access token
}
