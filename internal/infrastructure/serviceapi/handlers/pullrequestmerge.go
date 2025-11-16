package handlers

import (
	"avitostazhko/api"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) PostPullRequestMerge(ctx echo.Context) error {
	s.logger.Info("POST: Start Handling: PullRequestMerge")

	var req api.PostPullRequestMergeJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, BuildError("INVALID_REQUEST", "Invalid request body"))
	}

	pullRequest, reviewers, err := s.pullRequestUseCase.MergePR(ctx.Request().Context(), req.PullRequestId)
	if err != nil {
		return WrapError(ctx, err)
	}

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

	return ctx.JSON(201, map[string]interface{}{
		"pr": response,
	})
}
