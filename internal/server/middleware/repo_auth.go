package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"

	"github.com/gomvn/gomvn/internal/server/basicauth"
	"github.com/gomvn/gomvn/internal/service"
	"github.com/gomvn/gomvn/internal/service/users"
)

func NewRepoAuth(us *users.Service, ps *service.PathService, needsDeploy bool) func(*fiber.Ctx) error {
	return basicauth.New(basicauth.Config{
		Next: func(_ *fiber.Ctx) bool {
			permissions := us.GetDefaultPermissions()
			if needsDeploy {
				return permissions.Deploy
			}
			return permissions.View
		},
		Authorizer: func(c *fiber.Ctx, name string, token string) bool {
			u, err := us.GetByName(name)
			if err != nil {
				return false
			}
			if hashErr := bcrypt.CompareHashAndPassword([]byte(u.TokenHash), []byte(token)); hashErr != nil {
				return false
			}

			_, current, _, err := ps.ParsePathParts(c)
			if err != nil {
				if needsDeploy {
					return false
				}
				// Non-deploy request with unparseable path — fall through to path matching
			}

			paths, err := us.GetPaths(u)
			if err != nil {
				return false
			}

			current = "/" + current
			for _, path := range paths {
				if strings.HasPrefix(current, path.Path) && (path.Deploy || !needsDeploy) {
					return true
				}
			}

			return false
		},
	})
}
