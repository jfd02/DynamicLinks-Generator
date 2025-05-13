package utils

import (
	"crypto/rand"

	"github.com/rs/zerolog/log"
)

func GenerateDynamicLinkPath(length int) string {
	const alphanumeric = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Panic().Err(err).Msg("Failed to generate random bytes")
	}

	for i := range b {
		b[i] = alphanumeric[b[i]%byte(len(alphanumeric))]
	}

	id := string(b)

	log.Debug().
		Str("short_code", id).
		Msg("Generated alphanumeric short ID")

	return id
}
