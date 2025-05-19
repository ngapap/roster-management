package handler

import "roster-management/cmd/service-api/repository"

type JWTConfig struct {
	Key string
	Exp int
}

type Handler struct {
	repo   repository.Repository
	jwtCfg *JWTConfig
}

func NewHandlers(repo repository.Repository, jwtKey string, jwtExp int) *Handler {
	return &Handler{
		repo: repo,
		jwtCfg: &JWTConfig{
			Key: jwtKey,
			Exp: jwtExp,
		},
	}
}
