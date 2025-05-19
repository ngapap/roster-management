package middlewares

import (
	"context"
	"net/http"
	"roster-management/pkg/jwt"
	"strings"
)

func getTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if !strings.Contains(authHeader, "Bearer") {
		return ""
	}

	if len(authHeader) > 7 && strings.ToUpper(authHeader[0:6]) == "BEARER" {
		return authHeader[7:]
	}

	return ""
}

func VerifyAuthenticationToken(jwtKey string) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			// validate auth token
			tokenString := getTokenFromHeader(r)
			if tokenString == "" {
				http.Error(w, "token not found", http.StatusBadRequest)
				return
			}

			// validate jwt token
			claims, err := jwt.ValidateToken(tokenString, jwtKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			// pass several data on request context
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "user_id", claims["id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
