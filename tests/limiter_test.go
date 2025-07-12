package tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Ryan-18-system/rate-limit-golang/internal/config"
	"github.com/Ryan-18-system/rate-limit-golang/internal/limiter"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func flushDB(db int) {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   db,
	})
	_ = rdb.FlushDB(context.Background()).Err()
}
func setupLimiter() *limiter.RateLimiter {
	db := 15
	flushDB(db) // Limpa antes de cada teste

	cfg := config.Config{
		RateLimitIP:          2,
		RateLimitIPWindow:    1 * time.Second,
		RateLimitToken:       3,
		RateLimitTokenWindow: 1 * time.Second,
		BlockDuration:        2 * time.Second,
		RedisAddr:            "localhost:6379",
		RedisPassword:        "",
		RedisDB:              db,
	}

	strategy, err := limiter.NewRedisStrategy(cfg)
	if err != nil {
		panic(err)
	}

	return limiter.NewRateLimiter(strategy, cfg)
}

func TestAllowUnderLimitByIP(t *testing.T) {
	rl := setupLimiter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.100:1234"

	allowed, _ := rl.Allow(context.Background(), req)
	assert.True(t, allowed, "primeira requisição por IP deveria ser permitida")

	allowed, _ = rl.Allow(context.Background(), req)
	assert.True(t, allowed, "segunda requisição por IP deveria ser permitida")
}

func TestBlockAfterLimitByIP(t *testing.T) {
	rl := setupLimiter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.101:1234"

	// Excedendo o limite de 2
	rl.Allow(context.Background(), req)
	rl.Allow(context.Background(), req)
	allowed, _ := rl.Allow(context.Background(), req)

	assert.False(t, allowed, "terceira requisição por IP deve ser bloqueada")
}

func TestAllowUnderLimitByToken(t *testing.T) {
	rl := setupLimiter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "token123")

	allowed, _ := rl.Allow(context.Background(), req)
	assert.True(t, allowed)

	allowed, _ = rl.Allow(context.Background(), req)
	assert.True(t, allowed)

	allowed, _ = rl.Allow(context.Background(), req)
	assert.True(t, allowed)
}

func TestBlockAfterLimitByToken(t *testing.T) {
	rl := setupLimiter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("API_KEY", "token456")

	// Limite de 3
	rl.Allow(context.Background(), req)
	rl.Allow(context.Background(), req)
	rl.Allow(context.Background(), req)
	allowed, _ := rl.Allow(context.Background(), req)

	assert.False(t, allowed, "quarta requisição deve ser bloqueada por token")
}

func TestTokenOverridesIP(t *testing.T) {
	rl := setupLimiter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.200:1234"
	req.Header.Set("API_KEY", "super-token")

	// Limite IP: 2, Token: 3
	rl.Allow(context.Background(), req)
	rl.Allow(context.Background(), req)
	rl.Allow(context.Background(), req)

	// Deve ainda permitir por token (mesmo com IP excedido)
	allowed, _ := rl.Allow(context.Background(), req)
	assert.False(t, allowed, "quarta deve ser bloqueada por token")
}

func TestResetAfterWindow(t *testing.T) {
	rl := setupLimiter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "192.168.1.150:1234"

	rl.Allow(context.Background(), req)
	rl.Allow(context.Background(), req)
	assert.False(t, assertAllowed(rl, req), "deve ser bloqueado")

	// Espera janela e bloqueio expirar
	time.Sleep(3 * time.Second)

	assert.True(t, assertAllowed(rl, req), "deve permitir após o tempo de bloqueio")
}

func assertAllowed(rl *limiter.RateLimiter, req *http.Request) bool {
	allowed, _ := rl.Allow(context.Background(), req)
	return allowed
}
