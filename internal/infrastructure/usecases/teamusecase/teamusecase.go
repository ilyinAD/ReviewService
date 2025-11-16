package teamusecase

import (
	"avitostazhko/internal/domain"
	"avitostazhko/internal/infrastructure/repository"
	"avitostazhko/internal/infrastructure/repository/txs"
	"context"
	"fmt"
	"time"
)

type TeamUseCase struct {
	teamRepository *repository.TeamRepository
	userRepository *repository.UserRepository
	transactor     txs.Transactor
}

func NewTeamUseCase(teamRepository *repository.TeamRepository, userRepository *repository.UserRepository, transactor txs.Transactor) *TeamUseCase {
	return &TeamUseCase{teamRepository: teamRepository, userRepository: userRepository, transactor: transactor}
}

func (uc *TeamUseCase) AddTeam(ctx context.Context, team *domain.Team, users []*domain.User) (*domain.Team, []*domain.User, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var addedTeam *domain.Team
	var addedUsers []*domain.User

	err := uc.transactor.WithTransaction(timeOutCtx, func(opCtx context.Context) error {
		existedTeam, err := uc.teamRepository.GetTeam(opCtx, team.TeamName)
		if err != nil {
			return fmt.Errorf("teamRepository: GetTeam: %w", err)
		}

		if existedTeam != nil {
			return &domain.ErrTeamExist{
				TeamName: team.TeamName,
			}
		}

		addedTeam, err = uc.teamRepository.AddTeam(opCtx, team)
		if err != nil {
			return fmt.Errorf("teamRepository: AddTeam: %w", err)
		}

		addedUsers, err = uc.userRepository.AddUsers(opCtx, users)
		if err != nil {
			return fmt.Errorf("teamRepository: AddUsers: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return addedTeam, addedUsers, nil
}

func (uc *TeamUseCase) GetTeam(ctx context.Context, teamName string) (*domain.Team, []*domain.User, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var team *domain.Team
	var users []*domain.User
	var err error

	err = uc.transactor.WithTransaction(timeOutCtx, func(opCtx context.Context) error {
		team, err = uc.teamRepository.GetTeam(opCtx, teamName)
		if err != nil {
			return fmt.Errorf("teamRepository: GetTeam: %w", err)
		}

		if team == nil {
			return &domain.ErrNotFound{
				Source: "team",
				ID:     teamName,
			}
		}

		users, err = uc.userRepository.GetUsersByTeamName(opCtx, teamName)
		if err != nil {
			return fmt.Errorf("teamRepository: GetUsersByTeamName: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return team, users, nil
}
