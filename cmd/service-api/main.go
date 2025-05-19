package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"

	"roster-management/cmd/service-api/handler"
	"roster-management/cmd/service-api/repository/postgres"
	"roster-management/pkg/middlewares"
)

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	configPath := flag.String("config", "configs/local.env", "config file")
	config := viper.New()
	config.AutomaticEnv()
	config.SetConfigFile(*configPath)
	if err := config.ReadInConfig(); err != nil {
		panic(err)
	}
	repo, err := postgres.NewRepositoryFromConfig(config)
	if err != nil {
		panic(err)
	}
	jwtKey := config.GetString("JWT_KEY")
	hand := handler.NewHandlers(repo, jwtKey, config.GetInt("JWT_EXPIRY"))

	r := chi.NewRouter()

	r.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Route("/api", func(r chi.Router) {

		// auth
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", hand.RegisterWorker)
			r.Post("/login", hand.Login)
		})

		// auth verification
		r.Group(func(r chi.Router) {
			r.Use(middlewares.VerifyAuthenticationToken(jwtKey))

			//	shift
			r.Route("/shift", func(r chi.Router) {
				r.Post("/", hand.CreateShift)
				r.Put("/{shiftID}", hand.UpdateShift)
				r.Delete("/{shiftID}", hand.DeleteShift)
				r.Get("/available", hand.GetAvailableShifts)

				// worker
				r.Get("/worker/{workerID}", hand.GetShiftByWorker)

			})

			// shift-request
			r.Route("/shift-request", func(r chi.Router) {
				r.Post("/", hand.CreateShiftRequest)
				r.Put("/{requestID}", hand.UpdateShiftRequest)
				r.Delete("/{requestID}", hand.DeleteShiftRequest)

				// worker
				r.Get("/worker/{workerID}", hand.GetShiftRequestByWorker)
			})
		})
	})

	_, port, err := net.SplitHostPort(config.GetString("SERVICE_API_HOST"))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		log.Println("Starting server on :" + port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	// Listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Block until signal is received
	log.Println("Shutting down...")

	// Create context with timeout to force shutdown if it takes too long
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}
