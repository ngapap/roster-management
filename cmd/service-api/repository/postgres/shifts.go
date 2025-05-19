package postgres

import (
	"context"
	"database/sql"

	"roster-management/internal/models"
)

func (r *Repository) CreateShift(ctx context.Context, shift *models.Shift) error {
	query := `
		INSERT INTO shifts (id, start_time, end_time, role,  created_at,  is_available)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		shift.ID, shift.StartTime, shift.EndTime,
		shift.Role, shift.CreatedAt, shift.IsAvailable)

	return err
}

func (r *Repository) GetShiftByID(ctx context.Context, id string) (*models.Shift, error) {
	query := `
		SELECT id, date, start_time, end_time, role, location, created_at, created_by, assigned_to, is_available
		FROM shifts WHERE id = $1
	`
	shift := &models.Shift{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&shift.ID, &shift.StartTime, &shift.EndTime,
		&shift.Role, &shift.CreatedAt,
		&shift.AssignedTo, &shift.IsAvailable)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return shift, err
}

func (r *Repository) GetAvailableShift(ctx context.Context) ([]*models.Shift, error) {
	query := `
		SELECT id, date, start_time, end_time, role, location, created_at, created_by, assigned_to, is_available
		FROM shifts WHERE is_available = true
		ORDER BY date, start_time
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []*models.Shift
	for rows.Next() {
		shift := &models.Shift{}
		err := rows.Scan(
			&shift.ID, &shift.StartTime, &shift.EndTime,
			&shift.Role, &shift.CreatedAt,
			&shift.AssignedTo, &shift.IsAvailable)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (r *Repository) GetShiftByWorker(ctx context.Context, workerID string) ([]*models.Shift, error) {
	query := `
		SELECT id, date, start_time, end_time, role, location, created_at, created_by, assigned_to, is_available
		FROM shifts WHERE assigned_to = $1
		ORDER BY date, start_time
	`
	rows, err := r.db.QueryContext(ctx, query, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shifts []*models.Shift
	for rows.Next() {
		shift := &models.Shift{}
		err := rows.Scan(
			&shift.ID, &shift.StartTime, &shift.EndTime,
			&shift.Role, &shift.CreatedAt,
			&shift.AssignedTo, &shift.IsAvailable)
		if err != nil {
			return nil, err
		}
		shifts = append(shifts, shift)
	}
	return shifts, nil
}

func (r *Repository) UpdateShift(ctx context.Context, shift *models.Shift) error {
	query := `
		UPDATE shifts
		SET  start_time = $1, end_time = $2, role = $3, 
			assigned_to = $4, is_available = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		shift.StartTime, shift.EndTime, shift.Role,
		shift.AssignedTo, shift.IsAvailable, shift.ID)
	return err
}

func (r *Repository) Delete(ctx context.Context, ID string) error {
	query := `DELETE FROM shifts WHERE ID = $1`
	_, err := r.db.ExecContext(ctx, query, ID)
	return err
}
