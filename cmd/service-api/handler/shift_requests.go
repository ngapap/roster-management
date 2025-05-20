package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"roster-management/internal/models"
)

func (h *Handler) CreateShiftRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when fetch user info", http.StatusInternalServerError)
		return
	}

	payload := new(models.CreateShiftRequestPayload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.Error(err)
		http.Error(w, "err when parsing payload", http.StatusBadRequest)
		return
	}

	// check shift availability
	shifts, err := h.repo.GetShifts(ctx, models.WithID(payload.ShiftID))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when getting shifts", http.StatusInternalServerError)
		return
	}

	if len(shifts) == 0 {
		http.Error(w, "shift not found", http.StatusBadRequest)
		return
	}

	shift := shifts[0]
	// bypass shift availability when requesting as admin for further approval
	if !shift.IsAvailable && !user.IsAdmin {
		http.Error(w, "shift already assigned to someone else", http.StatusForbidden)
		return
	}

	// check worker latest shift
	lastShift, err := h.repo.GetLastShiftByWorker(ctx, payload.WorkerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
		http.Error(w, "err when getting last shift", http.StatusInternalServerError)
		return
	}

	// check weekly worker quota
	weeklyShift, err := h.repo.CountWeeklyShiftByWorker(ctx, payload.WorkerID)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when fetching weekly shift", http.StatusInternalServerError)
		return
	}

	if err := validateShiftRequest(shift, lastShift, weeklyShift); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res := &models.ShiftRequest{
		ShiftID:  payload.ShiftID,
		WorkerID: payload.WorkerID,
		Status:   models.PendingShiftRequest,
	}
	reqID, err := h.repo.CreateShiftRequest(ctx, res)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when creating payload request", http.StatusInternalServerError)
		return
	}
	res.ID = reqID

	requests, err := h.repo.GetShiftRequests(ctx, models.WithRequestID(res.ID))
	if err == nil && len(requests) > 0 {
		res = requests[0]
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	logrus.Println("successfully updates payload request")
}

func (h *Handler) UpdateShiftRequest(w http.ResponseWriter, r *http.Request) {
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

	reqID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	payload := new(models.UpdateShiftRequestPayload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.Error(err)
		http.Error(w, "err when parsing payload", http.StatusBadRequest)
		return
	}

	requests, err := h.repo.GetShiftRequests(ctx, models.WithRequestID(reqID.String()))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when fetching shift requests", http.StatusBadRequest)
		return
	}

	if len(requests) < 1 {
		http.Error(w, "shift request not found", http.StatusNotFound)
		return
	}

	res := requests[0]
	res.Status = payload.Status

	if payload.Status == models.ApprovedShiftRequest {

		shifts, err := h.repo.GetShifts(ctx, models.WithID(res.ShiftID))
		if err != nil {
			logrus.Error(err)
			http.Error(w, "err when getting shifts", http.StatusInternalServerError)
			return
		}
		if len(shifts) > 0 {
			shift := shifts[0]

			// check worker latest shift
			lastShift, err := h.repo.GetLastShiftByWorker(ctx, res.WorkerID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				logrus.Error(err)
				http.Error(w, "err when getting last shift", http.StatusInternalServerError)
				return
			}

			weeklyShift, err := h.repo.CountWeeklyShiftByWorker(ctx, res.WorkerID)
			if err != nil {
				logrus.Error(err)
				http.Error(w, "err when fetching weekly shift", http.StatusInternalServerError)
				return
			}

			if err := validateShiftRequest(shift, lastShift, weeklyShift); err != nil {
				logrus.Error(err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			shift.AssignedTo = res.WorkerID
			shift.IsAvailable = false
			if err := h.repo.UpdateShift(ctx, shift); err != nil {
				logrus.Error(err)
				http.Error(w, "err when updating shift", http.StatusInternalServerError)
				return
			}
		}
	}

	if err := h.repo.UpdateShiftRequest(ctx, res); err != nil {
		logrus.Error(err)
		http.Error(w, "err when updating shift", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)

	logrus.Println("successfully updates shift request")

}

func (h *Handler) DeleteShiftRequest(w http.ResponseWriter, r *http.Request) {
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

	reqID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "invalid uuid", http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteShiftRequest(ctx, reqID.String()); err != nil {
		logrus.Error(err)
		http.Error(w, "err when deleting shift", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("data has been deleted")

	logrus.Println("successfully deletes shift request")
}

func (h *Handler) GetShiftRequestByWorker(w http.ResponseWriter, r *http.Request) {}

// validateShiftRequest will check shift availability, daily and weekly worker quota, and overlapped shift
func validateShiftRequest(reqShift, lastShift *models.Shift, weeklyQuota int) error {

	sameDay := func(t1, t2 time.Time) bool {
		y1, m1, d1 := t1.UTC().Date()
		y2, m2, d2 := t2.UTC().Date()
		return y1 == y2 && m1 == m2 && d1 == d2
	}

	if sameDay(reqShift.StartTime, lastShift.EndTime) {
		return errors.New("maximum one shift per day")
	}

	if weeklyQuota >= 5 {
		return errors.New("worker reached weekly quota")
	}

	if reqShift.StartTime.Before(lastShift.EndTime) {
		return errors.New("shift is overlapped")
	}

	return nil
}
