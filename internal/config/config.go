package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port                 string
	RateLimitIP          int
	RateLimitIPWindow    time.Duration
	RateLimitToken       int
	RateLimitTokenWindow time.Duration
	BlockDuration        time.Duration
	RedisAddr            string
	RedisPassword        string
	RedisDB              int
}

func Load() Config {
	return Config{
		Port:                 getEnv("PORT", "8080"),
		RateLimitIP:          getEnvAsInt("RATE_LIMIT_IP", 10),
		RateLimitIPWindow:    getEnvAsDuration("RATE_LIMIT_IP_DURATION", time.Second),
		RateLimitToken:       getEnvAsInt("RATE_LIMIT_TOKEN", 100),
		RateLimitTokenWindow: getEnvAsDuration("RATE_LIMIT_TOKEN_DURATION", time.Second),
		BlockDuration:        getEnvAsDuration("BLOCK_DURATION", 300*time.Second),
		RedisAddr:            getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:        getEnv("REDIS_PASSWORD", ""),
		RedisDB:              getEnvAsInt("REDIS_DB", 0),
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}

func getEnvAsDuration(name string, defaultVal time.Duration) time.Duration {
	if valStr := os.Getenv(name); valStr != "" {
		if val, err := time.ParseDuration(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
