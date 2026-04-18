package configs

import (
	"errors"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type WhatsAppConfig struct {
	URL   string
	Token string
}

func LoadWhatsAppConfig() (WhatsAppConfig, error) {
	_ = godotenv.Load()

	url := strings.TrimSpace(os.Getenv("BLU_URL"))
	if url == "" {
		url = "https://api.blu.com/send"
	}

	// Matches the Laravel sample env name "blu_TOKEN", but normalized to uppercase.
	token := strings.TrimSpace(os.Getenv("BLU_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("blu_TOKEN"))
	}

	if token == "" {
		return WhatsAppConfig{}, errors.New("BLU_TOKEN is required")
	}

	return WhatsAppConfig{
		URL:   url,
		Token: token,
	}, nil
}
