package dto

import "github.com/hel1th/PR_Assigner/internal/domain"

type GetUserStatsResponse struct {
	Stats []UserStats `json:"stats"`
}

type UserStats struct {
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	TeamName         string `json:"team_name"`
	IsActive         bool   `json:"is_active"`
	TotalReviews     int    `json:"total_reviews_assigned"`
	OpenReviews      int    `json:"open_reviews"`
	CompletedReviews int    `json:"completed_reviews"`
}

func UserStatsFromDomain(stats []*domain.UserStats) []UserStats {
	dtos := make([]UserStats, len(stats))
	for i, s := range stats {
		dtos[i] = UserStats{
			UserID:           s.UserID,
			Username:         s.Username,
			TeamName:         s.TeamName,
			IsActive:         s.IsActive,
			TotalReviews:     s.TotalReviews,
			OpenReviews:      s.OpenReviews,
			CompletedReviews: s.CompletedReviews,
		}
	}
	return dtos
}
