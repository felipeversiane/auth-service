package http

import (
	"context"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/database"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func(config config.HttpServerConfig, db database.DatabaseInterface) HttpServerInterface {
			return New(config, db)
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, server HttpServerInterface) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					server.InitRoutes()
					if err := server.Start(); err != nil {
						panic(err)
					}
				}()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return server.Shutdown(ctx)
			},
		})
	}),
)
