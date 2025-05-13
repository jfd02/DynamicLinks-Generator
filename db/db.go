package db

import (
	"database/sql"
	"fmt"

	"dynamic-links-generator/config"

	"github.com/rs/zerolog/log"
)

type DB struct {
	*sql.DB
}

func New(cfg *config.Config) (*DB, error) {
	db, err := sql.Open(cfg.DBDriver, cfg.DBConnectionStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Info().Msg("Successfully connected to database")
	return &DB{DB: db}, nil
}
