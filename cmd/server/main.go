package main

import (
	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/database"
	"github.com/felipeversiane/auth-service/internal/infra/http"
	"github.com/felipeversiane/auth-service/internal/infra/telemetry"

	"go.uber.org/fx"
)

func main() {

	app := fx.New(
		config.Module,
		database.Module,
		telemetry.Module,
		http.Module,
		fx.NopLogger,
	)

	app.Run()
}
