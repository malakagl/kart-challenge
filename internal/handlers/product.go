package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/services"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(s *services.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, _ *http.Request) {
	products, err := h.service.GetAllProducts()
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
