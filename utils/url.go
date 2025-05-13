package utils

import (
	"fmt"
	"net/url"
	"slices"
	"strings"

	"github.com/rs/zerolog/log"
)

func ValidateURLScheme(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	validSchemes := []string{"http", "https"}
	if slices.Contains(validSchemes, u.Scheme) {
		return nil
	}

	return fmt.Errorf("link has invalid scheme. Must have schemes %v", validSchemes)
}

func IsURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func CleanHost(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("host is required")
	}

	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}

	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}
	host := u.Hostname()
	log.Debug().
		Str("host", host).
		Msg("Cleaned host")

	return host, nil
}
