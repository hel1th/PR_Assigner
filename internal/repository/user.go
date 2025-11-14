package repository

import (
	"context"
	"database/sql"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
)

type UserRepository interface {
	SetActive(ctx context.Context, id string, isActive bool) error
	GetReviewPRs(ctx context.Context, userID string) ([]*domain.PullRequest, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) SetActive(ctx context.Context, id string, isActive bool) error {
	query := `
		UPDATE users
		SET is_active = $2
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id, isActive)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return apperrors.NotFound
	}

	return nil
}

func (r *userRepo) GetReviewPRs(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	query := `
		SELECT pr.id, pr.name, pr.author_id, pr.status
		FROM pull_requests pr
		INNER JOIN reviewers r ON pr.id = r.pr_id
		WHERE r.reviewer_id = $1
		ORDER BY pr.created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prs := make([]*domain.PullRequest, 0)
	for rows.Next() {
		var pr domain.PullRequest
		if err := rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, &pr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prs, nil
}
