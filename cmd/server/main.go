package main

import (
	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/database"
	"github.com/felipeversiane/auth-service/internal/infra/logger"
	"github.com/felipeversiane/auth-service/internal/infra/server"
	"github.com/felipeversiane/auth-service/internal/infra/telemetry"

	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		config.Module,
		logger.Module,
		database.Module,
		telemetry.Module,
		server.Module,
		fx.NopLogger,
	)

	app.Run()
}
