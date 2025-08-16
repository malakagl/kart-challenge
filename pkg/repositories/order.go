package repositories

import (
	"sync"

	"github.com/google/uuid"
	"github.com/malakagl/kart-challenge/pkg/models"
)

var mu = &sync.RWMutex{}

type InMemoryOrderRepo struct {
	orders []models.Order
}

func NewInMemoryOrderRepo() *InMemoryOrderRepo {
	return &InMemoryOrderRepo{}
}

func (r *InMemoryOrderRepo) Create(order models.Order) (string, error) {
	id := uuid.New().String()
	order.ID = id
	mu.Lock()
	defer mu.Unlock()
	r.orders = append(r.orders, order)
	return id, nil
}
