package repository

import (
	"avitostazhko/internal/infrastructure/repository/txs"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Options(
		fx.Provide(NewTeamRepository),
		fx.Provide(NewUserRepository),
		fx.Provide(NewPullRequestRepository),
		fx.Provide(NewReviewersRepository),
		fx.Provide(txs.NewTxBeginner),
	)
}
