package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"roster-management/internal/models"
	"time"
)

// CreateShift creates a new shift, shift must be on 30 minutes interval and the duration must be between 4 and 12 hours.
func (h *Handler) CreateShift(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now()

	shift := new(models.Shift)
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if shift.StartTime.Sub(now).Minutes() < 30 {
		http.Error(w, "shift start_time must be at least 30 minutes from now", http.StatusBadRequest)
		return
	}

	isOn30MinInterval := func(t time.Time) bool {
		return t.Minute()%30 == 0 && t.Second() == 0 && t.Nanosecond() == 0
	}

	if !isOn30MinInterval(shift.StartTime) || !isOn30MinInterval(shift.EndTime) {
		http.Error(w, "shift start and end times must be on 30-minute intervals", http.StatusBadRequest)
		return
	}

	duration := shift.EndTime.Sub(shift.StartTime)
	if duration.Minutes() < 0 {
		http.Error(w, "shift start_time must be lesser than end_time", http.StatusBadRequest)
		return
	}

	if duration.Hours() < 4 {
		http.Error(w, "shift must be greater than equal 4h", http.StatusBadRequest)
		return
	}

	if duration.Hours() > 12 {
		http.Error(w, "shift must be lest than equal 12h", http.StatusBadRequest)
		return
	}

	shift.IsAvailable = true
	if shift.AssignedTo != "" {
		shift.IsAvailable = false
	}

	if err := h.repo.CreateShift(ctx, shift); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shift)

	logrus.Println("successfully creates shift")
}

func (h *Handler) GetAvailableShifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shifts, err := h.repo.GetShifts(ctx, models.WithIsAvailable(models.TrueStr))
	if err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(shifts)
}
