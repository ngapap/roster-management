package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"roster-management/internal/models"
	jwtPkg "roster-management/pkg/jwt"
	"roster-management/pkg/util"
)

// RegisterWorker handles new worker registration
func (h *Handler) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, err.Error())
		return
	}

	// Ensure user is not an admin
	user.IsAdmin = false

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "Error processing password")
		return
	}

	user.ID = uuid.New().String()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC()

	if err := h.repo.CreateUser(ctx, user); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "Error creating user")
		return
	}

	// Don't send password back
	user.Password = ""

	util.SendResponse(w, http.StatusOK, user, nil)
	logrus.Printf("successfully creates user with an email: %s ", user.Email)
}

// Login handles user authentication
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusBadRequest, nil, err)
		return
	}

	user, err := h.repo.GetUserByEmail(ctx, loginReq.Email)
	if errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
		util.SendResponse(w, http.StatusUnauthorized, nil, "user not found")
		return
	} else if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "database error")
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusUnauthorized, nil, "invalid credentials")
		return
	}

	// Don't send password back
	user.Password = ""

	// Generate JWT token
	jwtExpiry := h.jwtCfg.Exp
	if jwtExpiry == 0 {
		jwtExpiry = 60
	}
	expAt := time.Now().Add(time.Duration(jwtExpiry) * time.Minute)
	claims := jwt.MapClaims{
		"email": user.Email,
		"id":    user.ID,
		"exp":   expAt.Unix(),
	}

	jwtToken, err := jwtPkg.CreateToken(claims, h.jwtCfg.Key)
	if err != nil {
		logrus.Error(err)
		util.SendResponse(w, http.StatusInternalServerError, nil, "unable to login")
		return
	}

	response := models.LoginResponse{
		Token:     jwtToken,
		User:      *user,
		ExpiresAt: expAt,
	}

	util.SendResponse(w, http.StatusOK, response, nil)
	logrus.Printf("successfully login user with an email: %s ", user.Email)
}
