package main

import (
	"github.com/felipeversiane/auth-service/internal/infra/config"
	"github.com/felipeversiane/auth-service/internal/infra/database"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module,
		database.Module,
		fx.NopLogger,
	)

	app.Run()
}
