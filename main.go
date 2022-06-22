package main

import (
	"flag"

	"go.uber.org/fx"

	"github.com/gomvn/gomvn/internal/config"
	"github.com/gomvn/gomvn/internal/database"
	"github.com/gomvn/gomvn/internal/server"
	"github.com/gomvn/gomvn/internal/service"
)

// @title        Self-Hoster Maven Repository
// @version      1.0.0
// @description  Management API

// @contact.name   Aleksandr Soloshenko
// @contact.email  capcom@soft-c.ru
// @license.name   MIT
// @license.url    https://github.com/capcom6/gomvn/blob/master/LICENSE

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.basic  BasicAuth
func main() {
	cf := flag.String("config", "config.yml", "path to config file")
	flag.Parse()

	app := fx.New(
		fx.NopLogger,
		config.Module(*cf),
		database.Module,
		service.Module,
		server.Module,
	)
	app.Run()
}
