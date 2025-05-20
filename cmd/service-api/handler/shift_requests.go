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
	"roster-management/pkg/util"
)

func (h *Handler) CreateShiftRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetch user info")
		return
	}

	payload := new(models.CreateShiftRequestPayload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when parsing payload")
		return
	}

	// check shift availability
	shifts, err := h.repo.GetShifts(ctx, models.WithID(payload.ShiftID))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when getting shifts")
		return
	}

	if len(shifts) == 0 {
		util.SendResponse(w, http.StatusBadRequest, nil, "shift not found")
		return
	}

	shift := shifts[0]
	// bypass shift availability when requesting as admin for further approval
	if !shift.IsAvailable && !user.IsAdmin {
		util.SendResponse(w, http.StatusForbidden, nil, "shift already assigned to someone else")
		return
	}

	// check worker latest shift
	lastShift, err := h.repo.GetLastShiftByWorker(ctx, payload.WorkerID)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when getting last shift")
		return
	}

	// check weekly worker quota
	weeklyShift, err := h.repo.CountWeeklyShiftByWorker(ctx, payload.WorkerID)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetching weekly shift")
		return
	}

	if err := validateShiftRequest(shift, lastShift, weeklyShift); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, err)
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
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when creating payload request")
		return
	}
	res.ID = reqID

	requests, err := h.repo.GetShiftRequests(ctx, models.WithRequestID(res.ID))
	if err == nil && len(requests) > 0 {
		res = requests[0]
	}

	util.SendResponse(w, http.StatusOK, res, nil)
	logrus.Println("successfully updates shift request")
}

func (h *Handler) UpdateShiftRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetch user info")
		return
	}
	if !user.IsAdmin {
		util.SendResponse(w, http.StatusUnauthorized, nil, "require admin access")
		return
	}

	reqID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "invalid uuid")
		return
	}

	payload := new(models.UpdateShiftRequestPayload)
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when parsing payload")
		return
	}

	requests, err := h.repo.GetShiftRequests(ctx, models.WithRequestID(reqID.String()))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when fetching shift requests")
		return
	}

	if len(requests) < 1 {
		util.SendResponse(w, http.StatusNotFound, nil, "shift request not found")
		return
	}

	res := requests[0]
	res.Status = payload.Status

	if payload.Status == models.ApprovedShiftRequest {

		shifts, err := h.repo.GetShifts(ctx, models.WithID(res.ShiftID))
		if err != nil {
			logrus.Error(err)
			util.SendResponse(w, http.StatusInternalServerError, nil, "err when getting shifts")
			return
		}
		if len(shifts) > 0 {
			shift := shifts[0]

			// check worker latest shift
			lastShift, err := h.repo.GetLastShiftByWorker(ctx, res.WorkerID)
			if err != nil && !errors.Is(err, sql.ErrNoRows) {
				logrus.Error(err)
				util.SendResponse(w, http.StatusInternalServerError, nil, "err when getting last shift")
				return
			}

			weeklyShift, err := h.repo.CountWeeklyShiftByWorker(ctx, res.WorkerID)
			if err != nil {
				logrus.Error(err)
				util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetching weekly shift")
				return
			}

			if err := validateShiftRequest(shift, lastShift, weeklyShift); err != nil {
				logrus.Error(err)
				util.SendResponse(w, http.StatusBadRequest, nil, err)
				return
			}

			shift.AssignedTo = res.WorkerID
			shift.IsAvailable = false
			if err := h.repo.UpdateShift(ctx, shift); err != nil {
				logrus.Error(err)
				util.SendResponse(w, http.StatusInternalServerError, nil, "err when updating shift")
				return
			}
		}
	}

	if err := h.repo.UpdateShiftRequest(ctx, res); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when updating shift")
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
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetch user info")
		return
	}
	if !user.IsAdmin {
		util.SendResponse(w, http.StatusUnauthorized, nil, "require admin access")
		return
	}

	reqID, err := uuid.Parse(chi.URLParam(r, "requestID"))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "invalid uuid")
		return
	}

	if err := h.repo.DeleteShiftRequest(ctx, reqID.String()); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when deleting shift")
		return
	}

	util.SendResponse(w, http.StatusOK, "data has been deleted", nil)

	logrus.Println("successfully deletes shift request")
}

func (h *Handler) GetShiftRequestByWorker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workerID, err := uuid.Parse(chi.URLParam(r, "workerID"))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "invalid uuid")
		return
	}

	requests, err := h.repo.GetShiftRequests(ctx, models.WithWorkerID(workerID.String()))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when fetching shift requests")
		return
	}

	util.SendResponse(w, http.StatusOK, requests, nil)
	logrus.Println("successfully fetch  shift request")
}

func (h *Handler) GetPendingShiftRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userID := ctx.Value("user_id").(string)
	user, err := h.repo.GetUserByID(ctx, userID)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetch user info")
		return
	}

	opts := []models.ShiftRequestFilterOption{
		models.WithStatus([]string{string(models.PendingShiftRequest)}...),
	}
	if !user.IsAdmin {
		opts = append(opts, models.WithWorkerID(userID))
	}

	requests, err := h.repo.GetShiftRequests(ctx, opts...)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when fetching shift requests")
		return
	}

	util.SendResponse(w, http.StatusOK, requests, nil)
	logrus.Println("successfully fetch shift request")
}

// validateShiftRequest will check shift availability, daily and weekly worker quota, and overlapped shift
func validateShiftRequest(reqShift, lastShift *models.Shift, weeklyQuota int) error {

	sameDay := func(t1, t2 time.Time) bool {
		y1, m1, d1 := t1.UTC().Date()
		y2, m2, d2 := t2.UTC().Date()
		return y1 == y2 && m1 == m2 && d1 == d2
	}

	if sameDay(reqShift.StartTime, lastShift.StartTime) {
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
