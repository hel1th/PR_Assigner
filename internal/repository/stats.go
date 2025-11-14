package repository

import (
	"context"
	"database/sql"

	"github.com/hel1th/PR_Assigner/internal/domain"
)

type StatsRepository interface {
	GetUserStats(ctx context.Context, teamName *string) ([]*domain.UserStats, error)
}

type statsRepo struct {
	db *sql.DB
}

func NewStatsRepository(db *sql.DB) StatsRepository {
	return &statsRepo{db: db}
}

func (r *statsRepo) GetUserStats(ctx context.Context, teamName *string) ([]*domain.UserStats, error) {
	query := `
		SELECT 
			u.id,
			u.username,
			u.team_name,
			u.is_active,
			COUNT(rev.pr_id) as total_reviews,
			COUNT(CASE WHEN pr.status = 'open' THEN 1 END) as open_reviews,
			COUNT(CASE WHEN pr.status = 'merged' THEN 1 END) as completed_reviews
		FROM users u
		LEFT JOIN reviewers rev ON u.id = rev.reviewer_id
		LEFT JOIN pull_requests pr ON rev.pr_id = pr.id
	`

	args := make([]interface{}, 0)
	if teamName != nil {
		query += ` WHERE u.team_name = $1`
		args = append(args, *teamName)
	}

	query += `
		GROUP BY u.id, u.username, u.team_name, u.is_active
		ORDER BY total_reviews DESC
	`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make([]*domain.UserStats, 0)
	for rows.Next() {
		var s domain.UserStats
		err := rows.Scan(
			&s.UserID,
			&s.Username,
			&s.TeamName,
			&s.IsActive,
			&s.TotalReviews,
			&s.OpenReviews,
			&s.CompletedReviews,
		)
		if err != nil {
			return nil, err
		}
		stats = append(stats, &s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return stats, nil
}
