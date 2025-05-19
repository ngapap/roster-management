package postgres

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"roster-management/internal/models"
	"roster-management/pkg/postgres"
)

func (r *Repository) CreateShift(ctx context.Context, shift *models.Shift) (string, error) {
	query := `
		INSERT INTO shifts (start_time, end_time, assigned_to, role,  is_available)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	workerID := sql.NullString{
		String: shift.AssignedTo,
		Valid:  true,
	}
	if err := uuid.Validate(shift.AssignedTo); err != nil {
		workerID.Valid = false
	}
	row := r.db.QueryRowContext(ctx, query,
		shift.StartTime, shift.EndTime, workerID,
		shift.Role, shift.IsAvailable)

	var ID string
	if err := row.Scan(&ID); err != nil {
		return "", err
	}

	return ID, nil
}

func (r *Repository) GetShifts(ctx context.Context, opts ...models.ShiftFilterOption) ([]*models.Shift, error) {
	filter := new(models.ShiftFilter)
	for _, opt := range opts {
		opt(filter)
	}

	query, params := r.buildGetShiftFilter(filter)

	rows, err := r.db.NamedQueryContext(ctx, query, params)
	if err != nil {
		return nil, err
	}
	defer func(rows *sqlx.Rows) {
		_ = rows.Close()
	}(rows)

	var res []*models.Shift
	for rows.Next() {
		var item models.Shift
		workerID := sql.NullString{}
		if err := rows.Scan(
			&item.ID,
			&item.StartTime,
			&item.EndTime,
			&item.Role,
			&workerID,
			&item.IsAvailable,
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

func (r *Repository) UpdateShift(ctx context.Context, shift *models.Shift) error {
	workerID := sql.NullString{
		String: shift.AssignedTo,
		Valid:  true,
	}
	if err := uuid.Validate(shift.AssignedTo); err != nil {
		workerID.Valid = false
	}

	query := `
		UPDATE shifts
		SET  start_time = $1, end_time = $2, role = $3, 
			assigned_to = $4, is_available = $5
		WHERE id = $6
	`
	_, err := r.db.ExecContext(ctx, query,
		shift.StartTime, shift.EndTime, shift.Role,
		workerID, shift.IsAvailable, shift.ID)
	return err
}

func (r *Repository) DeleteShift(ctx context.Context, ID string) error {
	query := `DELETE FROM shifts WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, ID)
	return err
}

func (r *Repository) buildGetShiftFilter(filter *models.ShiftFilter) (string, map[string]interface{}) {
	params := map[string]interface{}{}
	query := bytes.NewBufferString(`SELECT id, start_time, end_time, role, assigned_to, is_available, created_at, updated_at FROM shifts`)

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
		if !filter.StartTime.IsZero() {
			colName := "start_time"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorGreaterThanEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})

			params[colName] = filter.StartTime.Format(dateFormat)
		}

		if !filter.EndTime.IsZero() {
			colName := "end_time"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorLessThanEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})

			params[colName] = filter.EndTime.Format(dateFormat)
		}

		if filter.Role != "" {
			colName := "role"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})
			params[colName] = filter.Role
		}

		if filter.AssignedTo != "" {
			colName := "assigned_to"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})
			params[colName] = filter.AssignedTo
		}

		if filter.IsAvailable != models.EmptyStrBool {
			colName := "is_available"
			conds = append(conds, postgres.Condition{
				Field:    colName,
				Operator: postgres.OperatorEqual,
				Value:    fmt.Sprintf(":%s", colName),
			})
			params[colName] = filter.IsAvailable
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
