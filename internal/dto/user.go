package dto

import (
	"errors"

	"github.com/hel1th/PR_Assigner/internal/domain"
)

type SetUserActiveRequest struct {
	UserID   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}

type SetUserActiveResponse struct {
	User UserResponse `json:"user"`
}

type GetUserReviewPRsResponse struct {
	UserID       string                     `json:"user_id"`
	PullRequests []PullRequestShortResponse `json:"pull_requests"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PullRequestShortResponse struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
	Status          string `json:"status"`
}

func (r *SetUserActiveRequest) Validate() error {
	if r.UserID == "" {
		return errors.New("user_id is required")
	}
	return nil
}

func UserFromDomain(user *domain.User) UserResponse {
	return UserResponse{
		UserID:   user.ID,
		Username: user.Username,
		TeamName: user.TeamName,
		IsActive: user.IsActive,
	}
}

func PullRequestsShortFromDomain(prs []*domain.PullRequest) []PullRequestShortResponse {
	resp := make([]PullRequestShortResponse, len(prs))
	for i, pr := range prs {
		resp[i] = PullRequestShortResponse{
			PullRequestID:   pr.ID,
			PullRequestName: pr.Name,
			AuthorID:        pr.AuthorID,
			Status:          string(pr.Status),
		}
	}
	return resp
}
