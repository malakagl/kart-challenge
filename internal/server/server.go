package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/couponcode"
	"github.com/malakagl/kart-challenge/internal/database"
	"github.com/malakagl/kart-challenge/internal/middleware"
	"github.com/malakagl/kart-challenge/internal/routes"
	"github.com/malakagl/kart-challenge/pkg/errors"
	"github.com/malakagl/kart-challenge/pkg/log"
	"gorm.io/gorm"
)

type Server struct {
	httpServer *http.Server
	ErrChan    chan error
	db         *gorm.DB
	cfg        *config.Config
}

func NewServer(c *config.Config) *Server {
	return &Server{
		ErrChan: make(chan error, 1),
		cfg:     c,
	}
}

// Start sets up the database, coupon codes, routes, and starts the HTTP server.
// Returns the server instance for later shutdown.
func (s *Server) Start() error {
	if err := database.RunMigrations(&s.cfg.Database); err != nil {
		log.Error().Err(err).Msgf("database migrations failed.")
		return err
	}

	ctx := context.Background()
	couponcode.SetCouponCodeFiles(s.cfg.CouponCode.FilePaths)
	go func() {
		log.Info().Msg("Started decompressing coupon code files in background")
		err := couponcode.SetupCouponCodeFiles(s.cfg.CouponCode.FilePaths)
		if err != nil {
			log.Error().Err(err).Msg("coupon code setup failed.")
		}
	}()

	couponcode.InitCache(s.cfg.Server.MaxCouponCodeCacheSize)
	log.Info().Msgf("connecting to database")
	var err error
	s.db, err = database.Connect(ctx, &s.cfg.Database)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to database.")
		return err
	}

	log.Info().Msgf("creating routes")
	middleware.SetRateLimits(s.cfg.Server.ReqLimitPerIPPerSec, s.cfg.Server.ReqBurstPerIPPerSec)
	r := chi.NewRouter()
	r.Use(middleware.Trace, middleware.Logging, middleware.Authentication, middleware.RateLimit)
	routes.AddHealthCheckRoutes(r)
	routes.AddProductRoutes(r, s.db)
	routes.AddOrderRoutes(r, s.db)

	serverAddr := fmt.Sprintf("%s:%d", s.cfg.Server.Host, s.cfg.Server.Port)
	s.httpServer = &http.Server{
		Addr:    serverAddr,
		Handler: r,
	}

	go func() {
		if errSrv := s.httpServer.ListenAndServe(); errSrv != nil && !errors.Is(errSrv, http.ErrServerClosed) {
			log.Error().Err(errSrv).Msgf("HTTP server failed.")
			s.ErrChan <- errSrv
		}
	}()

	log.Info().Msgf("Server started on %s", serverAddr)
	return nil
}

// Stop gracefully shuts down the server with a context timeout.
func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		log.Info().Msg("Shutting down server gracefully")
		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Server shutdown error.")
		}
	}

	if s.db != nil {
		db, err := s.db.DB()
		if err == nil {
			if err = db.Close(); err != nil {
				log.Error().Err(err).Msg("error while closing sql.DB")
			} else {
				log.Info().Msg("DB closed successfully")
			}
		} else {
			log.Error().Err(err).Msg("failed to retrieve sql.DB for closing")
		}
	}

	log.Info().Msg("Server stopped successfully")
	return nil
}
