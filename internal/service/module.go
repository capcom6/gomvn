package service

import (
	"go.uber.org/fx"

	"github.com/gomvn/gomvn/internal/service/storage"
	"github.com/gomvn/gomvn/internal/service/users"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(NewPathService),
		fx.Provide(storage.NewStorage),
		fx.Provide(NewRepoService),
		fx.Provide(users.New),
		fx.Invoke(users.Initialize),
	)
}
