package domain

import (
	"time"
)

type PRStatus string

const (
	PROpen       PRStatus = "Open"
	PRMerged     PRStatus = "Merged"
	MaxReviewers int      = 2
)

func (s PRStatus) IsValid() bool {
	return s == PROpen || s == PRMerged
}

func (s PRStatus) String() string {
	return string(s)
}

type PullRequest struct {
	ID                string
	Name              string
	AuthorID          string
	Status            PRStatus
	AssignedReviewers []string
	CreatedAt         time.Time
	MergedAt          *time.Time
}

func (pr *PullRequest) IsMerged() bool {
	return pr.Status == PRMerged
}

func (pr *PullRequest) IsOpen() bool {
	return pr.Status == PROpen
}

func (pr *PullRequest) ReviewersCount() int {
	return len(pr.AssignedReviewers)
}

func (pr *PullRequest) NeedMoreReviewers() bool {
	if pr.IsMerged() {
		return false
	}

	return pr.ReviewersCount() < MaxReviewers
}
