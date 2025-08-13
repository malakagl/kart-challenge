package services

import (
	"log"

	"github.com/malakagl/kart-challenge/pkg/models"
)

type OrderRepository interface {
	Create(order models.Order) (string, error)
}

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(r OrderRepository) *OrderService {
	return &OrderService{repo: r}
}

func (s *OrderService) Create(order models.Order) (string, error) {
	orderID, err := s.repo.Create(order)
	if err != nil {
		log.Println("Error creating order:", err)
		return "", err
	}

	return orderID, nil
}
