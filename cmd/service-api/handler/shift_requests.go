package handler

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"net/http"
	"roster-management/internal/models"
)

func (h *Handler) CreateShiftRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	shift := new(models.CreateShiftRequestPayload)
	if err := json.NewDecoder(r.Body).Decode(&shift); err != nil {
		logrus.Error(err)
		http.Error(w, "err when parsing payload", http.StatusBadRequest)
		return
	}
	res := &models.ShiftRequest{
		ShiftID:  shift.ShiftID,
		WorkerID: shift.WorkerID,
		Status:   models.PendingShiftRequest,
	}
	reqID, err := h.repo.CreateShiftRequest(ctx, res)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when creating shift request", http.StatusInternalServerError)
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

	logrus.Println("successfully updates shift request")
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
		http.Error(w, "err when parsing payld", http.StatusBadRequest)
		return
	}

	requests, err := h.repo.GetShiftRequests(ctx, models.WithRequestID(reqID.String()))
	if err != nil {
		logrus.Error(err)
		http.Error(w, "err when parsing payld", http.StatusBadRequest)
		return
	}

	if len(requests) < 1 {
		http.Error(w, "no requests found", http.StatusNotFound)
		return
	}

	res := requests[0]
	res.Status = payload.Status
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
