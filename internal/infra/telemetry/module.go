package telemetry

import (
	"context"

	"github.com/felipeversiane/auth-service/internal/infra/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		func(config config.TelemetryConfig) (TelemetryProviderInterface, error) {
			return New(config)
		},
	),
	fx.Invoke(func(lc fx.Lifecycle, observer TelemetryProviderInterface) {
		lc.Append(fx.Hook{
			OnStop: func(ctx context.Context) error {
				return observer.Shutdown(ctx)
			},
		})
	}),
)
