package server

import (
	"fmt"
	"log" //nolint:depguard // TODO
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"

	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/server/middleware"
	"github.com/gomvn/gomvn/internal/service"
	"github.com/gomvn/gomvn/internal/service/storage"
	"github.com/gomvn/gomvn/internal/service/users"
)

func New(
	conf *config.App,
	ps *service.PathService,
	storage *storage.Storage,
	us *users.Service,
	rs *service.RepoService,
) *Server {
	const defaultTimeout = time.Minute

	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		IdleTimeout:           defaultTimeout,
		DisableStartupMessage: true,
		Views:                 engine,
		StreamRequestBody:     true,
	})

	app.Use(recover.New())
	app.Use(compress.New())
	app.Use(logger.New())

	server := &Server{
		app:     app,
		name:    conf.Name,
		listen:  conf.Server.GetListenAddr(),
		ps:      ps,
		storage: storage,
		us:      us,
		rs:      rs,
	}

	registerAPI(app, us, server)

	app.Static("/admin", "./views/admin")

	app.Put("/*", middleware.NewRepoAuth(us, ps, true), server.handlePut)

	if *us.GetDefaultPermissions().Index {
		app.Get("/", server.handleIndex)
	}

	app.Use(middleware.NewRepoAuth(us, ps, false))

	if !*us.GetDefaultPermissions().Index {
		app.Get("/", server.handleIndex)
	}

	app.Get("/+", func(c *fiber.Ctx) error {
		pathname := c.Params("+")
		file, contentType, err := storage.Open(pathname)
		if err != nil {
			if os.IsNotExist(err) {
				return fiber.ErrNotFound
			}
			log.Printf("failed to open file at %s: %v", pathname, err)
			return fiber.ErrInternalServerError
		}

		c.Set(fiber.HeaderContentType, contentType)
		return c.SendStream(file)
	})

	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return server
}

func registerAPI(app *fiber.App, us *users.Service, server *Server) {
	api := app.Group("/api")
	api.Use(middleware.NewAPIAuth(us))
	api.Get("/users", server.handleAPIGetUsers)
	api.Post("/users", server.handleAPIPostUsers)
	api.Put("/users/:id", server.handleAPIPutUsers)
	api.Delete("/users/:id", server.handleAPIDeleteUsers)

	api.Get("/users/:id/refresh", server.handleAPIGetUsersRefresh)

	api.Put("/users/:id/paths", server.handleAPIPutUsersPaths)
}

type Server struct {
	app     *fiber.App
	name    string
	listen  string
	ps      *service.PathService
	storage *storage.Storage
	us      *users.Service
	rs      *service.RepoService
}

func (s *Server) Listen() error {
	log.Printf("GoMVN self-hosted repository listening on %s\n", s.listen)
	go func() {
		if err := s.app.Listen(s.listen); err != nil {
			log.Fatalf("server listen error: %v", err)
		}
	}()
	return nil
}

func (s *Server) Shutdown() error {
	if err := s.app.Shutdown(); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}
	return nil
}
