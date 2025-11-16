package repository

import (
	"avitostazhko/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool}
}

func (r *TeamRepository) AddTeam(ctx context.Context, team *domain.Team) (*domain.Team, error) {
	var addedTeam domain.Team
	err := r.pool.QueryRow(
		ctx,
		"insert into teams (team_name) values ($1) returning team_name", team.TeamName).
		Scan(&addedTeam.TeamName)
	if err != nil {

		return nil, fmt.Errorf("error: failed to insert team: %w", err)
	}

	return &addedTeam, nil
}

func (r *TeamRepository) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	var team domain.Team

	err := r.pool.QueryRow(ctx, "select team_name from teams where team_name = $1", teamName).Scan(&team.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("chatRepository.GetChatByID: %w", err)
	}

	return &team, nil
}
