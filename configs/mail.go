package configs

import (
	"errors"
	"net"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type MailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string

	// TLS enables implicit TLS (SMTPS).
	TLS bool
	// StartTLS upgrades a plain connection using STARTTLS (when supported).
	StartTLS bool
	// Insecure skips TLS cert verification (not recommended for production).
	Insecure bool
}

func LoadMailConfig() (MailConfig, error) {
	_ = godotenv.Load()

	host := strings.TrimSpace(os.Getenv("SMTP_HOST"))
	port := strings.TrimSpace(os.Getenv("SMTP_PORT"))
	user := strings.TrimSpace(os.Getenv("SMTP_USER"))
	pass := os.Getenv("SMTP_PASS")
	from := strings.TrimSpace(os.Getenv("SMTP_FROM"))

	if port == "" {
		port = "587"
	}
	if from == "" {
		from = user
	}

	tlsEnabled := envBool("SMTP_TLS", false)
	startTLS := envBool("SMTP_STARTTLS", !tlsEnabled)
	insecure := envBool("SMTP_INSECURE", false)

	if host == "" {
		return MailConfig{}, errors.New("SMTP_HOST is required")
	}
	if strings.TrimSpace(port) == "" {
		return MailConfig{}, errors.New("SMTP_PORT is required")
	}
	if strings.TrimSpace(from) == "" {
		return MailConfig{}, errors.New("SMTP_FROM (or SMTP_USER) is required")
	}

	return MailConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: pass,
		From:     from,
		TLS:      tlsEnabled,
		StartTLS: startTLS,
		Insecure: insecure,
	}, nil
}

func (c MailConfig) Addr() string {
	return net.JoinHostPort(c.Host, c.Port)
}

func envBool(key string, def bool) bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if v == "" {
		return def
	}
	switch v {
	case "1", "true", "t", "yes", "y", "on":
		return true
	case "0", "false", "f", "no", "n", "off":
		return false
	default:
		return def
	}
}
