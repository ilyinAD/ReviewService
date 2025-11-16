package handlers

import (
	"avitostazhko/api"
	"avitostazhko/internal/domain"
	"avitostazhko/internal/infrastructure/usecases/pullrequestusecase"
	"avitostazhko/internal/infrastructure/usecases/teamusecase"
	"avitostazhko/internal/infrastructure/usecases/userusecase"
	"errors"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type ReviewService struct {
	logger             *zap.Logger
	teamUseCase        *teamusecase.TeamUseCase
	userUseCase        *userusecase.UserUseCase
	pullRequestUseCase *pullrequestusecase.PullRequestUseCase
}

func NewReviewService(teamUseCase *teamusecase.TeamUseCase, userUseCase *userusecase.UserUseCase, pullRequestUseCase *pullrequestusecase.PullRequestUseCase) *ReviewService {
	return &ReviewService{
		zap.NewExample(),
		teamUseCase,
		userUseCase,
		pullRequestUseCase,
	}
}

func BuildError(code, message string) api.ErrorResponse {
	return api.ErrorResponse{
		Error: struct {
			Code    api.ErrorResponseErrorCode `json:"code"`
			Message string                     `json:"message"`
		}{
			Code:    api.ErrorResponseErrorCode(code),
			Message: message,
		},
	}
}

func WrapError(ctx echo.Context, err error) error {
	var errNotFound *domain.ErrNotFound
	var errPRExist *domain.ErrPullRequestExists
	var errPRMerged *domain.ErrPRMerged
	var errNotAssigned *domain.ErrNotAssigned
	var errNoCandidate *domain.ErrNoCandidate
	var errTeamExists *domain.ErrTeamExist

	if errors.As(err, &errNotFound) {
		return ctx.JSON(404, BuildError("NOT_FOUND", err.Error()))
	}
	if errors.As(err, &errPRExist) {
		return ctx.JSON(404, BuildError("PR_EXIST", err.Error()))
	}
	if errors.As(err, &errPRMerged) {
		return ctx.JSON(404, BuildError("PR_MERGED", err.Error()))
	}
	if errors.As(err, &errNotAssigned) {
		return ctx.JSON(404, BuildError("NOT_ASSIGNED", err.Error()))
	}
	if errors.As(err, &errNoCandidate) {
		return ctx.JSON(404, BuildError("NO_CANDIDATE", err.Error()))
	}
	if errors.As(err, &errTeamExists) {
		return ctx.JSON(404, BuildError("TEAM_EXISTS", err.Error()))
	}

	return ctx.JSON(500, BuildError("INTERNAL_ERROR", err.Error()))
}
