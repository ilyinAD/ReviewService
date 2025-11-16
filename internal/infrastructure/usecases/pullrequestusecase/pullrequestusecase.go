package pullrequestusecase

import (
	"avitostazhko/internal/domain"
	"avitostazhko/internal/infrastructure/repository"
	"avitostazhko/internal/infrastructure/repository/txs"
	"context"
	"fmt"
	"math/rand"
	"time"
)

type PullRequestUseCase struct {
	pullRequestRepository *repository.PullRequestRepository
	userRepository        *repository.UserRepository
	reviewersRepository   *repository.ReviewersRepository
	transactor            txs.Transactor
}

func NewPullRequestUseCase(
	pullRequestRepository *repository.PullRequestRepository,
	userRepository *repository.UserRepository,
	reviewersRepository *repository.ReviewersRepository,
	transactor txs.Transactor,
) *PullRequestUseCase {
	return &PullRequestUseCase{
		pullRequestRepository: pullRequestRepository,
		userRepository:        userRepository,
		reviewersRepository:   reviewersRepository,
		transactor:            transactor,
	}
}

func (uc *PullRequestUseCase) CreatePR(ctx context.Context, pullRequestID, pullRequestName, authorID string) (*domain.PullRequest, []*domain.Reviewer, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var createdPR *domain.PullRequest
	var insertedReviewers []*domain.Reviewer

	err := uc.transactor.WithTransaction(timeOutCtx, func(opCtx context.Context) error {
		pr, err := uc.pullRequestRepository.Get(opCtx, pullRequestID)
		if err != nil {
			return err
		}
		if pr != nil {
			return &domain.ErrPullRequestExists{
				PullRequestID: pullRequestID,
			}
		}

		author, err := uc.userRepository.GetUserByID(opCtx, authorID)
		if err != nil {
			return err
		}
		if author == nil {
			return &domain.ErrNotFound{
				Source: "user",
				ID:     authorID,
			}
		}

		activeUsers, err := uc.userRepository.GetActiveUsers(opCtx, author.TeamName, authorID)
		if err != nil {
			return fmt.Errorf("userRepository: GetActiveUsers: %w", err)
		}

		var reviewers []*domain.User
		if len(activeUsers) > 0 {
			shuffled := make([]*domain.User, len(activeUsers))
			copy(shuffled, activeUsers)
			rand.Shuffle(len(shuffled), func(i, j int) {
				shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
			})
			reviewers = shuffled[:min(2, len(activeUsers))]
		}

		pr = &domain.PullRequest{
			PullRequestID:   pullRequestID,
			PullRequestName: pullRequestName,
			AuthorID:        authorID,
			Status:          domain.StatusOpen,
		}

		createdPR, err = uc.pullRequestRepository.Create(opCtx, pr)
		if err != nil {
			return fmt.Errorf("pullRequestRepository: Create: %w", err)
		}

		for _, reviewer := range reviewers {
			insertedReviewer, err := uc.reviewersRepository.InsertReviewer(ctx, &domain.Reviewer{
				ReviewerID:    reviewer.UserID,
				PullRequestID: pullRequestID,
			})
			if err != nil {
				return fmt.Errorf("reviewersRepository: InsertReviewer: %w", err)
			}

			insertedReviewers = append(insertedReviewers, insertedReviewer)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return createdPR, insertedReviewers, nil
}

func (uc *PullRequestUseCase) MergePR(ctx context.Context, pullRequestID string) (*domain.PullRequest, []*domain.Reviewer, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var updatedPR *domain.PullRequest
	var reviewers []*domain.Reviewer

	err := uc.transactor.WithTransaction(timeOutCtx, func(opCtx context.Context) error {
		pr, err := uc.pullRequestRepository.Get(ctx, pullRequestID)
		if err != nil {
			return fmt.Errorf("pullRequestRepository: Get: %w", err)
		}
		if pr == nil {
			return &domain.ErrNotFound{
				Source: "pull request",
				ID:     pullRequestID,
			}
		}

		updatedPR = pr

		if pr.Status != domain.StatusMerged {
			now := time.Now()
			pr.Status = domain.StatusMerged
			pr.MergedAt = &now
			updatedPR, err = uc.pullRequestRepository.Update(opCtx, pr)
			if err != nil {
				return fmt.Errorf("pullRequestRepository: Update: %w", err)
			}
		}

		reviewers, err = uc.reviewersRepository.GetReviewersByPR(opCtx, pullRequestID)
		if err != nil {
			return fmt.Errorf("reviewersRepository: GetReviewers: %w", err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return updatedPR, reviewers, nil
}

func (uc *PullRequestUseCase) ReassignReviewer(ctx context.Context, pullRequestID, oldUserID string) (*domain.PullRequestWithReviewers, string, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var result *domain.PullRequestWithReviewers
	var newReviewer *domain.User

	err := uc.transactor.WithTransaction(timeOutCtx, func(opCtx context.Context) error {
		pr, err := uc.pullRequestRepository.Get(opCtx, pullRequestID)
		if err != nil {
			return fmt.Errorf("pullRequestRepository: Get: %w", err)
		}
		if pr == nil {
			return &domain.ErrNotFound{
				Source: "pull request",
				ID:     pullRequestID,
			}
		}

		if pr.Status == domain.StatusMerged {
			return &domain.ErrPRMerged{
				PullRequestID: pullRequestID,
			}
		}

		currentReviewers, err := uc.reviewersRepository.GetReviewersByPR(opCtx, pullRequestID)
		if err != nil {
			return fmt.Errorf("reviewersRepository: GetReviewersByPR: %w", err)
		}

		ok := false
		for _, reviewer := range currentReviewers {
			if reviewer.ReviewerID == oldUserID {
				ok = true
				break
			}
		}
		if !ok {
			return &domain.ErrNotAssigned{
				OldUserID: oldUserID,
			}
		}

		oldReviewer, err := uc.userRepository.GetUserByID(opCtx, oldUserID)
		if err != nil {
			return fmt.Errorf("userRepository: GetUserByID: %w", err)
		}
		if oldReviewer == nil {
			return &domain.ErrNotFound{
				Source: "user",
				ID:     oldUserID,
			}
		}

		activeUsers, err := uc.userRepository.GetActiveUsers(opCtx, oldReviewer.TeamName, pr.AuthorID)
		if err != nil {
			return fmt.Errorf("userRepository: GetActiveUsers: %w", err)
		}

		var available []int

		for idx, activeUser := range activeUsers {
			isAvailable := true
			for _, reviewer := range currentReviewers {
				if reviewer.ReviewerID == activeUser.UserID {
					isAvailable = false
					break
				}
			}

			if isAvailable {
				available = append(available, idx)
			}
		}

		if len(available) == 0 {
			return &domain.ErrNoCandidate{}
		}

		newReviewer = activeUsers[available[rand.Intn(len(available))]]

		if err := uc.reviewersRepository.RemoveReviewer(opCtx, pullRequestID, oldUserID); err != nil {
			return fmt.Errorf("reviewersRepository: RemoveReviewer: %w", err)
		}

		_, err = uc.reviewersRepository.InsertReviewer(opCtx, &domain.Reviewer{
			ReviewerID:    newReviewer.UserID,
			PullRequestID: pullRequestID,
		})
		if err != nil {
			return fmt.Errorf("reviewersRepository: InsertReviewer: %w", err)
		}

		updatedReviewers, err := uc.reviewersRepository.GetReviewersByPR(opCtx, pullRequestID)
		if err != nil {
			return fmt.Errorf("reviewersRepository: GetReviewersByPR: %w", err)
		}

		result = &domain.PullRequestWithReviewers{
			PullRequest: pr,
			Reviewers:   updatedReviewers,
		}

		return nil
	})
	if err != nil {
		return nil, "", err
	}

	return result, newReviewer.UserID, nil
}

func (uc *PullRequestUseCase) GetPullRequestsByUserID(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	timeOutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var result []*domain.PullRequest

	err := uc.transactor.WithTransaction(timeOutCtx, func(opCtx context.Context) error {
		prIDs, err := uc.pullRequestRepository.GetPullRequestsByUserID(opCtx, userID)
		if err != nil {
			return fmt.Errorf("pullRequestRepository: GetPullRequestsByUserID: %w", err)
		}

		result = make([]*domain.PullRequest, len(prIDs))
		for idx, prID := range prIDs {
			pr, err := uc.pullRequestRepository.Get(ctx, prID)
			if err != nil {
				return fmt.Errorf("pullRequestRepository: Get: %w", err)
			}

			result[idx] = pr
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}
