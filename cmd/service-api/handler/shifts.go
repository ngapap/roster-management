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
	"roster-management/pkg/util"
)

// CreateShift creates a new shift, shift must be on 30 minutes interval and the duration must be between 4 and 12 hours.
func (h *Handler) CreateShift(w http.ResponseWriter, r *http.Request) {
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

	shift := new(models.CreateShiftPayload)
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when parsing payld")
		return
	}

	if err := validateShiftDuration(shift.StartTime, shift.EndTime); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, err)
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
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when creating shift")
		return
	}

	shifts, err := h.repo.GetShifts(ctx, models.WithID(shiftID))
	if err == nil && len(shifts) > 0 {
		res = shifts[0]
	}

	util.SendResponse(w, http.StatusOK, res, nil)
	logrus.Println("successfully creates shift")
}

func (h *Handler) UpdateShift(w http.ResponseWriter, r *http.Request) {
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

	shift := new(models.UpdateShiftPayload)
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "err when parsing payld")
		return
	}

	if err := validateShiftDuration(shift.StartTime, shift.EndTime); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	shiftID, err := uuid.Parse(chi.URLParam(r, "shiftID"))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "invalid uuid")
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
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when updating shift")
		return
	}

	shifts, err := h.repo.GetShifts(ctx, models.WithID(res.ID))
	if err == nil && len(shifts) > 0 {
		res = shifts[0]
	}

	util.SendResponse(w, http.StatusOK, res, nil)
	logrus.Println("successfully updates shift")
}

func (h *Handler) DeleteShift(w http.ResponseWriter, r *http.Request) {
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

	shiftID, err := uuid.Parse(chi.URLParam(r, "shiftID"))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "invalid uuid")
		return
	}

	if err := h.repo.DeleteShift(ctx, shiftID.String()); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when deleting shift")
		return
	}

	util.SendResponse(w, http.StatusOK, "data has been deleted", nil)
	logrus.Println("successfully deletes shift")
}

// GetAvailableShifts will return all of unassigned shifts
func (h *Handler) GetAvailableShifts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	now := time.Now().UTC()

	shifts, err := h.repo.GetShifts(ctx, models.WithIsAvailable(models.TrueStr), models.WithStartTime(now))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetching available shifts")
		return
	}

	util.SendResponse(w, http.StatusOK, shifts, nil)
	logrus.Println("successfully fetch available shifts")
}

// GetAssignedShifts will return all of assigned shifts
func (h *Handler) GetAssignedShifts(w http.ResponseWriter, r *http.Request) {
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

	now := time.Now().UTC()

	shifts, err := h.repo.GetShifts(ctx, models.WithIsAvailable(models.FalseStr), models.WithStartTime(now))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetching assigned shifts")
		return
	}

	util.SendResponse(w, http.StatusOK, shifts, nil)
	logrus.Println("successfully fetch assigned shifts")
}

func (h *Handler) GetShiftByWorker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	workerID, err := uuid.Parse(chi.URLParam(r, "workerID"))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, "invalid uuid")
		return
	}

	shifts, err := h.repo.GetShifts(ctx, models.WithAssignedTo(workerID.String()))
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "err when fetching shift by worker")
		return
	}

	util.SendResponse(w, http.StatusOK, shifts, nil)
	logrus.Println("successfully fetch shift by worker")
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
