package middleware

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/gomvn/gomvn/internal/server/basicauth"
	"github.com/gomvn/gomvn/internal/service/users"
)

func NewAPIAuth(us *users.Service) func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Next: nil,
		Authorizer: func(_ *fiber.Ctx, name string, token string) bool {
			u, err := us.GetByName(name)
			if err != nil || !u.Admin {
				return false
			}
			if hashErr := bcrypt.CompareHashAndPassword([]byte(u.TokenHash), []byte(token)); hashErr != nil {
				return false
			}
			return true
		},
	})
}
