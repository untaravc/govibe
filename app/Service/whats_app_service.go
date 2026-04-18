package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"govibe/configs"
)

func SendWhatsAppMessage(number, message string) error {
	number = strings.TrimSpace(number)
	message = strings.TrimSpace(message)

	if number == "" {
		return errors.New("number is required")
	}
	if message == "" {
		return errors.New("message is required")
	}

	cfg, err := configs.LoadWhatsAppConfig()
	if err != nil {
		return err
	}

	form := url.Values{}
	form.Set("target", number)
	form.Set("message", message)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.URL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", cfg.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		msg := strings.TrimSpace(string(body))
		if msg == "" {
			msg = http.StatusText(resp.StatusCode)
		}
		return fmt.Errorf("whatsapp request failed (%d): %s", resp.StatusCode, msg)
	}

	return nil
}
