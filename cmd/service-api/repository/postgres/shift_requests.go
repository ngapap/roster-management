package postgres

import (
	"context"
	"database/sql"
	"roster-management/internal/models"
)

func (r *Repository) CreateShiftRequest(ctx context.Context, request *models.ShiftRequest) error {
	query := `
		INSERT INTO shift_requests (id, shift_id, worker_id, status, created_at,)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(ctx, query,
		request.ID, request.ShiftID, request.WorkerID,
		request.Status, request.CreatedAt)

	return err
}

func (r *Repository) GetShiftRequestByID(ctx context.Context, ID string) (*models.ShiftRequest, error) {
	query := `
		SELECT ID, shift_id, worker_id, status, created_at
		FROM shift_requests WHERE ID = $1
	`
	request := &models.ShiftRequest{}
	err := r.db.QueryRowContext(ctx, query, ID).Scan(
		&request.ID, &request.ShiftID, &request.WorkerID,
		&request.Status, &request.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return request, err
}

func (r *Repository) GetPendingShiftRequest(ctx context.Context) ([]*models.ShiftRequest, error) {
	query := `
		SELECT id, shift_id, worker_id, status, created_at, updated_at
		FROM shift_requests WHERE status = 'pending'
		ORDER BY created_at
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ShiftRequest
	for rows.Next() {
		request := &models.ShiftRequest{}
		err := rows.Scan(
			&request.ID, &request.ShiftID, &request.WorkerID,
			&request.Status, &request.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (r *Repository) GetShiftRequestByShift(ctx context.Context, shiftID string) ([]*models.ShiftRequest, error) {
	query := `
		SELECT id, shift_id, worker_id, status, created_at, updated_at
		FROM shift_requests WHERE shift_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, shiftID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ShiftRequest
	for rows.Next() {
		request := &models.ShiftRequest{}
		err := rows.Scan(
			&request.ID, &request.ShiftID, &request.WorkerID,
			&request.Status, &request.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}

func (r *Repository) GetShiftRequestByWorker(ctx context.Context, workerID string) ([]*models.ShiftRequest, error) {
	query := `
		SELECT id, shift_id, worker_id, status, created_at, updated_at
		FROM shift_requests WHERE worker_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, workerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []*models.ShiftRequest
	for rows.Next() {
		request := &models.ShiftRequest{}
		err := rows.Scan(
			&request.ID, &request.ShiftID, &request.WorkerID,
			&request.Status, &request.CreatedAt)
		if err != nil {
			return nil, err
		}
		requests = append(requests, request)
	}
	return requests, nil
}
