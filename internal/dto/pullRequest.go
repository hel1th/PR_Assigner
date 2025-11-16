package dto

import (
	"errors"
	"time"

	"github.com/hel1th/PR_Assigner/internal/domain"
)

type CreatePullRequestRequest struct {
	PullRequestID   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorID        string `json:"author_id"`
}

type CreatePullRequestResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type MergePullRequestRequest struct {
	PullRequestID string `json:"pull_request_id"`
}

type MergePullRequestResponse struct {
	PR PullRequestResponse `json:"pr"`
}

type ReassignReviewerRequest struct {
	PullRequestID string `json:"pull_request_id"`
	OldUserID     string `json:"old_user_id"`
}

type ReassignReviewerResponse struct {
	PR         PullRequestResponse `json:"pr"`
	ReplacedBy string              `json:"replaced_by"`
}

type PullRequestResponse struct {
	PullRequestID     string     `json:"pull_request_id"`
	PullRequestName   string     `json:"pull_request_name"`
	AuthorID          string     `json:"author_id"`
	Status            string     `json:"status"`
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         time.Time  `json:"created_at"`
	MergedAt          *time.Time `json:"merged_at,omitempty"`
}

func (r *CreatePullRequestRequest) Validate() error {
	if r.PullRequestID == "" {
		return errors.New("pull_request_id is required")
	}
	if r.PullRequestName == "" {
		return errors.New("pull_request_name is required")
	}
	if r.AuthorID == "" {
		return errors.New("author_id is required")
	}
	return nil
}

func (r *MergePullRequestRequest) Validate() error {
	if r.PullRequestID == "" {
		return errors.New("pull_request_id is required")
	}
	return nil
}

func (r *ReassignReviewerRequest) Validate() error {
	if r.PullRequestID == "" {
		return errors.New("pull_request_id is required")
	}
	if r.OldUserID == "" {
		return errors.New("old_user_id is required")
	}
	return nil
}

func PullRequestFromDomain(pr *domain.PullRequest) PullRequestResponse {
	return PullRequestResponse{
		PullRequestID:     pr.ID,
		PullRequestName:   pr.Name,
		AuthorID:          pr.AuthorID,
		Status:            string(pr.Status),
		AssignedReviewers: pr.AssignedReviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func PullRequestsFromDomain(prs []*domain.PullRequest) []PullRequestResponse {
	resp := make([]PullRequestResponse, len(prs))
	for i, pr := range prs {
		resp[i] = PullRequestFromDomain(pr)
	}
	return resp
}
