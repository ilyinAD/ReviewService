package handlers

import (
	"avitostazhko/api"
	"avitostazhko/internal/domain"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) PostTeamAdd(ctx echo.Context) error {
	s.logger.Info("POST: Start Handling: TeamAdd")
	var team api.Team
	if err := ctx.Bind(&team); err != nil {
		return ctx.JSON(http.StatusBadRequest, api.ErrorResponse{
			Error: struct {
				Code    api.ErrorResponseErrorCode `json:"code"`
				Message string                     `json:"message"`
			}{
				Code:    "BAD_REQUEST",
				Message: err.Error(),
			},
		})
	}

	var users []*domain.User
	for _, user := range team.Members {
		users = append(users, &domain.User{
			UserID:   user.UserId,
			IsActive: user.IsActive,
			Username: user.Username,
			TeamName: team.TeamName,
		})
	}

	addedTeam, addedUsers, err := s.teamUseCase.AddTeam(ctx.Request().Context(), domain.NewTeam(team.TeamName), users)
	if err != nil {
		return WrapError(ctx, err)
	}

	var teamResponse api.Team
	teamResponse.TeamName = addedTeam.TeamName
	for _, user := range addedUsers {
		teamResponse.Members = append(teamResponse.Members, api.TeamMember{
			IsActive: user.IsActive,
			UserId:   user.UserID,
			Username: user.Username,
		})
	}
	return ctx.JSON(http.StatusOK, teamResponse)
}
