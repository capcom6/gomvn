// Copyright 2022 Aleksandr Soloshenko
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gomvn/gomvn/internal/entity"
)

// @Summary      Replace user's allowed paths
// @Description  replaces user's allowed paths
// @Tags         Users,Paths
// @Security     BasicAuth
// @Produce      json
// @Param        id     path      int                             true  "User ID"
// @Param        paths  body      []apiPuthUsersPathsRequestItem  true  "Allowed paths"
// @Success      200    {array}   apiPuthUsersPathsResponseItem   "Current allowed paths"
// @Failure      400    {object}  string
// @Failure      401    {object}  string
// @Failure      500    {object}  string
// @Router       /api/users/{id}/paths [put]
func (s *Server) handleApiPutUsersPaths(c *fiber.Ctx) error {
	paths := make([]*apiPuthUsersPathsRequestItem, 0)
	if err := c.BodyParser(&paths); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if len(paths) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No paths provided")
	}

	newPaths := make([]entity.Path, 0)
	for _, path := range paths {
		if path.Path[0] != '/' {
			return c.Status(fiber.StatusBadRequest).SendString("Paths must start with '/'")
		}
		newPaths = append(newPaths, entity.Path{
			Path:   path.Path,
			Deploy: path.Deploy,
		})
	}

	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if newPaths, err = s.us.ReplacePaths(uint(id), newPaths); err != nil {
		log.Printf("[ERROR] %v", err)
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	resp := make([]*apiPuthUsersPathsResponseItem, len(newPaths))
	for i, path := range newPaths {
		resp[i] = &apiPuthUsersPathsResponseItem{
			Path:   path.Path,
			Deploy: path.Deploy,
		}
	}

	return c.JSON(resp)
}

type apiPuthUsersPathsRequestItem struct {
	Path   string `json:"path"`
	Deploy bool   `json:"deploy"`
}

type apiPuthUsersPathsResponseItem struct {
	Path   string `json:"path"`
	Deploy bool   `json:"deploy"`
}
