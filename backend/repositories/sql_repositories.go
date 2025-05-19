package repositories

import (
	"context"
	"database/sql"

	"roster-management/backend/models"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) CreateShift(ctx context.Context, shift *models.Shift) (string, error) {
	query := `
		INSERT INTO shifts (id, start_at, end_at, role, assigned_to, created_at, is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id string
	err := r.db.QueryRowContext(ctx, query,
		shift.ID, shift.StartAt, shift.EndAt, shift.Role, shift.AssignedTo,
		shift.CreatedAt, shift.IsAvailable).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}
