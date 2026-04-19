package configs

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port string

	// AccessTokenPeriod is access-token validity in minutes.
	AccessTokenPeriod int
	// RefreshTokenPeriod is refresh-token validity in minutes.
	RefreshTokenPeriod int
}

func LoadAppConfig() AppConfig {
	_ = godotenv.Load()

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}

	accessMin := envInt("ACCESS_TOKEN_PERIOD_MIN", 60)
	if accessMin <= 0 {
		accessMin = 60
	}
	refreshMin := envInt("REFRESH_TOKEN_PERIOD_MIN", 1440)
	if refreshMin <= 0 {
		refreshMin = 1440
	}

	return AppConfig{
		Port:               port,
		AccessTokenPeriod:  accessMin,
		RefreshTokenPeriod: refreshMin,
	}
}
