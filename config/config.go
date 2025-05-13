package config

import (
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

type Config struct {
	Port                  string
	DBDriver              string
	DBConnectionStr       string
	ShortPathLength       int
	UnguessablePathLength int
	URLScheme             string
	DomainAllowList       []string
	LogLevel              string
}

func New() *Config {
	return &Config{
		Port:                  getEnv("PORT", "9010"),
		DBDriver:              getEnv("DB_DRIVER", "postgres"),
		DBConnectionStr:       getEnv("DATABASE_URL", ""),
		ShortPathLength:       getEnvAsInt("SHORT_PATH_LENGTH", 6),
		UnguessablePathLength: getEnvAsInt("UNGUESSABLE_PATH_LENGTH", 10),
		URLScheme:             getEnv("URL_SCHEME", "https"),
		DomainAllowList:       getEnvAsSlice("DOMAIN_ALLOW_LIST", []string{}),
		LogLevel:              getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsSlice(key string, defaultVal []string) []string {
	if value, exists := os.LookupEnv(key); exists {
		return strings.Split(value, ",")
	}
	return defaultVal
}

func getEnvAsInt(name string, defaultVal int) int {
	if valStr, ok := os.LookupEnv(name); ok {
		if val, err := strconv.Atoi(valStr); err == nil {
			return val
		}
	}
	return defaultVal
}
