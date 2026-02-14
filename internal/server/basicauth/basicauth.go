package basicauth

import (
	"encoding/base64"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Config defines the config for BasicAuth middleware.
type Config struct {
	// Next defines a function to skip middleware.
	// Optional. Default: nil
	Next func(*fiber.Ctx) bool
	// Authorizer defines a function you can pass
	// to check the credentials however you want.
	// It will be called with a username and password
	// and is expected to return true or false to indicate
	// that the credentials were approved or not.
	// Optional. Default: nil.
	Authorizer func(*fiber.Ctx, string, string) bool
}

func New(config ...Config) func(*fiber.Ctx) error {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}
	// Return middleware handler
	return func(c *fiber.Ctx) error {
		// Filter request to skip middleware
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		// Get authorization header
		auth := c.Get(fiber.HeaderAuthorization)

		// Check if header is valid
		if len(auth) < 6 || strings.ToLower(auth[:6]) != "basic " {
			return unauthorized(c)
		}

		// Try to decode
		raw, err := base64.StdEncoding.DecodeString(auth[6:])
		if err != nil {
			return unauthorized(c)
		}

		// Convert to string
		cred := string(raw)

		// Split on first colon (password may contain colons per RFC 7617)
		if before, after, ok := strings.Cut(cred, ":"); ok {
			user := before
			pass := after
			if cfg.Authorizer != nil && cfg.Authorizer(c, user, pass) {
				return c.Next()
			}
		}

		// Authentication failed
		return unauthorized(c)
	}
}

func unauthorized(c *fiber.Ctx) error {
	c.Set(fiber.HeaderWWWAuthenticate, `basic realm="Restricted"`)
	return fiber.ErrUnauthorized
}
