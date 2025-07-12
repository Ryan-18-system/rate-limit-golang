package main

import (
	"log"
	"net/http"

	"github.com/Ryan-18-system/rate-limit-golang/internal/config"
	"github.com/Ryan-18-system/rate-limit-golang/internal/limiter"
	"github.com/Ryan-18-system/rate-limit-golang/internal/middleware"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	strategy, err := limiter.NewRedisStrategy(cfg)
	if err != nil {
		log.Fatalf("failed to initialize Redis strategy: %v", err)
	}

	rl := limiter.NewRateLimiter(strategy, cfg)
	r := mux.NewRouter()

	r.Use(middleware.LimitMiddleware(rl))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Request OK"))
	})

	log.Printf("Starting server on port %s...", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
