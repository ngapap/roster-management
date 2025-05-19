package repository

import (
	"context"

	"roster-management/internal/models"
)

type Repository interface {
	UserRepository
	ShiftRepository
	ShiftRequestRepository
}

// UserRepository defines the interface for user-related database operations
type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
}

// ShiftRepository defines the interface for shift-related database operations
type ShiftRepository interface {
	CreateShift(ctx context.Context, shift *models.Shift) (string, error)
	GetShifts(ctx context.Context, opts ...models.ShiftFilterOption) ([]*models.Shift, error)
	UpdateShift(ctx context.Context, shift *models.Shift) error
	DeleteShift(ctx context.Context, id string) error
}

// ShiftRequestRepository defines the interface for shift request operations
type ShiftRequestRepository interface {
	CreateShiftRequest(ctx context.Context, request *models.ShiftRequest) error
	GetShiftRequestByID(ctx context.Context, id string) (*models.ShiftRequest, error)
	GetPendingShiftRequest(ctx context.Context) ([]*models.ShiftRequest, error)
	GetShiftRequestByShift(ctx context.Context, shiftID string) ([]*models.ShiftRequest, error)
	GetShiftRequestByWorker(ctx context.Context, workerID string) ([]*models.ShiftRequest, error)
}

// WorkerAvailabilityRepository defines the interface for worker availability operations
//type WorkerAvailabilityRepository interface {
//	Create(ctx context.Context, availability *models.WorkerAvailability) error
//	GetByWorker(ctx context.Context, workerID string) ([]*models.WorkerAvailability, error)
//	Update(ctx context.Context, availability *models.WorkerAvailability) error
//	Delete(ctx context.Context, id string) error
//}
