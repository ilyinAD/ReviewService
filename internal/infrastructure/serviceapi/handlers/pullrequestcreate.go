package handlers

import (
	"avitostazhko/api"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) PostPullRequestCreate(ctx echo.Context) error {
	s.logger.Info("POST: Start Handling: PullRequestCreate")

	var req api.PostPullRequestCreateJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, BuildError("INVALID_REQUEST", "Invalid request body"))
	}

	pullRequest, reviewers, err := s.pullRequestUseCase.CreatePR(ctx.Request().Context(), req.PullRequestId, req.PullRequestName, req.AuthorId)
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
