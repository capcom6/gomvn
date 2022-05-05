package server

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

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
	Deploy  bool     `json:"deploy"`
	Allowed []string `json:"allowed"`
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
	Id   uint   `json:"id"`
	Name string `json:"name"`
}
