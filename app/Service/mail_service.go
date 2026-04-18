package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"

	"govibe/configs"
)

func SendEmail(receiver, title, body string) error {
	receiver = strings.TrimSpace(receiver)
	title = strings.TrimSpace(title)

	if receiver == "" {
		return errors.New("receiver is required")
	}
	if title == "" {
		return errors.New("title is required")
	}

	cfg, err := configs.LoadMailConfig()
	if err != nil {
		return err
	}

	addr := cfg.Addr()
	tlsCfg := &tls.Config{
		ServerName:         cfg.Host,
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: cfg.Insecure,
	}

	var c *smtp.Client
	if cfg.TLS {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, tlsCfg)
		if err != nil {
			return fmt.Errorf("smtp tls dial failed: %w", err)
		}
		client, err := smtp.NewClient(conn, cfg.Host)
		if err != nil {
			_ = conn.Close()
			return fmt.Errorf("smtp client init failed: %w", err)
		}
		c = client
	} else {
		client, err := smtp.Dial(addr)
		if err != nil {
			return fmt.Errorf("smtp dial failed: %w", err)
		}
		c = client
	}
	defer func() { _ = c.Close() }()

	if !cfg.TLS && cfg.StartTLS {
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsCfg); err != nil {
				return fmt.Errorf("smtp starttls failed: %w", err)
			}
		}
	}

	if cfg.User != "" && cfg.Password != "" {
		if ok, _ := c.Extension("AUTH"); ok {
			auth := smtp.PlainAuth("", cfg.User, cfg.Password, cfg.Host)
			if err := c.Auth(auth); err != nil {
				return fmt.Errorf("smtp auth failed: %w", err)
			}
		}
	}

	if err := c.Mail(cfg.From); err != nil {
		return fmt.Errorf("smtp mail-from failed: %w", err)
	}
	if err := c.Rcpt(receiver); err != nil {
		return fmt.Errorf("smtp rcpt-to failed: %w", err)
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("smtp data failed: %w", err)
	}

	subject := sanitizeHeader(title)
	msg := strings.Join([]string{
		fmt.Sprintf("From: %s", cfg.From),
		fmt.Sprintf("To: %s", receiver),
		fmt.Sprintf("Subject: %s", subject),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
		"",
	}, "\r\n")

	if _, err := w.Write([]byte(msg)); err != nil {
		_ = w.Close()
		return fmt.Errorf("smtp write failed: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp data close failed: %w", err)
	}

	if err := c.Quit(); err != nil {
		return fmt.Errorf("smtp quit failed: %w", err)
	}

	return nil
}

func sanitizeHeader(v string) string {
	v = strings.ReplaceAll(v, "\r", " ")
	v = strings.ReplaceAll(v, "\n", " ")
	return strings.TrimSpace(v)
}
