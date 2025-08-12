package server

import (
	"net/http"

	"github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/handlers"
	"github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/models"
	"github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/repositories"
	"github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/services"
	"github.com/go-chi/chi/v5"
)

func Start() error {
	repo := repositories.NewInMemoryProductRepo([]models.Product{
		{ID: "1", Name: "Burger", Price: 9.99, Category: "Fast Food"},
		{ID: "2", Name: "Pizza", Price: 14.50, Category: "Italian"},
	})
	service := services.NewProductService(repo)
	handler := handlers.NewProductHandler(service)

	r := chi.NewRouter()
	r.Get("/products", handler.ListProducts)

	return http.ListenAndServe(":8080", r)
}
