package jwt

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
)

var SIGNING_METHOD = jwt.SigningMethodHS256

var (
	ErrInvalidMethod  = errors.New("invalid signing method")
	ErrInvalidSession = errors.New("invalid session")
)

func CreateToken(claim jwt.MapClaims, jwtKey string) (string, error) {
	tkn := jwt.NewWithClaims(SIGNING_METHOD, claim)

	return tkn.SignedString([]byte(jwtKey))
}

func ValidateToken(tokenString string, jwtKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidMethod
		} else if method != SIGNING_METHOD {
			return nil, ErrInvalidMethod
		}

		return jwtKey, nil
	})

	if err != nil {
		return nil, ErrInvalidSession
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidSession
}
