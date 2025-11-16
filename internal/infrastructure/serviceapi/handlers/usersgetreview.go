package handlers

import (
	"avitostazhko/api"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) GetUsersGetReview(ctx echo.Context, params api.GetUsersGetReviewParams) error {
	s.logger.Info("GET: Start Handling: UsersGetReview")

	prs, err := s.pullRequestUseCase.GetPullRequestsByUserID(ctx.Request().Context(), params.UserId)
	if err != nil {
		return WrapError(ctx, err)
	}

	response := make([]api.PullRequestShort, len(prs))
	for i, pr := range prs {
		response[i] = api.PullRequestShort{
			PullRequestId:   pr.PullRequestID,
			PullRequestName: pr.PullRequestName,
			AuthorId:        pr.AuthorID,
			Status:          api.PullRequestShortStatus(pr.Status),
		}
	}

	return ctx.JSON(200, map[string]interface{}{
		"user_id":       params.UserId,
		"pull_requests": response,
	})
}
