package handlers

import (
	"avitostazhko/api"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) PostPullRequestReassign(ctx echo.Context) error {
	s.logger.Info("POST: Start Handling: PullRequestReassign")

	var req api.PostPullRequestReassignJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, BuildError("INVALID_REQUEST", "Invalid request body"))
	}

	pullRequestWithReviewers, newReviewer, err := s.pullRequestUseCase.ReassignReviewer(ctx.Request().Context(), req.PullRequestId, req.OldUserId)
	if err != nil {
		return WrapError(ctx, err)
	}

	pullRequest := pullRequestWithReviewers.PullRequest
	reviewers := pullRequestWithReviewers.Reviewers

	response := api.PullRequest{
		AuthorId:        pullRequest.AuthorID,
		CreatedAt:       pullRequest.CreatedAt,
		MergedAt:        pullRequest.MergedAt,
		PullRequestId:   pullRequest.PullRequestID,
		PullRequestName: pullRequest.PullRequestName,
		Status:          api.PullRequestStatus(pullRequest.Status),
	}
	response.AssignedReviewers = []string{}
	for _, reviewer := range reviewers {
		response.AssignedReviewers = append(response.AssignedReviewers, reviewer.ReviewerID)
	}

	return ctx.JSON(200, map[string]interface{}{
		"pr":          response,
		"replaced_by": newReviewer,
	})
}
