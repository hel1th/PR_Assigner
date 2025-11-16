package dto

import (
	"errors"

	"github.com/hel1th/PR_Assigner/internal/domain"
)

type CreateTeamRequest struct {
	TeamName string              `json:"team_name"`
	Members  []TeamMemberRequest `json:"members"`
}

type GetTeamRequest struct {
	TeamName string `json:"team_name"`
}

type TeamMemberRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type CreateTeamResponse struct {
	Team TeamResponse `json:"team"`
}

type TeamResponse struct {
	TeamName string               `json:"team_name"`
	Members  []TeamMemberResponse `json:"members"`
}

type TeamMemberResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func (r *CreateTeamRequest) Validate() error {
	if r.TeamName == "" {
		return errors.New("team_name is required")
	}

	if len(r.Members) == 0 {
		return errors.New("at least one member is required")
	}

	for i, m := range r.Members {
		if m.UserID == "" {
			return errors.New("member user_id is required at index " + string(rune(i)))
		}
		if m.Username == "" {
			return errors.New("member username is required at index " + string(rune(i)))
		}
	}

	return nil
}

func (r *GetTeamRequest) Validate() error {
	if r.TeamName == "" {
		return errors.New("team_name is required")
	}
	return nil
}

func (r *CreateTeamRequest) ToDomainMembers() []domain.TeamMember {
	members := make([]domain.TeamMember, len(r.Members))
	for i, m := range r.Members {
		members[i] = domain.TeamMember{
			ID:       m.UserID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}
	return members
}

func TeamFromDomain(team *domain.Team) TeamResponse {
	members := make([]TeamMemberResponse, len(team.Members))
	for i, m := range team.Members {
		members[i] = TeamMemberResponse{
			UserID:   m.ID,
			Username: m.Username,
			IsActive: m.IsActive,
		}
	}

	return TeamResponse{
		TeamName: team.Name,
		Members:  members,
	}
}
