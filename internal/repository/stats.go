package repository

import (
	"context"
	"database/sql"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
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
			COUNT(rev.pr_id) AS total_reviews,
			COUNT(*) FILTER (WHERE pr.status = 'open') AS open_reviews,
			COUNT(*) FILTER (WHERE pr.status = 'merged') AS completed_reviews
		FROM users u
		LEFT JOIN reviewers rev ON u.id = rev.reviewer_id
		LEFT JOIN pull_requests pr ON rev.pr_id = pr.id
	`

	var args []interface{}
	if teamName != nil {
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`
		if err := r.db.QueryRowContext(ctx, checkQuery, teamName).Scan(&exists); err != nil {
			return nil, err
		}
		if !exists {
			return nil, apperrors.NotFound
		}
		query += " WHERE u.team_name = $1 "
		args = append(args, *teamName)
	}

	query += `GROUP BY u.id ORDER BY total_reviews DESC`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*domain.UserStats

	for rows.Next() {
		var s domain.UserStats
		if err := rows.Scan(
			&s.UserID,
			&s.Username,
			&s.TeamName,
			&s.IsActive,
			&s.TotalReviews,
			&s.OpenReviews,
			&s.CompletedReviews,
		); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}

	return out, rows.Err()
}
