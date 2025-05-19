package models

import (
	"time"
)

type StrBool string

const (
	EmptyStrBool StrBool = ""
	TrueStr      StrBool = "true"
	FalseStr     StrBool = "false"
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

type CreateShiftPayload struct {
	StartTime  time.Time `json:"start_time" `
	EndTime    time.Time `json:"end_time" `
	Role       string    `json:"role" `
	AssignedTo string    `json:"assigned_to,omitempty" `
}

type UpdateShiftPayload struct {
	StartTime  time.Time `json:"start_time" db:"start_time"`
	EndTime    time.Time `json:"end_time" db:"end_time"`
	Role       string    `json:"role" db:"role"`
	AssignedTo string    `json:"assigned_to,omitempty" db:"assigned_to"`
}

type ShiftFilter struct {
	ID          string    `json:"id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Role        string    `json:"role"`
	AssignedTo  string    `json:"assigned_to"`
	IsAvailable StrBool   `json:"is_available"`
}

type ShiftFilterOption func(*ShiftFilter)

func WithID(id string) ShiftFilterOption {
	return func(filter *ShiftFilter) {
		filter.ID = id
	}
}

func WithStartTime(startTime time.Time) ShiftFilterOption {
	return func(filter *ShiftFilter) {
		filter.StartTime = startTime
	}
}

func WithEndTime(endTime time.Time) ShiftFilterOption {
	return func(filter *ShiftFilter) {
		filter.EndTime = endTime
	}
}

func WithRole(role string) ShiftFilterOption {
	return func(filter *ShiftFilter) {
		filter.Role = role
	}
}

func WithAssignedTo(assignedTo string) ShiftFilterOption {
	return func(filter *ShiftFilter) {
		filter.AssignedTo = assignedTo
	}
}

func WithIsAvailable(isAvailable StrBool) ShiftFilterOption {
	return func(filter *ShiftFilter) {
		filter.IsAvailable = isAvailable
	}
}
