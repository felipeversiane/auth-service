package logger

import (
	"context"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Options(
	fx.Provide(
		func(cfg config.LogConfig) *zap.Logger {
			return New(cfg)
		},
	),
	fx.Invoke(func(lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				StopFlush()
				return nil
			},
		})
	}),
)
