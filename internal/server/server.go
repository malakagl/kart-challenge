package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/database"
	"github.com/malakagl/kart-challenge/internal/middleware"
	"github.com/malakagl/kart-challenge/internal/routes"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
)

func Start(cfg *config.Config) error {
	if err := database.RunMigrations(cfg.Database); err != nil {
		log.Error().Msgf("database migrations failed: %v", err)
		return err
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Error().Msgf("failed to connect to database: %v", err)
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.AuthenticationMiddleware, middleware.LoggingMiddleware)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Success(w, "ok")
	})

	routes.AddProductRoutes(r, db)
	routes.AddOrderRoutes(r, db)

	serverURL := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info().Msgf("Server starting on %s", serverURL)
	return http.ListenAndServe(serverURL, r)
}
