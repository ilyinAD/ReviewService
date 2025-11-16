package handlers

import (
	"avitostazhko/api"

	"github.com/labstack/echo/v4"
)

func (s *ReviewService) GetTeamGet(ctx echo.Context, params api.GetTeamGetParams) error {
	s.logger.Info("GET: Start Handling: TeamGet")

	team, users, err := s.teamUseCase.GetTeam(ctx.Request().Context(), params.TeamName)
	if err != nil {
		return WrapError(ctx, err)
	}

	var apiUsers []api.TeamMember
	for _, user := range users {
		apiUsers = append(apiUsers, api.TeamMember{
			UserId:   user.UserID,
			Username: user.Username,
			IsActive: user.IsActive,
		})
	}

	response := api.Team{
		TeamName: team.TeamName,
		Members:  apiUsers,
	}

	return ctx.JSON(200, response)
}
