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
	"roster-management/pkg/middlewares"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"roster-management/cmd/service-api/handler"
	"roster-management/cmd/service-api/repository/postgres"
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

type customMux struct {
	http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func (cm *customMux) Use(middleware func(http.Handler) http.Handler) {
	cm.middlewares = append(cm.middlewares, middleware)
}

func (cm *customMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var current http.Handler = &cm.ServeMux

	for _, next := range cm.middlewares {
		current = next(current)
	}

	current.ServeHTTP(w, r)
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

	hand := handler.NewHandlers(repo, config.GetString("JWT_KEY"), config.GetInt("JWT_EXPIRY"))

	mux := new(customMux)

	//set middlewares
	mux.Use(middlewares.VerifyAuthenticationToken(config.GetString("JWT_KEY")))

	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})
	mux.HandleFunc("POST /api/auth/register", hand.RegisterWorker)
	mux.HandleFunc("POST /api/auth/login", hand.Login)

	_, port, err := net.SplitHostPort(config.GetString("SERVICE_API_HOST"))
	if err != nil {
		panic(err)
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
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
