package repository

import (
	"avitostazhko/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReviewersRepository struct {
	pool *pgxpool.Pool
}

func NewReviewersRepository(pool *pgxpool.Pool) *ReviewersRepository {
	return &ReviewersRepository{pool}
}

func (r *ReviewersRepository) InsertReviewer(ctx context.Context, reviewer *domain.Reviewer) (*domain.Reviewer, error) {
	var addedReviewer domain.Reviewer
	query := `insert into reviewers (reviewer_id, pull_request_id) VALUES ($1, $2) returning reviewer_id, pull_request_id`
	err := r.pool.QueryRow(ctx, query, reviewer.ReviewerID, reviewer.PullRequestID).Scan(&addedReviewer.ReviewerID, &addedReviewer.PullRequestID)
	if err != nil {
		return nil, err
	}

	return &addedReviewer, nil
}

func (r *ReviewersRepository) GetReviewersByPR(ctx context.Context, pullRequestID string) ([]*domain.Reviewer, error) {
	query := `select reviewer_id, pull_request_id from reviewers where pull_request_id = $1`
	rows, err := r.pool.Query(ctx, query, pullRequestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get the reviewers by pull request id: %w", err)
	}

	defer rows.Close()

	reviewers := make([]*domain.Reviewer, 0)

	for rows.Next() {
		var reviewer domain.Reviewer
		err = rows.Scan(&reviewer.ReviewerID, &reviewer.PullRequestID)
		if err != nil {
			return nil, fmt.Errorf("failed to get the reviewer by pull request id: %w", err)
		}

		reviewers = append(reviewers, &reviewer)
	}

	return reviewers, nil
}

func (r *ReviewersRepository) RemoveReviewer(ctx context.Context, pullRequestID, reviewerID string) error {
	query := `DELETE FROM reviewers WHERE pull_request_id = $1 AND reviewer_id = $2`
	_, err := r.pool.Exec(ctx, query, pullRequestID, reviewerID)
	if err != nil {
		return fmt.Errorf("failed to remove reviewer by pull request id: %w", err)
	}

	return nil
}
