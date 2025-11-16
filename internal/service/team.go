package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
	"github.com/hel1th/PR_Assigner/internal/repository"
)

type TeamService interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.TeamMember) error
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
}

func NewTeamService(teamRepo repository.TeamRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, teamName string, members []domain.TeamMember) error {
	if teamName == "" {
		return apperrors.InvalidInput
	}

	for _, member := range members {
		if member.ID == "" || member.Username == "" {
			return apperrors.InvalidInput
		}
	}

	exists, err := s.teamRepo.Exists(ctx, teamName)
	if err != nil {
		return fmt.Errorf("failed to check team exisitance: %w", err)
	}
	if exists {
		return apperrors.TeamExists
	}
	if err := s.teamRepo.CreateTeam(ctx, teamName, members); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	return nil
}

func (s *teamService) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	if teamName == "" {
		return nil, apperrors.InvalidInput
	}

	team, err := s.teamRepo.GetTeam(ctx, teamName)
	if err != nil {
		if errors.Is(err, apperrors.NotFound) {
			return nil, apperrors.NotFound
		}
		return nil, fmt.Errorf("failed to get team: %w", err)
	}

	return team, nil
}
