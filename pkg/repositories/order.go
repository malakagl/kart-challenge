package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/models/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepo struct {
	db *gorm.DB
}

func NewOrderRepo(db *gorm.DB) OrderRepo {
	return OrderRepo{db: db}
}

// Create inserts a new order with products
func (r *OrderRepo) Create(order *db.Order) error {
	if err := r.db.Clauses(clause.Returning{}).Create(&order).Error; err != nil {
		return err
	}

	return nil
}
