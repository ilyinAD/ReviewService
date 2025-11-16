package usecases

import (
	"avitostazhko/internal/infrastructure/usecases/pullrequestusecase"
	"avitostazhko/internal/infrastructure/usecases/teamusecase"
	"avitostazhko/internal/infrastructure/usecases/userusecase"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Options(
		fx.Provide(teamusecase.NewTeamUseCase),
		fx.Provide(userusecase.NewUserUseCase),
		fx.Provide(pullrequestusecase.NewPullRequestUseCase),
	)
}
