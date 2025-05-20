package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"roster-management/internal/models"
)

// CreateShift creates a new shift, shift must be on 30 minutes interval and the duration must be between 4 and 12 hours.
func (h *Handler) CreateShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when fetch user info", http.StatusInternalServerError)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "require admin access", http.StatusUnauthorized)
		return
	}

	shift := new(models.CreateShiftPayload)
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		logrus.Error(err)
		http.Error(w, "err when parsing payld", http.StatusBadRequest)
		return
	}

	if err := validateShiftDuration(shift.StartTime, shift.EndTime); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isAvailable := true
	if shift.AssignedTo != "" {
		isAvailable = false
	}

	res := &models.Shift{
		StartTime:   shift.StartTime,
		EndTime:     shift.EndTime,
		Role:        shift.Role,
		AssignedTo:  shift.AssignedTo,
		IsAvailable: isAvailable,
	}

	shiftID, err := h.repo.CreateShift(ctx, res)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when creating shift", http.StatusInternalServerError)
		return
	}

	shifts, err := h.repo.GetShifts(ctx, models.WithID(shiftID))
	if err == nil && len(shifts) > 0 {
		res = shifts[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	logrus.Println("successfully creates shift")
}

func (h *Handler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when fetch user info", http.StatusInternalServerError)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "require admin access", http.StatusUnauthorized)
		return
	}

	shift := new(models.UpdateShiftPayload)
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		logrus.Error(err)
		http.Error(w, "err when parsing payld", http.StatusBadRequest)
		return
	}

	if err := validateShiftDuration(shift.StartTime, shift.EndTime); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	shiftID, err := uuid.Parse(chi.URLParam(r, "shiftID"))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	isAvailable := true
	if shift.AssignedTo != "" {
		isAvailable = false
	}

	res := &models.Shift{
		ID:          shiftID.String(),
		StartTime:   shift.StartTime,
		EndTime:     shift.EndTime,
		Role:        shift.Role,
		AssignedTo:  shift.AssignedTo,
		IsAvailable: isAvailable,
	}
	if err := h.repo.UpdateShift(ctx, res); err != nil {
		logrus.Error(err)
		http.Error(w, "err when updating shift", http.StatusInternalServerError)
		return
	}

	shifts, err := h.repo.GetShifts(ctx, models.WithID(res.ID))
	if err == nil && len(shifts) > 0 {
		res = shifts[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	logrus.Println("successfully updates shift")
}

func (h *Handler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when fetch user info", http.StatusInternalServerError)
		return
	}
	if !user.IsAdmin {
		http.Error(w, "require admin access", http.StatusUnauthorized)
		return
	}

	shiftID, err := uuid.Parse(chi.URLParam(r, "shiftID"))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteShift(ctx, shiftID.String()); err != nil {
		logrus.Error(err)
		http.Error(w, "err when deleting shift", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("data has been deleted")

	logrus.Println("successfully deletes shift")
}

// GetAvailableShifts will return all of unassigned shifts
func (h *Handler) GetAvailableShifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now().UTC()

	shifts, err := h.repo.GetShifts(ctx, models.WithIsAvailable(models.TrueStr), models.WithStartTime(now))
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shifts)
}

func (h *Handler) GetShiftByWorker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workerID, err := uuid.Parse(chi.URLParam(r, "workerID"))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	shifts, err := h.repo.GetShifts(ctx, models.WithAssignedTo(workerID.String()))
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(shifts)
}

func validateShiftDuration(start, end time.Time) error {
	now := time.Now().UTC()

	if start.Sub(now).Minutes() < 30 {
		return errors.New("shift start_time must be at least 30 minutes from now")
	}

	isOn30MinInterval := func(t time.Time) bool {
		return t.Minute()%30 == 0 && t.Second() == 0 && t.Nanosecond() == 0
	}

	if !isOn30MinInterval(start) || !isOn30MinInterval(end) {
		return errors.New("shift start and end times must be on 30-minute intervals")
	}

	duration := end.Sub(start)
	if duration.Minutes() < 0 {
		return errors.New("shift start time must be lesser than end_time")
	}

	if duration.Hours() < 4 {
		return errors.New("shift durations must be greater or equal 4h")
	}

	if duration.Hours() > 12 {
		return errors.New("shift start time must be greater or equal 12h")
	}

	return nil
}
