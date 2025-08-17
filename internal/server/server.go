package server

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	handlers2 "github.com/malakagl/kart-challenge/internal/api/handlers"
	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/couponcode"
	"github.com/malakagl/kart-challenge/internal/database"
	"github.com/malakagl/kart-challenge/internal/middleware"
	"github.com/malakagl/kart-challenge/pkg/log"
	repositories2 "github.com/malakagl/kart-challenge/pkg/repositories"
	services2 "github.com/malakagl/kart-challenge/pkg/services"
)

func Start(cfg *config.Config) error {
	if err := database.RunMigrations(cfg.Database); err != nil {
		log.Error().Msgf("database migrations failed: %v", err)
		return err
	}

	if !cfg.CouponCode.Unzipped {
		if err := couponcode.SetupCouponCodeFiles(cfg.CouponCode.FilePaths); err != nil {
			log.Error().Msgf("failed to load coupon codes: %v", err)
			return err
		}
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Error().Msgf("failed to connect to database: %v", err)
		return err
	}

	r := chi.NewRouter()

	r.Use(middleware.AuthenticationMiddleware, middleware.LoggingMiddleware, middleware.ResponseHeadersMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	productRepo := repositories2.NewProductRepo(db)
	productService := services2.NewProductService(productRepo)
	productHandler := handlers2.NewProductHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{productID}", productHandler.GetProductByID)

	v := couponcode.NewValidator(cfg.CouponCode.FilePaths)
	orderRepo := repositories2.NewOrderRepo(db)
	orderService := services2.NewOrderService(orderRepo)
	couponCodeRepo := repositories2.NewCouponCodeRepository(db)
	couponCodeService := services2.NewCouponCodeService(couponCodeRepo)
	orderHandler := handlers2.NewOrderHandler(orderService, productService, v, couponCodeService)
	r.Post("/orders", orderHandler.CreateOrder)

	serverURL := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info().Msgf("Server starting on %s", serverURL)
	return http.ListenAndServe(serverURL, r)
}
