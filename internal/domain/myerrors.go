package domain

import (
	"fmt"
)

type ErrPullRequestExists struct {
	PullRequestID string
}

func (e *ErrPullRequestExists) Error() string {
	return fmt.Sprintf("pull request with id %s already exists", e.PullRequestID)
}

type ErrNotFound struct {
	Source string
	ID     string
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("%s with id: %s not found", e.Source, e.ID)
}

type ErrPRMerged struct {
	PullRequestID string
}

func (e *ErrPRMerged) Error() string {
	return fmt.Sprintf("pull request with id %s merged", e.PullRequestID)
}

type ErrNotAssigned struct {
	OldUserID string
}

func (e *ErrNotAssigned) Error() string {
	return fmt.Sprintf("PR not assigned to user %s", e.OldUserID)
}

type ErrNoCandidate struct{}

func (e *ErrNoCandidate) Error() string {
	return fmt.Sprintf("No candidate found")
}

type ErrTeamExist struct {
	TeamName string
}

func (e *ErrTeamExist) Error() string {
	return fmt.Sprintf("Team - %s already exists", e.TeamName)
}
