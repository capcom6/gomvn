package server

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/gomvn/gomvn/internal/entity"
)

//	@Summary		List users
//	@Description	returns list of users
//	@Tags			Users
//	@Security		BasicAuth
//	@Produce		json
//	@Param			limit	query		int	false	"Limit on page"		Default(50)
//	@Param			offset	query		int	false	"Offset of page"	Default(0)
//	@Success		200		{object}	apiGetUsersResponse
//	@Failure		500		{object}	string
//	@Router			/api/users [get]
//
// List users.
func (s *Server) handleAPIGetUsers(c *fiber.Ctx) error {
	const defaultLimit = 50

	limit := getQueryInt(c, "limit", defaultLimit)
	offset := getQueryInt(c, "offset", 0)

	users, count, err := s.us.GetAll(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return c.JSON(&apiGetUsersResponse{
		Total: count,
		Items: mapToAPIGetUsersItem(users),
	})
}

func mapToAPIGetUsersItem(users []entity.User) []apiGetUsersItem {
	items := make([]apiGetUsersItem, len(users))
	for i, user := range users {
		items[i] = apiGetUsersItem{
			ID:        user.ID,
			Name:      user.Name,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Paths:     mapToAPIGetUsersPathItem(user.Paths),
		}
	}
	return items
}

func mapToAPIGetUsersPathItem(paths []entity.Path) []apiGetUsersPathItem {
	items := make([]apiGetUsersPathItem, len(paths))
	for i, path := range paths {
		items[i] = apiGetUsersPathItem{
			Path:      path.Path,
			Deploy:    path.Deploy,
			CreatedAt: path.CreatedAt,
			UpdatedAt: path.UpdatedAt,
		}
	}
	return items
}

type apiGetUsersResponse struct {
	Total int64             `json:"total"` // Total count of users
	Items []apiGetUsersItem `json:"items"` // List of users
}

type apiGetUsersItem struct {
	ID        uint                  `json:"id"`        // User ID
	Name      string                `json:"name"`      // User name
	CreatedAt time.Time             `json:"createdAt"` // User created at
	UpdatedAt time.Time             `json:"updatedAt"` // User updated at
	Paths     []apiGetUsersPathItem `json:"allowed"`   // List of allowed paths
}

type apiGetUsersPathItem struct {
	Path      string    `json:"name"`      // Path
	Deploy    bool      `json:"deploy"`    // Allowed to deploy
	CreatedAt time.Time `json:"createdAt"` // Path created at
	UpdatedAt time.Time `json:"updatedAt"` // Path updated at
}
