package handlers

import (
	"avitostazhko/api"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) PostUsersSetIsActive(ctx echo.Context) error {
	s.logger.Info("POST: Start Handling: UsersSetIsActive")
	var req api.PostUsersSetIsActiveJSONRequestBody
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(400, BuildError("INVALID_REQUEST", "Invalid request body"))
	}

	user, err := s.userUseCase.SetUserIsActive(ctx.Request().Context(), req.UserId, req.IsActive)
	if err != nil {
		return WrapError(ctx, err)
	}

	response := api.User{
		UserId:   user.UserID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}

	return ctx.JSON(200, map[string]interface{}{
		"user": response,
	})
}
