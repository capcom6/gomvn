package server

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func getQueryUint64(c *fiber.Ctx, name string, def uint64) uint64 {
	if val, err := strconv.ParseUint(c.Query(name), 10, 64); err == nil {
		return val
	}
	return def
}

func getQueryInt64(c *fiber.Ctx, name string, def int64) int64 {
	if val, err := strconv.ParseInt(c.Query(name), 10, 64); err == nil {
		return val
	}
	return def
}

func getQueryInt(c *fiber.Ctx, name string, def int) int {
	if val, err := strconv.ParseInt(c.Query(name), 10, 32); err == nil {
		return int(val)
	}
	return def
}
