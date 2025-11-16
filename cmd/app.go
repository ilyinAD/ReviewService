package main

import (
	"avitostazhko/database"
	"avitostazhko/internal/infrastructure"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func BuildApp() fx.Option {
	return fx.Options(
		database.FxModule(),
		infrastructure.FxModule(),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)
}
