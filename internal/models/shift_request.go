package models

import "time"

type ShiftRequest struct {
	ID        string    `json:"id" db:"id"`
	ShiftID   string    `json:"shift_id" db:"shift_id"`
	WorkerID  string    `json:"worker_id" db:"worker_id"`
	Status    string    `json:"status" db:"status"` // pending, approved, rejected
	Note      string    `json:"note" db:"note"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
