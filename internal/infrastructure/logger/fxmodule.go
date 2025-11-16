package logger

import "go.uber.org/fx"

func FxModule() fx.Option {
	return fx.Options(
		fx.Provide(NewLogger),
	)
}
