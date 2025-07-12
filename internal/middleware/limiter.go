package middleware

import (
	"context"
	"net/http"

	"github.com/Ryan-18-system/rate-limit-golang/internal/limiter"
)

func LimitMiddleware(rl *limiter.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, _ := rl.Allow(context.Background(), r)
			if !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
