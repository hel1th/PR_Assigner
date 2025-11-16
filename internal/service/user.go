package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
	"github.com/hel1th/PR_Assigner/internal/repository"
)

type UserService interface {
	GetUser(ctx context.Context, userID string) (*domain.User, error)
	SetUserActive(ctx context.Context, userID string, isActive bool) error
	GetUserReviewPRs(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUser(ctx context.Context, userID string) (*domain.User, error) {
	if userID == "" {
		return nil, apperrors.InvalidInput
	}

	user, err := s.userRepo.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *userService) SetUserActive(ctx context.Context, userID string, isActive bool) error {
	if userID == "" {
		return apperrors.InvalidInput
	}

	if err := s.userRepo.SetActive(ctx, userID, isActive); err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return apperrors.NotFound
		}
		return fmt.Errorf("failed to set user active status: %w", err)
	}

	return nil
}

func (s *userService) GetUserReviewPRs(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	if userID == "" {
		return nil, apperrors.InvalidInput
	}

	prs, err := s.userRepo.GetReviewPRs(ctx, userID)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, fmt.Errorf("failed to get review PRs: %w", err)
	}

	return prs, nil
}
