package models

import (
	"time"
)

type Shift struct {
	ID          string    `json:"id" db:"id"`
	StartTime   time.Time `json:"start_time" db:"start_time"`
	EndTime     time.Time `json:"end_time" db:"end_time"`
	Role        string    `json:"role" db:"role"`
	AssignedTo  string    `json:"assigned_to,omitempty" db:"assigned_to"`
	IsAvailable bool      `json:"is_available" db:"is_available"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type ShiftRequest struct {
	ID        string    `json:"id" db:"id"`
	ShiftID   string    `json:"shift_id" db:"shift_id"`
	WorkerID  string    `json:"worker_id" db:"worker_id"`
	Status    string    `json:"status" db:"status"` // pending, approved, rejected
	Note      string    `json:"note" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
