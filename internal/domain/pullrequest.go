package domain

import "time"

type Status string

const (
	StatusOpen   Status = "OPEN"
	StatusMerged Status = "MERGED"
)

type PullRequest struct {
	PullRequestID   string
	PullRequestName string
	AuthorID        string
	Status          Status
	CreatedAt       *time.Time
	MergedAt        *time.Time
}

type PullRequestWithReviewers struct {
	PullRequest *PullRequest
	Reviewers   []*Reviewer
}
