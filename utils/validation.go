package utils

import (
	"net/url"
	"strings"

	"github.com/rs/zerolog/log"
)

func IsNumericString(s string) bool {
	if s == "" {
		return true
	}

	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

func IsDomainAllowed(allowList []string, rawLink string) bool {
	u, err := url.Parse(rawLink)
	if err != nil {
		log.Error().
			Str("raw_link", rawLink).
			Msg("Invalid link")
		return false
	}
	host := strings.ToLower(u.Hostname())

	for _, allowed := range allowList {
		allowed = strings.ToLower(strings.TrimSpace(allowed))
		if host == allowed {
			return true
		}
	}
	return false
}
