package configs

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port string
}

func LoadAppConfig() AppConfig {
	_ = godotenv.Load()

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3000"
	}

	return AppConfig{Port: port}
}

