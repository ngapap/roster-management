package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"roster-management/internal/models"
	"roster-management/pkg/postgres"
)

func (r *Repository) CreateShiftRequest(ctx context.Context, request *models.ShiftRequest) (string, error) {
	query := `
		INSERT INTO shift_requests (shift_id, worker_id, status)
		VALUES ($1, $2, $3)  RETURNING id;
	`

	row := r.db.QueryRowContext(ctx, query,
		request.ShiftID, request.WorkerID, request.Status)

	var ID string
	if err := row.Scan(&ID); err != nil {
		return "", err
	}

	return ID, nil
}

func (r *Repository) UpdateShiftRequest(ctx context.Context, req *models.ShiftRequest) error {

	query := `
		UPDATE shift_requests
		SET  status = $1
		WHERE id = $2
	`
	_, err := r.db.ExecContext(ctx, query, req.Status, req.ID)
	if err != nil {
		return err
	}

	if req.Status == models.ApprovedShiftRequest {
		query := `
		UPDATE shift_requests
		SET  status = $1
		WHERE id != $2 AND shift_id = $3 
	`
		_, err := r.db.ExecContext(ctx, query, models.NotSelectedShiftRequest, req.ID, req.ShiftID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Repository) GetShiftRequests(ctx context.Context, opts ...models.ShiftRequestFilterOption) ([]*models.ShiftRequest, error) {
	filter := new(models.ShiftRequestFilter)
	for _, opt := range opts {
		opt(filter)
	}

	query, params := r.buildGetShiftRequestFilter(filter)

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var res []*models.ShiftRequest
	for rows.Next() {
		var item models.ShiftRequest
		if err := rows.Scan(
			&item.ID,
			&item.ShiftID,
			&item.WorkerID,
			&item.Status,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		res = append(res, &item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (r *Repository) buildGetShiftRequestFilter(filter *models.ShiftRequestFilter) (string, map[string]interface{}) {
	params := map[string]interface{}{}
	query := bytes.NewBufferString(`SELECT id, shift_id, worker_id, status, created_at, updated_at FROM shift_requests`)

	conds := postgres.Conditions{}
	if filter.ID != "" {
		colName := "id"
		conds = append(conds, postgres.Condition{
			Field:    colName,
			Operator: postgres.OperatorEqual,
			Value:    fmt.Sprintf(":%s", colName),
		})
		params[colName] = filter.ID
	} else {
		if filter.ShiftID != "" {
			colName := "shift_id"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})
			params[colName] = filter.ShiftID
		}

		if filter.WorkerID != "" {
			colName := "worker_id"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})
			params[colName] = filter.WorkerID
		}

		if len(filter.Status) > 0 {
			var condVal string
			colName := "status"
			sLen := len(filter.Status)
			for i, v := range filter.Status {
				key := fmt.Sprintf("status_%d", i)
				condVal += fmt.Sprintf(":%s", key)
				params[key] = v
				if i != sLen-1 {
					condVal += ","
				}
			}
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorIn,
				Value:    condVal,
			})
		}

	}

	for k, cond := range conds {
		if k == 0 {
			query.WriteString(" WHERE ")
		} else {
			query.WriteString(" AND ")
		}
		format := " %s %s %s "
		if cond.Operator == postgres.OperatorIn {
			format = " %s %s (%s) "
		}

		query.WriteString(fmt.Sprintf(format, cond.Field, cond.Operator, cond.Value))
	}

	query.WriteString(fmt.Sprintf(" ORDER BY created_at %s ", postgres.OrderByDateAscending.String()))

	return query.String(), params
}

func (r *Repository) DeleteShiftRequest(ctx context.Context, reqID string) error {
	query := `DELETE FROM shift_requests WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, reqID)
	return err
}

func (r *Repository) GetShiftRequestByID(ctx context.Context, ID string) (*models.ShiftRequest, error) {
	query := `
		SELECT ID, shift_id, worker_id, status, created_at
		FROM shift_requests WHERE id = $1
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
