package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getQueryInt(c *fiber.Ctx, name string, def int) int {
	if val, err := strconv.ParseInt(c.Query(name), 10, 32); err == nil {
		return int(val)
	}
	return def
}
