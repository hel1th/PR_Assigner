package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
	"github.com/hel1th/PR_Assigner/internal/repository"
)

type PullRequestService interface {
	CreatePullReq(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error)
	GetPullReq(ctx context.Context, prID string) (*domain.PullRequest, error)
	MergePullReq(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) error
	ListPullReqsByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error)
}

type pullReqSvc struct {
	pullReqRepo repository.PullReqRepository
	userRepo    repository.UserRepository
	teamRepo    repository.TeamRepository
	nowFn       func() time.Time
}

func (s *pullReqSvc) CreatePullReq(ctx context.Context, prID, prName, authorID string) (*domain.PullRequest, error) {
	if prID == "" || prName == "" || authorID == "" {
		return nil, apperrors.InvalidInput
	}

	author, err := s.userRepo.GetUser(ctx, authorID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, err
	}

	exists, err := s.pullReqRepo.Exists(ctx, prID)
	if err != nil {
		if errors.Is(err, apperrors.PRExists) {
			return nil, apperrors.PRExists
		}
		return nil, err
	}
	if exists {
		return nil, apperrors.PRExists
	}

	reviewers, err := s.selectReviewers(ctx, author)
	if err != nil {
		if errors.Is(err, apperrors.NoCandidate) {
			return nil, apperrors.NoCandidate
		}
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, fmt.Errorf("failed to get select reviewers: %w", err)
	}

	pr := &domain.PullRequest{
		ID:                prID,
		Name:              prName,
		AuthorID:          authorID,
		Status:            domain.PROpen,
		AssignedReviewers: reviewers,
		CreatedAt:         s.nowFn(),
		MergedAt:          nil,
	}

	if err := s.pullReqRepo.Create(ctx, pr); err != nil {
		if errors.Is(err, apperrors.PRExists) {
			return nil, apperrors.PRExists
		}

		return nil, fmt.Errorf("failed to create PR: %w", err)
	}

	return pr, nil
}

func (s *pullReqSvc) selectReviewers(ctx context.Context, author *domain.User) (revsID []string, err error) {
	if author.TeamName == "" {
		return []string{}, nil
	}

	candidates, err := s.teamRepo.GetActiveTeamMembers(ctx, author.ID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		if errors.Is(err, apperrors.NoCandidate) {
			return nil, apperrors.NoCandidate
		}
		return nil, fmt.Errorf("failed to get active members: %w", err)
	}

	if candidatesCount := len(candidates); candidatesCount <= domain.MaxReviewers {
		out := make([]string, candidatesCount)
		copy(out, candidates)
		return out, nil
	}

	return s.randomSelect(candidates, domain.MaxReviewers), nil
}

func (s *pullReqSvc) randomSelect(items []string, count int) []string {
	shuffled := make([]string, len(items))
	copy(shuffled, items)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.IntN(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	return shuffled[:count]
}

func (s *pullReqSvc) GetPullReq(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if prID == "" {
		return nil, apperrors.InvalidInput
	}

	pr, err := s.pullReqRepo.GetPullReq(ctx, prID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, fmt.Errorf("failed to get PR: %w", err)
	}
	return pr, nil
}

func (s *pullReqSvc) MergePullReq(ctx context.Context, prID string) (*domain.PullRequest, error) {
	if prID == "" {
		return nil, apperrors.InvalidInput
	}

	_, err := s.pullReqRepo.PullReqMerge(ctx, prID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		if errors.Is(err, apperrors.PRMerged) {
			return s.pullReqRepo.GetPullReq(ctx, prID)
		}
		return nil, fmt.Errorf("failed to merge PR: %w", err)
	}
	return s.pullReqRepo.GetPullReq(ctx, prID)
}

func (s *pullReqSvc) ReassignReviewer(ctx context.Context, prID, oldReviewerID string) error {
	if prID == "" || oldReviewerID == "" {
		return apperrors.InvalidInput
	}

	pr, err := s.pullReqRepo.GetPullReq(ctx, prID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return apperrors.NotFound
		}
		return fmt.Errorf("failed to get PR: %w", err)
	}

	if pr.IsMerged() {
		return apperrors.PRMerged
	}

	oldRevFound := false
	for _, revID := range pr.AssignedReviewers {
		if revID == oldReviewerID {
			oldRevFound = true
		}
	}

	if !oldRevFound {
		return apperrors.NotAssigned
	}

	oldRev, err := s.userRepo.GetUser(ctx, oldReviewerID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return apperrors.NotFound
		}
		return fmt.Errorf("failed to get old reviewer: %w", err)
	}

	newRevsID, err := s.selectReviewers(ctx, oldRev)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return apperrors.NotFound
		}
		if errors.Is(err, apperrors.NoCandidate) {
			return apperrors.NoCandidate
		}
		return err
	}
	newRev, err := s.userRepo.GetUser(ctx, newRevsID[0])
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return apperrors.NotFound
		}
		return fmt.Errorf("failed to get new reviewer: %w", err)
	}

	if newRev.ID == pr.AuthorID {
		return apperrors.InvalidInput
	}

	if !newRev.IsActive {
		return apperrors.NoCandidate
	}

	for _, existingReviewerID := range pr.AssignedReviewers {
		if existingReviewerID == newRev.ID {
			return apperrors.ReviewerAlreadyAssigned
		}
	}

	if oldRev.TeamName == "" || newRev.TeamName == "" {
		return apperrors.InvalidInput
	}
	if oldRev.TeamName != newRev.TeamName {
		return apperrors.NoCandidate
	}

	if err := s.pullReqRepo.ReassignReviewer(ctx, prID, oldReviewerID, newRev.ID); err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return apperrors.NotFound
		}
		return fmt.Errorf("failed to reassign reviewer: %w", err)
	}

	return nil
}

func (s *pullReqSvc) ListPullReqsByReviewer(ctx context.Context, reviewerID string) ([]*domain.PullRequest, error) {
	if reviewerID == "" {
		return nil, apperrors.InvalidInput
	}

	_, err := s.userRepo.GetUser(ctx, reviewerID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, fmt.Errorf("failed to get reviewer: %w", err)
	}

	prs, err := s.pullReqRepo.ListByUser(ctx, reviewerID)
	if err != nil {
		return nil, fmt.Errorf("failed to list PRs: %w", err)
	}

	return prs, nil
}

func NewPullReqSvc(
	prRepo repository.PullReqRepository,
	userRepo repository.UserRepository,
	teamRepo repository.TeamRepository,
) PullRequestService {
	return &pullReqSvc{
		pullReqRepo: prRepo,
		userRepo:    userRepo,
		teamRepo:    teamRepo,
		nowFn:       time.Now,
	}
}
