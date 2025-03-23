package database

import (
	"context"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func(config config.DatabaseConfig) (DatabaseInterface, error) {
			return New(config)
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, db DatabaseInterface) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				db.Close()
				return nil
			},
		})
	}),
)
