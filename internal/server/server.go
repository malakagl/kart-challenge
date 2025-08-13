package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/internal/handlers"
	"github.com/malakagl/kart-challenge/pkg/models"
	repositories2 "github.com/malakagl/kart-challenge/pkg/repositories"
	services2 "github.com/malakagl/kart-challenge/pkg/services"
)

func Start() error {
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
	productHandler := handlers.NewProductHandler(productService)
	r.Get("/products", productHandler.ListProducts)
	r.Get("/products/{productID}", productHandler.GetProductByID)

	orderRepo := repositories2.NewInMemoryOrderRepo()
	orderService := services2.NewOrderService(orderRepo)
	orderHandler := handlers.NewOrderHandler(orderService, productService)
	r.Post("/orders", orderHandler.CreateOrder)

	log.Println("Server starting on port 8080")
	return http.ListenAndServe(":8080", r)
}
