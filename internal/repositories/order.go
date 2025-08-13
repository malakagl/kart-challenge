package repositories

import (
	"github.com/google/uuid"
	"github.com/malakagl/kart-challenge/internal/models"
)

type InMemoryOrderRepo struct {
	orders []models.Order
}

func NewInMemoryOrderRepo() *InMemoryOrderRepo {
	return &InMemoryOrderRepo{}
}

func (r *InMemoryOrderRepo) Create(order models.Order) (string, error) {
	id := uuid.New().String()
	order.ID = id
	r.orders = append(r.orders, order)
	return id, nil
}
