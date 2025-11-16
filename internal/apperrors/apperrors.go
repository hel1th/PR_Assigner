package apperrors

import (
	"fmt"
)

type ErrResponse string

const (
	TeamExists              ErrResponse = "TEAM_EXISTS"
	PRExists                ErrResponse = "PR_EXISTS"
	PRMerged                ErrResponse = "PR_MERGED"
	NotAssigned             ErrResponse = "NOT_ASSIGNED"
	NoCandidate             ErrResponse = "NO_CANDIDATE"
	NotFound                ErrResponse = "NOT_FOUND"
	ReviewerAlreadyAssigned ErrResponse = "ALREADY_ASSIGNED"
	MaxReviewers            ErrResponse = "TO_MANY_REVIEWERS"
	InvalidInput            ErrResponse = "INVALID_INPUT"
)

func (err ErrResponse) Error() string {
	return fmt.Sprintf("Error: %s", string(err))
}
