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

	url := strings.TrimSpace(os.Getenv("FONNTE_URL"))
	if url == "" {
		url = "https://api.blu.com/send"
	}

	token := strings.TrimSpace(os.Getenv("FONNTE_TOKEN"))

	if token == "" {
		return WhatsAppConfig{}, errors.New("FONNTE_TOKEN is required")
	}

	return WhatsAppConfig{
		URL:   url,
		Token: token,
	}, nil
}
