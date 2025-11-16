package repository

import (
	"avitostazhko/internal/domain"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool}
}

func (r *PullRequestRepository) Create(ctx context.Context, pullRequest *domain.PullRequest) (*domain.PullRequest, error) {
	var addedPR domain.PullRequest
	query := `
		insert into pull_requests (pull_request_id, pull_request_name, author_id, status, created_at) 
		values ($1, $2, $3, $4, $5) returning pull_request_id, pull_request_name, author_id, status, created_at;
	`
	err := r.pool.QueryRow(ctx, query,
		pullRequest.PullRequestID, pullRequest.PullRequestName, pullRequest.AuthorID, pullRequest.Status, time.Now()).
		Scan(&addedPR.PullRequestID, &addedPR.PullRequestName, &addedPR.AuthorID, &addedPR.Status, &addedPR.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert pull request: %w", err)
	}

	return &addedPR, nil
}

func (r *PullRequestRepository) Get(ctx context.Context, pullRequestID string) (*domain.PullRequest, error) {
	query := `select pull_request_id, pull_request_name, author_id, status, created_at, merged_at FROM pull_requests WHERE pull_request_id = $1`
	var pr domain.PullRequest
	err := r.pool.QueryRow(ctx, query, pullRequestID).
		Scan(&pr.PullRequestID, &pr.PullRequestName, &pr.AuthorID, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get pull request: %w", err)
	}

	return &pr, nil
}

func (r *PullRequestRepository) Update(ctx context.Context, pr *domain.PullRequest) (*domain.PullRequest, error) {
	query := `
		update pull_requests 
		set pull_request_name = $1, status = $2, merged_at = $3 
		where pull_request_id = $4 returning pull_request_id, pull_request_name, author_id, status, created_at, merged_at;
	`
	var updatedPR domain.PullRequest
	err := r.pool.QueryRow(ctx, query,
		pr.PullRequestName, pr.Status, pr.MergedAt, pr.PullRequestID).
		Scan(&updatedPR.PullRequestID, &updatedPR.PullRequestName, &updatedPR.AuthorID, &updatedPR.Status, &updatedPR.CreatedAt, &updatedPR.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to update PR: %w", err)
	}

	return &updatedPR, nil
}

func (r *PullRequestRepository) GetPullRequestsByUserID(ctx context.Context, userID string) ([]string, error) {
	query := `
		select pull_request_id from reviewers where reviewer_id = $1
	`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []string
	for rows.Next() {
		var prID string

		err := rows.Scan(&prID)
		if err != nil {
			return nil, err
		}

		prs = append(prs, prID)
	}
	return prs, nil
}
