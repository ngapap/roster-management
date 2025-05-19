package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"

	"roster-management/internal/models"
	jwtPkg "roster-management/pkg/jwt"
)

// RegisterWorker handles new worker registration
func (h *Handler) RegisterWorker(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user := new(models.User)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Ensure user is not an admin
	user.IsAdmin = false

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	user.ID = uuid.New().String()
	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now().UTC()

	err = h.repo.CreateUser(ctx, user)

	if err != nil {
		logrus.Error(err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Don't send password back
	user.Password = ""

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(user)

	logrus.Printf("successfully creates user with an email: %s ", user.Email)
}

// Login handles user authentication
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var loginReq models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		logrus.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByEmail(ctx, loginReq.Email)
	if errors.Is(err, sql.ErrNoRows) {
		logrus.Error(err)
		http.Error(w, "user not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		logrus.Error(err)
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginReq.Password)); err != nil {
		logrus.Error(err)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	// hide password
	user.Password = ""

	// Generate JWT token
	jwtExpiry, _ := strconv.Atoi(os.Getenv("JWT_EXPIRY"))
	if jwtExpiry == 0 {
		jwtExpiry = 60
	}
	expAt := time.Now().Add(time.Duration(jwtExpiry) * time.Minute)
	claims := jwt.MapClaims{
		"email":     user.Email,
		"id":        user.ID,
		"expire_at": expAt.Format(time.RFC3339Nano),
	}

	jwtToken, err := jwtPkg.CreateToken(claims, h.jwtCfg.Key)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "unable to login", http.StatusInternalServerError)
		return
	}

	response := models.LoginResponse{
		Token:     jwtToken,
		User:      *user,
		ExpiresAt: expAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	logrus.Printf("successfully login user with an email: %s ", user.Email)
}

// AuthMiddleware checks for valid JWT token and user role
//func AuthMiddleware(next http.HandlerFunc, requireAdmin bool) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//		token := r.Header.Get("Authorization")
//		if token == "" {
//			http.Error(w, "Unauthorized", http.StatusUnauthorized)
//			return
//		}
//
//		// Remove "Bearer " prefix if present
//		if len(token) > 7 && token[:7] == "Bearer " {
//			token = token[7:]
//		}
//
//		claims, err := validateJWT(token)
//		if err != nil {
//			http.Error(w, "Invalid token", http.StatusUnauthorized)
//			return
//		}
//
//		// Check if user exists and has correct role
//		var isAdmin bool
//		err = db.DB.QueryRow("SELECT is_admin FROM users WHERE id = $1", claims.UserID).Scan(&isAdmin)
//		if err != nil {
//			http.Error(w, "User not found", http.StatusUnauthorized)
//			return
//		}
//
//		if requireAdmin && !isAdmin {
//			http.Error(w, "Admin access required", http.StatusForbidden)
//			return
//		}
//
//		// Add user info to request context
//		ctx := context.WithValue(r.Context(), "userID", claims.UserID)
//		ctx = context.WithValue(ctx, "isAdmin", isAdmin)
//		next.ServeHTTP(w, r.WithContext(ctx))
//	}
//}
