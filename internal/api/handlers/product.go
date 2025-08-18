package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	errors2 "github.com/malakagl/kart-challenge/pkg/errors"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
	"github.com/malakagl/kart-challenge/pkg/services"
	"github.com/malakagl/kart-challenge/pkg/util"
)

type ProductHandler struct {
	service services.IProductService
}

func NewProductHandler(s services.IProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	products, err := h.service.FindAll(ctx)
	if err != nil {
		log.WithCtx(ctx).Error().Msgf("Error fetching products: %v", err)
		if errors.Is(err, errors2.ErrProductNotFound) {
			response.Error(w, http.StatusNotFound, "No products found", "No products available in the database")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Error fetching products", "Error fetching products")
		return
	}

	response.Success(w, products)
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	pID := chi.URLParam(r, "productID")
	productId, err := util.StringToUint(pID)
	if err != nil || productId == 0 {
		log.WithCtx(ctx).Error().Msgf("Invalid product ID in order request: %s", pID)
		response.Error(w, http.StatusBadRequest, "Invalid product ID", "Invalid product ID")
		return
	}

	product, err := h.service.FindByID(ctx, productId)
	if err != nil {
		log.WithCtx(ctx).Error().Msgf("Error fetching product: %v", err)
		if errors.Is(err, errors2.ErrProductNotFound) {
			response.Error(w, http.StatusNotFound, "No products found", "No products found in the database")
			return
		}

		response.Error(w, http.StatusInternalServerError, "Error fetching products", "Error fetching products")
		return
	}

	response.Success(w, product)
}
