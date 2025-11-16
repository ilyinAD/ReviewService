package repository

import (
	"avitostazhko/internal/domain"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool}
}

func (r *UserRepository) AddUsers(ctx context.Context, users []*domain.User) ([]*domain.User, error) {
	var addedUsers []*domain.User
	for _, user := range users {
		addedUser, err := r.AddUser(ctx, user)
		if err != nil {
			return nil, fmt.Errorf("error adding user: %v", err)
		}

		addedUsers = append(addedUsers, addedUser)
	}

	return addedUsers, nil
}

func (r *UserRepository) AddUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	var addedUser domain.User
	err := r.pool.QueryRow(
		ctx,
		"insert into users (id, username, is_active, team_name) values ($1, $2, $3, $4) returning id, username, is_active, team_name",
		user.UserID, user.Username, user.IsActive, user.TeamName).
		Scan(&addedUser.UserID, &addedUser.Username, &addedUser.IsActive, &addedUser.TeamName)
	if err != nil {
		return nil, fmt.Errorf("error: failed to insert user: %w", err)
	}

	return &addedUser, nil
}

func (r *UserRepository) GetUsersByTeamName(ctx context.Context, teamName string) ([]*domain.User, error) {
	query := `select id, username, is_active, team_name from users where team_name = $1`
	rows, err := r.pool.Query(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(&user.UserID, &user.Username, &user.IsActive, &user.TeamName)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	query := `update users set username = $1, is_active = $2, team_name = $3 where id = $4`

	_, err := r.pool.Exec(ctx, query, user.Username, user.IsActive, user.TeamName, user.UserID)
	if err != nil {
		return nil, fmt.Errorf("error: failed to update user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	query := `select id, username, is_active, team_name from users where id = $1`
	var user domain.User
	err := r.pool.QueryRow(ctx, query, userID).Scan(&user.UserID, &user.Username, &user.IsActive, &user.TeamName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("error: failed to query user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetActiveUsers(ctx context.Context, teamName string, authorID string) ([]*domain.User, error) {
	query := `select id, username, is_active, team_name from users where team_name = $1 and is_active = true and id != $2`

	rows, err := r.pool.Query(ctx, query, teamName, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User

	for rows.Next() {
		var user domain.User
		err = rows.Scan(&user.UserID, &user.Username, &user.IsActive, &user.TeamName)
		if err != nil {
			return nil, fmt.Errorf("error: failed to query user: %w", err)
		}

		users = append(users, &user)
	}

	return users, nil
}
