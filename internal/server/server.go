package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	handlers2 "github.com/malakagl/kart-challenge/internal/api/handlers"
	"github.com/malakagl/kart-challenge/internal/config"
	"github.com/malakagl/kart-challenge/internal/couponcode"
	"github.com/malakagl/kart-challenge/internal/db"
	"github.com/malakagl/kart-challenge/pkg/models"
	repositories2 "github.com/malakagl/kart-challenge/pkg/repositories"
	services2 "github.com/malakagl/kart-challenge/pkg/services"
)

func Start(cfg *config.Config) error {
	if err := db.RunMigrations(); err != nil {
		log.Fatalf("db migrations failed: %v", err)
	}

	if !cfg.CouponCodeConfig.Unzipped {
		if err := couponcode.SetupCouponCodeFiles(cfg.CouponCodeConfig.FilePaths); err != nil {
			log.Fatalf("failed to load coupon codes: %v", err)
		}
	}

	//couponcode.LoadCouponCodes()

	r := chi.NewRouter()

	r.Use(AuthenticationMiddleware, LoggingMiddleware, ResponseHeadersMiddleware)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	})

	productRepo := repositories2.NewInMemoryProductRepo([]models.Product{
		{ID: "1", Name: "Burger", Price: 9.99, Category: "Fast Food"},
		{ID: "2", Name: "Pizza", Price: 14.50, Category: "Italian"},
	})
	productService := services2.NewProductService(productRepo)
	productHandler := handlers2.NewProductHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{productID}", productHandler.GetProductByID)

	v := couponcode.NewValidator(cfg.CouponCodeConfig.FilePaths)
	orderRepo := repositories2.NewInMemoryOrderRepo()
	orderService := services2.NewOrderService(orderRepo)
	orderHandler := handlers2.NewOrderHandler(orderService, productService, v)
	r.Post("/orders", orderHandler.CreateOrder)

	serverURL := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Println("Server starting on ", serverURL)
	return http.ListenAndServe(serverURL, r)
}
