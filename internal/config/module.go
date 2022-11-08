package config

import (
	"go.uber.org/fx"
)

func Module(configFile string) fx.Option {
	return fx.Options(
		fx.Provide(func() (*App, error) {
			return NewAppConfig(configFile)
		}),
		fx.Provide(provideServerConfig),
		fx.Provide(provideStorageConfig),
	)
}

func provideServerConfig(app *App) *Server {
	return app.Server
}

func provideStorageConfig(app *App) *Storage {
	return &app.Storage
}
