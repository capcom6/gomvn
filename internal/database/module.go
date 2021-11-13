package database

import (
	"context"

	"go.uber.org/fx"
	"gorm.io/gorm"
)

var Module = fx.Options(
	fx.Provide(New),
	fx.Invoke(register),
)

func register(lifecycle fx.Lifecycle, db *gorm.DB) {
	lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return nil //db.Close()
		},
	})
}
