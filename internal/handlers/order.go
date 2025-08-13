package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/malakagl/kart-challenge/internal/models"
	"github.com/malakagl/kart-challenge/internal/services"
)

type OrderHandler struct {
	orderService   services.OrderRepository
	productService services.ProductRepository
}

func NewOrderHandler(o services.OrderRepository, p services.ProductRepository) *OrderHandler {
	return &OrderHandler{orderService: o, productService: p}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq models.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Println("Received order request:", orderReq)
	order := models.Order{
		Items:      orderReq.Items,
		CouponCode: orderReq.CouponCode,
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
	}

	orderID, err := h.orderService.Create(order)
	if err != nil {
		log.Println("Error creating orderReq:", err)
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
		log.Println("Error encoding response:", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
