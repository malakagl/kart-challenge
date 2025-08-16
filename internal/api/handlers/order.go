package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/malakagl/kart-challenge/internal/couponcode"
	"github.com/malakagl/kart-challenge/pkg/models"
	"github.com/malakagl/kart-challenge/pkg/services"
)

type OrderHandler struct {
	orderService    services.OrderRepository
	productService  services.ProductRepository
	couponValidator couponcode.CouponValidator
}

func NewOrderHandler(o services.OrderRepository, p services.ProductRepository, v couponcode.CouponValidator) *OrderHandler {
	return &OrderHandler{orderService: o, productService: p, couponValidator: v}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if orderReq.CouponCode == "" || !h.couponValidator.ValidateCouponCode(orderReq.CouponCode) {
		log.Println("Invalid coupon code:", orderReq.CouponCode)
		http.Error(w, "Invalid coupon code", http.StatusUnprocessableEntity)
		return
	}

	order := models.Order{
		Items: orderReq.Items,
	}
	products := make([]*models.Product, len(orderReq.Items))
	for i, item := range orderReq.Items {
		if item.ProductID == "" {
			log.Println("Invalid product ID in order request")
			http.Error(w, "Invalid product ID", http.StatusBadRequest)
			return
		}

		product, err := h.productService.FindByID(item.ProductID)
		if err != nil || product == nil {
			log.Println("Error fetching product:", err)
			http.Error(w, "Failed to fetch product", http.StatusBadRequest)
			return
		}

		products[i] = product
		order.Total += product.Price * float64(item.Quantity)
	}

	orderID, err := h.orderService.Create(order)
	if err != nil {
		log.Println("Error creating orderReq:", err)
		http.Error(w, "Failed to create orderReq", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(models.Order{
		ID:       orderID,
		Items:    orderReq.Items,
		Products: products,
	})
	if err != nil {
		log.Println("Error encoding response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
