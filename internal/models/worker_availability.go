package models

import "time"

type WorkerAvailability struct {
	ID        string    `json:"id" db:"id"`
	WorkerID  string    `json:"worker_id" db:"worker_id"`
	DayOfWeek int       `json:"day_of_week" db:"day_of_week"` // 0-6 (Sunday-Saturday)
	StartTime time.Time `json:"start_time" db:"start_time"`
	EndTime   time.Time `json:"end_time" db:"end_time"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
