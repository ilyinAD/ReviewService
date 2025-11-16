package database

import (
	"avitostazhko/migrations"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Options(
		fx.Provide(NewDBConfigFromEnv),
		fx.Provide(NewPgxPoolConfig),
		fx.Provide(NewPgxPool),
		migrations.FxModule(),
	)
}
