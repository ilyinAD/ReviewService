package serviceapi

import (
	"avitostazhko/internal/infrastructure/serviceapi/handlers"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Options(
		fx.Provide(NewRouter),
		fx.Provide(NewConfig),
		handlers.FxModule(),
		fx.Invoke(RunServer),
	)
}
