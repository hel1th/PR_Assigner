package service

import (
	"context"
	"fmt"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
	"github.com/hel1th/PR_Assigner/internal/repository"
)

type StatsService interface {
	GetUserStats(ctx context.Context, teamName *string) ([]*domain.UserStats, error)
}

type statsService struct {
	statsRepo repository.StatsRepository
}

func NewStatsService(statsRepo repository.StatsRepository) StatsService {
	return &statsService{
		statsRepo: statsRepo,
	}
}

func (s *statsService) GetUserStats(ctx context.Context, teamName *string) ([]*domain.UserStats, error) {
	if teamName != nil && *teamName == "" {
		return nil, apperrors.InvalidInput
	}

	stats, err := s.statsRepo.GetUserStats(ctx, teamName)
	if err != nil {
		return nil, fmt.Errorf("failed to get user stats: %w", err)
	}

	return stats, nil
}