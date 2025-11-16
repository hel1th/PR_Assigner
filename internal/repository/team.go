package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hel1th/PR_Assigner/internal/apperrors"
	"github.com/hel1th/PR_Assigner/internal/domain"
)

type TeamRepository interface {
	GetTeam(ctx context.Context, teamName string) (*domain.Team, error)
	CreateTeam(ctx context.Context, teamName string, Members []domain.TeamMember) error
	GetActiveTeamMembers(ctx context.Context, authorID string) ([]string, error)
}
type teamRepo struct {
	db *sql.DB
}

func (r *teamRepo) GetTeam(ctx context.Context, teamName string) (*domain.Team, error) {
	query := `SELECT u.id, u.username, u.is_active 
				FROM users u
				JOIN teams t ON t.name = u.team_name
				WHERE t.name = $1`

	rows, err := r.db.QueryContext(ctx, query, teamName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teamMembers := make([]*domain.TeamMember, 0)
	for rows.Next() {
		var m domain.TeamMember
		if err := rows.Scan(&m.ID, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		teamMembers = append(teamMembers, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(teamMembers) == 0 {
		var exists bool
		checkQuery := `SELECT EXISTS(SELECT 1 FROM teams WHERE name = $1)`
		if err := r.db.QueryRowContext(ctx, checkQuery, teamName).Scan(&exists); err != nil {
			return nil, err
		}
		if !exists {
			return nil, apperrors.NotFound
		}
		return &domain.Team{Name: teamName, Members: teamMembers}, nil
	}

	return &domain.Team{Name: teamName, Members: teamMembers}, nil
}

func (r *teamRepo) GetActiveTeamMembers(ctx context.Context, authorID string) (revsID []string, err error) {
	query := `
        SELECT u2.id 
        FROM users u1
        JOIN users u2 ON u2.team_name = u1.team_name
        WHERE u1.id = $1 
          AND u1.is_active = true
          AND u2.is_active = true 
          AND u2.id != $1
    `

	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]string, 0)
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		members = append(members, id)
	}
	if rows.Err() != nil {
		return nil, err
	}

	if len(members) == 0 {
		return nil, apperrors.NoCandidate
	}
	return members, nil
}

func (r *teamRepo) CreateTeam(ctx context.Context, teamName string, members []domain.TeamMember) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	queryTeam := `INSERT INTO teams (name) 
					VALUES ($1)
					ON CONFLICT (name) DO NOTHING`

	if _, err = tx.ExecContext(ctx, queryTeam, teamName); err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}

	if len(members) > 0 {
		queryMembers := `
			INSERT INTO users (id, username, team_name, is_active) 
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				team_name = excluded.team_name,
				username = excluded.username,
				is_active = excluded.is_active
		`

		stmt, err := tx.PrepareContext(ctx, queryMembers)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}
		defer stmt.Close()

		for _, member := range members {
			_, err = stmt.ExecContext(ctx, member.ID, member.Username, teamName, member.IsActive)
			if err != nil {
				return fmt.Errorf("failed to upsert member %s: %w", member.Username, err)
			}
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func NewTeamRepository(db *sql.DB) TeamRepository {
	return &teamRepo{db: db}
}
