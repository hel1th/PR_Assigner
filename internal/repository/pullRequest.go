package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
)

type PullReqRepository interface {
	Create(ctx context.Context, pr *domain.PullRequest) error
	GetByID(ctx context.Context, id string) (*domain.PullRequest, error)
	Update(ctx context.Context, pr *domain.PullRequest) error
	Exists(ctx context.Context, id string) (bool, error)

	AssignReviewer(ctx context.Context, prID, reviewerID string) error
	ListReviewers(ctx context.Context, prID string) ([]string, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldRevID, newRevID string) error
	PullReqMerge(ctx context.Context, id string) (*domain.PullRequest, error)
}

type pullReqRepo struct {
	db *sql.DB
}

func (r *pullReqRepo) PullReqMerge(ctx context.Context, id string) (*domain.PullRequest, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var pr domain.PullRequest
	getQuery := `
		SELECT id, name, author_id, 
		       status, created_at, merged_at
		FROM pull_requests 
		WHERE id = $1
		FOR UPDATE
	`

	err = tx.QueryRowContext(ctx, getQuery, id).Scan(
		&pr.ID,
		&pr.Name,
		&pr.AuthorID,
		&pr.Status,
		&pr.CreatedAt,
		&pr.MergedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NotFound
		}
		return nil, err
	}

	if pr.Status == domain.PRMerged {
		return nil, apperrors.PRMerged
	}

	updateQuery := `UPDATE pull_requests SET status = 'merged', merged_at = NOW()
					WHERE id = $1 RETURNING merged_at`

	err = tx.QueryRowContext(ctx, updateQuery, id).Scan(&pr.MergedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to merge pull request: %w", err)
	}
	pr.Status = domain.PRMerged

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *pullReqRepo) Exists(ctx context.Context, prID string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM pull_requests WHERE id=$1)`

	var exists bool
	err := r.db.QueryRowContext(ctx, query, prID).Scan(&exists)

	return exists, err
}

func (r *pullReqRepo) Create(ctx context.Context, pr *domain.PullRequest) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	prQuery := `INSERT INTO pull_requests (id, name, author_id, status, created_at, merged_at)
				VALUES ($1,$2,$3,$4,$5,$6)`

	revQuery := `INSERT INTO reviewers (pr_id, reviewer_id) VALUES ($1,$2)`

	_, err = tx.ExecContext(ctx,
		prQuery, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt)
	if err != nil {
		return err
	}

	for _, reviewerID := range pr.AssignedReviewers {
		_, err = tx.ExecContext(ctx, revQuery, pr.ID, reviewerID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *pullReqRepo) GetByID(ctx context.Context, prID string) (*domain.PullRequest, error) {
	var (
		pr       domain.PullRequest
		status   domain.PRStatus
		mergedAt sql.NullTime
	)

	prQuery := `SELECT id, name, author_id, status, created_at, merged_at FROM pull_requests
				WHERE id=$1`

	revQuery := `SELECT reviewer_id FROM reviewers WHERE pr_id=$1`

	err := r.db.QueryRowContext(ctx, prQuery, prID).Scan(
		&pr.ID, &pr.Name, &pr.AuthorID,
		&status, &pr.CreatedAt, &mergedAt)

	if err != nil {
		return nil, err
	}

	pr.Status = status
	if mergedAt.Valid {
		pr.MergedAt = &mergedAt.Time
	}

	rows, err := r.db.QueryContext(ctx, revQuery, prID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.NotFound
		}
		return nil, err
	}
	defer rows.Close()

	var revID string
	pr.AssignedReviewers = make([]string, 0, 2)
	for rows.Next() {
		if err := rows.Scan(&revID); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, revID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (r *pullReqRepo) Update(ctx context.Context, pr *domain.PullRequest) error {
	query := `UPDATE pull_requests SET name=$2, author_id=$3, status=$4, merged_at=$5 WHERE id=$1`

	result, err := r.db.ExecContext(ctx, query, pr.ID, pr.Name, pr.AuthorID, pr.Status, pr.MergedAt)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperrors.NotFound
	}
	return nil
}

func (r *pullReqRepo) AssignReviewer(ctx context.Context, prID, reviewerID string) error {
	query := `INSERT INTO reviewers (pr_id, reviewer_id) VALUES ($1,$2)`

	_, err := r.db.ExecContext(ctx, query, prID, reviewerID)

	return err
}

func (r *pullReqRepo) ReassignReviewer(ctx context.Context, prID, oldRevID, newRevID string) error {
	query := `UPDATE reviewers SET reviewer_id=$1
				WHERE pr_id=$2 and reviewer_id=$3`

	result, err := r.db.ExecContext(ctx, query, newRevID, prID, oldRevID)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return apperrors.NotFound
	}
	return nil
}

func (r *pullReqRepo) ListByUser(ctx context.Context, userID string) ([]*domain.PullRequest, error) {
	query := `SELECT p.id, p.name, p.author_id, p.status, p.created_at, p.merged_at
				FROM pull_requests p
				JOIN reviewers r ON r.pr_id = p.id
				WHERE r.reviewer_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	prsByUser := make([]*domain.PullRequest, 0)
	for rows.Next() {
		var (
			pr       domain.PullRequest
			mergedAt sql.NullTime
			status   domain.PRStatus
		)
		err = rows.Scan(&pr.ID, &pr.Name, &pr.AuthorID, &status, &pr.CreatedAt, &mergedAt)
		if err != nil {
			return nil, err
		}

		if mergedAt.Valid {
			pr.MergedAt = &mergedAt.Time
		}

		pr.Status = status

		prsByUser = append(prsByUser, &pr)

	}

	return prsByUser, rows.Err()
}

func (r *pullReqRepo) ListReviewers(ctx context.Context, prID string) ([]string, error) {
	query := `SELECT r.reviewer_id FROM reviewers r WHERE r.pr_id=$1`

	rows, err := r.db.QueryContext(ctx, query, prID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	revsList := make([]string, 0)
	for rows.Next() {
		var rev string
		err = rows.Scan(&rev)
		if err != nil {
			return nil, err
		}

		revsList = append(revsList, rev)
	}
	return revsList, rows.Err()
}

func NewPRRepo(db *sql.DB) PullReqRepository {
	return &pullReqRepo{db: db}
}
