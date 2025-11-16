package infrastructure

import (
	"avitostazhko/internal/infrastructure/logger"
	"avitostazhko/internal/infrastructure/repository"
	"avitostazhko/internal/infrastructure/serviceapi"
	"avitostazhko/internal/infrastructure/usecases"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Options(
		serviceapi.FxModule(),
		logger.FxModule(),
		repository.FxModule(),
		usecases.FxModule(),
	)
}
