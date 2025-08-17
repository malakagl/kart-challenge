package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/services"
	"github.com/malakagl/kart-challenge/pkg/util"
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
	_ = json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	pID := chi.URLParam(r, "productID")
	productId, err := util.StringToUint(pID)
	if err != nil || productId == 0 {
		log.Error().Msgf("Invalid product ID in order request: %s", pID)
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.service.FindByID(productId)
	if err != nil {
		log.Error().Msgf("Error fetching product: %v", err)
		if errors.Is(err, constants.ErrProductNotFound) {
			http.Error(w, "Product not found", http.StatusNotFound)
			return
		}

		http.Error(w, "failed to fetch product", http.StatusInternalServerError)
		return
	}

	if product == nil {
		log.Error().Msgf("Product not found for ID: %d", productId)
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(product)
}
