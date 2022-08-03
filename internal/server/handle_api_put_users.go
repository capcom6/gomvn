package server

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// @Summary      Update user
// @Description  updates single user without changing token
// @Tags         Users
// @Security     BasicAuth
// @Produce      json
// @Param        id    path      int                 true  "User ID"
// @Param        user  body      apiPutUsersRequest  true  "Edited user"
// @Success      200   {object}  apiGetUsersResponse
// @Failure      400   {object}  string
// @Failure      500   {object}  string
// @Router       /api/users/{id} [put]
func (s *Server) handleApiPutUsers(c *fiber.Ctx) error {
	r := new(apiPutUsersRequest)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := r.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, err := s.us.Update(uint(id), r.Deploy, r.Allowed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(&apiPutUsersResponse{
		Id:   user.ID,
		Name: user.Name,
	})
}

type apiPutUsersRequest struct {
	Deploy  bool     `json:"deploy"`  // Is alowed to deploy
	Allowed []string `json:"allowed"` // Allowed paths
}

func (r *apiPutUsersRequest) validate() error {
	if len(r.Allowed) < 1 {
		return fmt.Errorf("field 'allowed' must contain at least one string")
	}
	for _, path := range r.Allowed {
		if path[0] != '/' {
			return fmt.Errorf("paths in field 'allowed' must start with '/'")
		}
	}
	return nil
}

type apiPutUsersResponse struct {
	Id   uint   `json:"id"`   // User ID
	Name string `json:"name"` // User name
}
