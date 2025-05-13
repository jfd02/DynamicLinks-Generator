package api

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"dynamic-links-generator/api/repository"
	"dynamic-links-generator/api/service"
	"dynamic-links-generator/config"
)

func NewRouter(database *sql.DB, cfg *config.Config) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	linkRepository := repository.NewLinkRepository(database)
	linkService := service.NewLinkService(linkRepository, cfg)
	handler := NewHandler(linkService)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/shortLinks", handler.CreateLink)
		r.Post("/exchangeShortLink", handler.ExchangeShortLink)
	})

	return r
}
