package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/models"
	"gorm.io/gorm"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) *OrderRepo {
	return &OrderRepo{db: db}
}

// Create inserts a new order with products
func (r *OrderRepo) Create(order models.Order) (string, error) {
	if err := r.db.Create(&order).Error; err != nil {
		return "", err
	}

	return order.ID.String(), nil
}
