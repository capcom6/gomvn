package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// @Summary      Create user
// @Description  creates new user and returns access token
// @Tags         Users
// @Security     BasicAuth
// @Produce      json
// @Param        user  body      apiPostUsersRequest  true  "New user"
// @Success      200   {object}  apiPostUsersResponse
// @Failure      400   {object}  string
// @Failure      500   {object}  string
// @Router       /api/users [post]
func (s *Server) handleApiPostUsers(c *fiber.Ctx) error {
	r := new(apiPostUsersRequest)
	if err := c.BodyParser(r); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	if err := r.validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	user, token, err := s.us.Create(r.Name, r.Admin, r.Deploy, r.Allowed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(&apiPostUsersResponse{
		Id:    user.ID,
		Name:  user.Name,
		Token: token,
	})
}

type apiPostUsersRequest struct {
	Name    string   `json:"name"`    // User name
	Admin   bool     `json:"admin"`   // Is admin user
	Deploy  bool     `json:"deploy"`  // Is allowed to deploy
	Allowed []string `json:"allowed"` // Allowed paths
}

func (r *apiPostUsersRequest) validate() error {
	if r.Name == "" {
		return fmt.Errorf("field 'name' cannot be empty")
	}
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

type apiPostUsersResponse struct {
	Id    uint   `json:"id"`    // User ID
	Name  string `json:"name"`  // User name
	Token string `json:"token"` // Access token
}
