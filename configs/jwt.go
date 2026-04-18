package configs

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type JWTConfig struct {
	Secret string
	// TTLMin is token validity duration in minutes.
	TTLMin int
}

func LoadJWTConfig() (JWTConfig, error) {
	_ = godotenv.Load()

	secret := strings.TrimSpace(os.Getenv("JWT_SECRET"))
	if secret == "" {
		return JWTConfig{}, errors.New("JWT_SECRET is required")
	}

	ttlMin := envInt("JWT_TTL_MIN", 60)
	if ttlMin <= 0 {
		ttlMin = 60
	}

	return JWTConfig{
		Secret: secret,
		TTLMin: ttlMin,
	}, nil
}

func envInt(key string, def int) int {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
