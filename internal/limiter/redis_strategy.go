package limiter

import (
	"context"
	"fmt"
	"time"

	"github.com/Ryan-18-system/rate-limit-golang/internal/config"
	"github.com/go-redis/redis/v8"
)

type Strategy interface {
	Allow(ctx context.Context, key string, limit int, window, blockTime time.Duration) (bool, error)
}

type RedisStrategy struct {
	client *redis.Client
}

func NewRedisStrategy(cfg config.Config) (*RedisStrategy, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("could not connect to Redis: %w", err)
	}

	return &RedisStrategy{client: rdb}, nil
}

func (r *RedisStrategy) Allow(ctx context.Context, key string, limit int, window, blockTime time.Duration) (bool, error) {
	blockedKey := "block:" + key

	// Verifica se está bloqueado
	blocked, err := r.client.Get(ctx, blockedKey).Result()
	if err == nil && blocked == "1" {
		return false, nil
	}

	// Incrementa contador
	count, err := r.client.Incr(ctx, key).Result()
	if err != nil {
		return false, err
	}

	// Se for a primeira vez, define o tempo de expiração
	if count == 1 {
		err = r.client.Expire(ctx, key, window).Err()
		if err != nil {
			return false, err
		}
	}

	if int(count) > limit {
		_ = r.client.Set(ctx, blockedKey, "1", blockTime).Err()
		return false, nil
	}

	return true, nil
}
