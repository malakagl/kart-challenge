package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/couponcode"
	"github.com/malakagl/kart-challenge/internal/database"
	"github.com/malakagl/kart-challenge/internal/middleware"
	"github.com/malakagl/kart-challenge/internal/routes"
	"github.com/malakagl/kart-challenge/pkg/errors"
	"github.com/malakagl/kart-challenge/pkg/log"
)

var httpServer *http.Server

// Start sets up the database, coupon codes, routes, and starts the HTTP server.
// Returns the server instance for later shutdown.
func Start(cfg *config.Config) error {
	if err := database.RunMigrations(cfg.Database); err != nil {
		log.Error().Msgf("database migrations failed: %v", err)
		return err
	}

	couponcode.SetCouponCodeFiles(cfg.CouponCode.FilePaths)
	go func() {
		log.Info().Msg("Started decompressing coupon code files in background")
		err := couponcode.SetupCouponCodeFiles(cfg.CouponCode.FilePaths)
		if err != nil {
			log.Error().Msgf("couponcode setup failed: %v", err)
		}
	}()

	couponcode.InitCache(cfg.Server.MaxCouponCodeCacheSize)
	log.Info().Msgf("connecting to database")
	db, err := database.Connect(context.Background(), &cfg.Database)
	if err != nil {
		log.Error().Msgf("failed to connect to database: %v", err)
		return err
	}

	log.Info().Msgf("creating routes")
	r := chi.NewRouter()
	r.Use(middleware.TraceMiddleware, middleware.AuthenticationMiddleware, middleware.LoggingMiddleware)
	routes.AddHealthCheckRoutes(r)
	routes.AddProductRoutes(r, db)
	routes.AddOrderRoutes(r, db)

	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	httpServer = &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	// Start server in a goroutine to make it stoppable
	go func() {
		if errSrv := httpServer.ListenAndServe(); errSrv != nil && !errors.Is(errSrv, http.ErrServerClosed) {
			log.Error().Msgf("HTTP server failed: %v", errSrv)
		}
	}()

	log.Info().Msgf("Server started on %s", serverAddr)
	return nil
}

// Stop gracefully shuts down the server with a context timeout.
func Stop(ctx context.Context) error {
	if httpServer == nil {
		return nil
	}

	log.Info().Msg("Shutting down server gracefully")
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Error().Msgf("Server shutdown error: %v", err)
		return err
	}

	log.Info().Msg("Server stopped successfully")
	return nil
}
