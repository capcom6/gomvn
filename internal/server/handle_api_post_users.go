package server

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

//	@Summary		Create user
//	@Description	creates new user and returns access token
//	@Tags			Users
//	@Security		BasicAuth
//	@Produce		json
//	@Param			user	body		apiPostUsersRequest	true	"New user"
//	@Success		200		{object}	apiPostUsersResponse
//	@Failure		400		{object}	string
//	@Failure		500		{object}	string
//	@Router			/api/users [post]
//
// Create user.
func (s *Server) handleAPIPostUsers(c *fiber.Ctx) error {
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
		ID:    user.ID,
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
		return fmt.Errorf("%w: field 'name' cannot be empty", ErrValidationFailed)
	}
	if len(r.Allowed) < 1 {
		return fmt.Errorf("%w: field 'allowed' must contain at least one string", ErrValidationFailed)
	}
	for _, path := range r.Allowed {
		if len(path) == 0 {
			return fmt.Errorf("%w: paths in field 'allowed' must not be empty", ErrValidationFailed)
		}
		if path[0] != '/' {
			return fmt.Errorf("%w: paths in field 'allowed' must start with '/'", ErrValidationFailed)
		}
	}
	return nil
}

type apiPostUsersResponse struct {
	ID    uint   `json:"id"`    // User ID
	Name  string `json:"name"`  // User name
	Token string `json:"token"` // Access token
}
