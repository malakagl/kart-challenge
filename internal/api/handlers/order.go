package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	errors2 "github.com/malakagl/kart-challenge/pkg/errors"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/dto/request"
	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
	"github.com/malakagl/kart-challenge/pkg/services"
)

type OrderHandler struct {
	orderService services.IOrderService
	validator    *validator.Validate
}

func NewOrderHandler(o services.IOrderService) *OrderHandler {
	return &OrderHandler{
		orderService: o,
		validator:    validator.New(),
	}
}

func (o *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var orderReq request.OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&orderReq); err != nil {
		log.WithCtx(ctx).Error().Msgf("Error decoding request body: %v", err)
		response.Error(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	if err := o.validator.Struct(orderReq); err != nil {
		log.WithCtx(ctx).Error().Msgf("Validation error: %v", err)
		response.Error(w, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	orderRes, err := o.orderService.Create(ctx, &orderReq)
	if err != nil {
		log.WithCtx(ctx).Error().Msgf("Error creating order: %v", err)
		if errors.Is(err, errors2.ErrInvalidCouponCode) {
			response.Error(w, http.StatusUnprocessableEntity, "Invalid coupon code", err.Error())
		} else {
			response.Error(w, http.StatusInternalServerError, "Failed to create order", err.Error())
		}
		return
	}

	response.Success(w, orderRes)
}
