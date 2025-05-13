package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"dynamic-links-generator/api/apperrors"

	"github.com/rs/zerolog/log"
)

type LinkRepository interface {
	GetQueryParamsByHostAndPath(ctx context.Context, host, path string) (string, error)
	FindExistingShortLink(ctx context.Context, host, rawQS string) (string, error)
	CreateShortLink(ctx context.Context, host, path, rawQS string, unguessable bool) error
}

type linkRepository struct {
	db *sql.DB
}

func NewLinkRepository(db *sql.DB) LinkRepository {
	return &linkRepository{
		db: db,
	}
}

func (r *linkRepository) GetQueryParamsByHostAndPath(ctx context.Context, host, path string) (string, error) {
	var rawQueryStr string

	row := r.db.QueryRowContext(
		ctx,
		`SELECT query_params
           FROM dynamic_links
          WHERE host = $1 AND path = $2`,
		host,
		path,
	)
	if err := row.Scan(&rawQueryStr); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Debug().
				Str("path", path).
				Msg("Link not found in database")
			return "", apperrors.ErrLinkNotFound
		}
		log.Error().
			Err(err).
			Str("path", path).
			Msg("Failed to retrieve link from database")
		return "", fmt.Errorf("database error: %w", err)
	}

	return rawQueryStr, nil
}

func (r *linkRepository) FindExistingShortLink(ctx context.Context, host, rawQS string) (string, error) {
	var path string
	const q = `
    SELECT path
      FROM dynamic_links
     WHERE host                = $1
       AND query_params        = $2
       AND is_unguessable_path = FALSE
     LIMIT 1`
	err := r.db.QueryRowContext(ctx, q, host, rawQS).Scan(&path)
	return path, err
}

func (r *linkRepository) CreateShortLink(ctx context.Context, host, path, rawQS string, unguessable bool) error {
	const stmt = `
    INSERT INTO dynamic_links
      (host, path, query_params, is_unguessable_path)
    VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(
		ctx,
		stmt,
		host,
		path,
		rawQS,
		unguessable,
	)
	return err
}
