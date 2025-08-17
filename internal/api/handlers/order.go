package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/dto/request"
	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
	"github.com/malakagl/kart-challenge/pkg/services"
)

type OrderHandler struct {
	orderService services.OrderService
}

func NewOrderHandler(o services.OrderService) *OrderHandler {
	return &OrderHandler{orderService: o}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var orderReq request.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		log.Error().Msgf("Error decoding request body: %v", err)
		response.Error(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := validator.New().Struct(orderReq); err != nil {
		log.Error().Msgf("Validation error: %v", err)
		response.Error(w, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	orderRes, err := h.orderService.Create(&orderReq)
	if err != nil {
		log.Error().Msgf("Error creating order: %v", err)
		if errors.Is(err, constants.ErrInvalidCouponCode) {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid coupon code", err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to create order", err.Error())
		}
		return
	}

	response.Success(w, orderRes)
}
