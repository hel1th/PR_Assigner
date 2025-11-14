package domain

type UserStats struct {
	UserID           string `json:"user_id"`
	Username         string `json:"username"`
	TeamName         string `json:"team_name"`
	IsActive         bool   `json:"is_active"`
	TotalReviews     int    `json:"total_reviews_assigned"`
	OpenReviews      int    `json:"open_reviews"`
	CompletedReviews int    `json:"completed_reviews"`
}
