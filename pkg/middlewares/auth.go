package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"roster-management/pkg/jwt"
	"roster-management/pkg/util"
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
			logrus.Println("tokenString", tokenString)
			if tokenString == "" {
				logrus.Error(errors.New("token not found"))
				util.SendResponse(w, http.StatusBadRequest, nil, "token not found")
				return
			}

			// validate jwt token
			claims, err := jwt.ValidateToken(tokenString, jwtKey)
			if err != nil {
				logrus.Error(err)
				util.SendResponse(w, http.StatusUnauthorized, nil, err.Error())
				return
			}

			// pass several data on request context
			ctx = context.WithValue(ctx, "email", claims["email"])
			ctx = context.WithValue(ctx, "user_id", claims["id"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
