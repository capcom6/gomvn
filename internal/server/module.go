package server

import (
	"context"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(New),
		fx.Invoke(register),
	)
}

func register(lifecycle fx.Lifecycle, server *Server) {
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return server.Listen()
		},
		OnStop: func(ctx context.Context) error {
			return server.app.ShutdownWithContext(ctx)
		},
	})
}
