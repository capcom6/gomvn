package server

import (
	"log" //nolint:depguard // TODO
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gomvn/gomvn/internal/entity"
)

//	@Summary		Replace user's allowed paths
//	@Description	replaces user's allowed paths
//	@Tags			Users,Paths
//	@Security		BasicAuth
//	@Produce		json
//	@Param			id		path		int								true	"User ID"
//	@Param			paths	body		[]apiPutUsersPathsRequestItem	true	"Allowed paths"
//	@Success		200		{array}		apiPutUsersPathsResponseItem	"Current allowed paths"
//	@Failure		400		{object}	string
//	@Failure		401		{object}	string
//	@Failure		500		{object}	string
//	@Router			/api/users/{id}/paths [put]
//
// Replace user's allowed paths.
func (s *Server) handleAPIPutUsersPaths(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	paths := make([]*apiPutUsersPathsRequestItem, 0)
	if bodyErr := c.BodyParser(&paths); bodyErr != nil {
		return c.Status(fiber.StatusBadRequest).SendString(bodyErr.Error())
	}

	if len(paths) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No paths provided")
	}

	newPaths := make([]entity.Path, 0)
	for _, path := range paths {
		if len(path.Path) == 0 || path.Path[0] != '/' {
			return c.Status(fiber.StatusBadRequest).SendString("Paths must be non-empty and start with '/'")
		}
		newPaths = append(newPaths, entity.NewPath(uint(id), path.Path, path.Deploy))
	}

	if newPaths, err = s.us.ReplacePaths(uint(id), newPaths); err != nil {
		log.Printf("[ERROR] %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := make([]*apiPutUsersPathsResponseItem, len(newPaths))
	for i, path := range newPaths {
		resp[i] = &apiPutUsersPathsResponseItem{
			Path:   path.Path,
			Deploy: path.Deploy,
		}
	}

	return c.JSON(resp)
}

type apiPutUsersPathsRequestItem struct {
	Path   string `json:"name"`   // Path
	Deploy bool   `json:"deploy"` // Allowed to deploy
}

type apiPutUsersPathsResponseItem struct {
	Path   string `json:"name"`   // Path
	Deploy bool   `json:"deploy"` // Allowed to deploy
}
