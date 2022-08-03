package server

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/template/html"

	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/server/middleware"
	"github.com/gomvn/gomvn/internal/service"
	"github.com/gomvn/gomvn/internal/service/user"
)

func New(conf *config.App, ps *service.PathService, storage *service.Storage, us *user.Service, rs *service.RepoService) *Server {
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		IdleTimeout:           5 * time.Second,
		DisableStartupMessage: true,
		Views:                 engine,
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

	registerApi(app, us, server)

	app.Put("/*", middleware.NewRepoAuth(us, ps, true), server.handlePut)

	if *us.GetDefaultPermissions().Index {
		app.Get("/", server.handleIndex)
	}

	app.Use(middleware.NewRepoAuth(us, ps, false))

	if !*us.GetDefaultPermissions().Index {
		app.Get("/", server.handleIndex)
	}

	app.Static("/", storage.GetRoot(), fiber.Static{
		Browse: true,
	})
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})

	return server
}

func registerApi(app *fiber.App, us *user.Service, server *Server) {
	api := app.Group("/api")
	api.Use(middleware.NewApiAuth(us))
	api.Get("/users", server.handleApiGetUsers)
	api.Post("/users", server.handleApiPostUsers)
	api.Put("/users/:id", server.handleApiPutUsers)
	api.Delete("/users/:id", server.handleApiDeleteUsers)

	api.Get("/users/:id/refresh", server.handleApiGetUsersRefresh)

	api.Put("/users/:id/paths", server.handleApiPutUsersPaths)
}

type Server struct {
	app     *fiber.App
	name    string
	listen  string
	ps      *service.PathService
	storage *service.Storage
	us      *user.Service
	rs      *service.RepoService
}

func (s *Server) Listen() error {
	log.Printf("GoMVN self-hosted repository listening on %s\n", s.listen)
	go s.app.Listen(s.listen)
	return nil
}

func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
