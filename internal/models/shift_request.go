package models

import "time"

type ShiftRequestStatus string

const (
	ApprovedShiftRequest    ShiftRequestStatus = "approved"
	PendingShiftRequest     ShiftRequestStatus = "pending"
	RejectedShiftRequest    ShiftRequestStatus = "rejected"
	NotSelectedShiftRequest ShiftRequestStatus = "not_selected"
)

type ShiftRequest struct {
	ID        string             `json:"id" db:"id"`
	ShiftID   string             `json:"shift_id" db:"shift_id"`
	WorkerID  string             `json:"worker_id" db:"worker_id"`
	Status    ShiftRequestStatus `json:"status" db:"status"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
}

type CreateShiftRequestPayload struct {
	ShiftID  string `json:"shift_id"`
	WorkerID string `json:"worker_id"`
}

type UpdateShiftRequestPayload struct {
	Status ShiftRequestStatus `json:"status"`
}

type ShiftRequestFilter struct {
	ID       string
	ShiftID  string
	WorkerID string
	Status   []string
}

type ShiftRequestFilterOption func(*ShiftRequestFilter)

func WithRequestID(id string) ShiftRequestFilterOption {
	return func(f *ShiftRequestFilter) {
		f.ID = id
	}
}

func WithShiftID(id string) ShiftRequestFilterOption {
	return func(f *ShiftRequestFilter) {
		f.ShiftID = id
	}
}

func WithWorkerID(workerID string) ShiftRequestFilterOption {
	return func(f *ShiftRequestFilter) {
		f.WorkerID = workerID
	}
}

func WithStatus(status ...string) ShiftRequestFilterOption {
	return func(f *ShiftRequestFilter) {
		f.Status = status
	}
}
