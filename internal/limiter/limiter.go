package limiter

import (
	"context"
	"net/http"
	"strings"

	"github.com/Ryan-18-system/rate-limit-golang/internal/config"
)

type RateLimiter struct {
	Strategy Strategy
	Config   config.Config
}

func NewRateLimiter(strategy Strategy, cfg config.Config) *RateLimiter {
	return &RateLimiter{
		Strategy: strategy,
		Config:   cfg,
	}
}

func (rl *RateLimiter) Allow(ctx context.Context, r *http.Request) (bool, string) {
	token := r.Header.Get("API_KEY")
	if token != "" {
		key := "token:" + token
		allowed, err := rl.Strategy.Allow(ctx, key, rl.Config.RateLimitToken, rl.Config.RateLimitTokenWindow, rl.Config.BlockDuration)
		if err == nil {
			return allowed, key
		}
	}

	ip := extractIP(r)
	key := "ip:" + ip
	allowed, err := rl.Strategy.Allow(ctx, key, rl.Config.RateLimitIP, rl.Config.RateLimitIPWindow, rl.Config.BlockDuration)
	if err != nil {
		return false, key
	}
	return allowed, key
}

func extractIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip != "" {
		return strings.Split(ip, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}
