package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/internal/services"
	"github.com/malakagl/kart-challenge/pkg/constants"
)

type ProductHandler struct {
	service services.ProductRepository
}

func NewProductHandler(s services.ProductRepository) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, _ *http.Request) {
	products, err := h.service.FindAll()
	if err != nil {
		http.Error(w, "failed to fetch products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	productID := chi.URLParam(r, "productID")
	if productID == "" {
		log.Default().Println("Invalid product ID")
		http.Error(w, "Invalid ID supplied", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindByID(productID)
	if err != nil {
		log.Default().Println("Error fetching product:", err)
		if errors.Is(err, constants.ErrProductNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch product", http.StatusInternalServerError)
		return
	}

	if product == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
