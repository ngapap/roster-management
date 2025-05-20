package middlewares

import (
	"github.com/go-chi/cors"
)

var Cors = cors.Handler(cors.Options{
	AllowedOrigins:   []string{"http://localhost:9002", "http://localhost:9001"},
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Origin", "Authorization", "Content-Type"},
	ExposedHeaders:   []string{"Link"},
	AllowCredentials: true,
	MaxAge:           300,
})
