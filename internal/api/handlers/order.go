package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/malakagl/kart-challenge/internal/couponcode"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models"
	"github.com/malakagl/kart-challenge/pkg/services"
	"github.com/malakagl/kart-challenge/pkg/util"
)

type OrderHandler struct {
	orderService      services.OrderRepository
	productService    services.ProductRepository
	couponValidator   couponcode.CouponValidator
	couponCodeService services.CouponCodeRepository
}

func NewOrderHandler(o services.OrderRepository, p services.ProductRepository, v couponcode.CouponValidator, c services.CouponCodeRepository) *OrderHandler {
	return &OrderHandler{orderService: o, productService: p, couponValidator: v, couponCodeService: c}
}

func (h *OrderHandler) isCouponCodeValid(code string) bool {
	if len(code) < 8 || len(code) > 10 {
		return false
	}

	if h.couponCodeService.CountFilesByCode(code) > 1 {
		log.Error().Msgf("Coupon code %s is valid: found in multiple files", code)
		return true
	}

	return false
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		log.Error().Msgf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if orderReq.CouponCode == "" || h.couponCodeService.CountFilesByCode(orderReq.CouponCode) > 1 { // !h.couponValidator.ValidateCouponCode(orderReq.CouponCode)
		log.Error().Msgf("Invalid coupon code: %s", orderReq.CouponCode)
		http.Error(w, "Invalid coupon code", http.StatusUnprocessableEntity)
		return
	}

	order := models.Order{}
	orderProducts := make([]*models.OrderProduct, len(orderReq.Items))
	products := make([]*models.Product, len(orderReq.Items))
	for i, item := range orderReq.Items {
		productId, err := util.StringToUint(item.ProductID)
		if err != nil || productId == 0 {
			log.Error().Msg("Invalid product ID in order request")
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		product, err := h.productService.FindByID(productId)
		if err != nil || product == nil {
			log.Error().Msgf("Error fetching product: %v", err)
			http.Error(w, "Failed to fetch product", http.StatusBadRequest)
			return
		}

		orderProducts[i] = &models.OrderProduct{ProductID: item.ProductID, Quantity: item.Quantity}
		products[i] = product
		order.Total += product.Price * float64(item.Quantity)
	}
	order.Products = orderProducts

	orderID, err := h.orderService.Create(order)
	if err != nil {
		log.Error().Msgf("Error creating orderReq: %v", err)
		http.Error(w, "Failed to create orderReq", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(models.OrderResponse{
		ID:       orderID,
		Items:    orderReq.Items,
		Products: products,
	})
	if err != nil {
		log.Error().Msgf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
