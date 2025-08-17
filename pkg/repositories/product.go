package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/db"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) ProductRepo {
	return ProductRepo{db: db}
}

func (r *ProductRepo) FindAll() ([]db.Product, error) {
	var products []db.Product
	if err := r.db.Preload("Image").Find(&products).Error; err != nil {
		return nil, constants.ErrProductNotFound
	}

	return products, nil
}

func (r *ProductRepo) FindByID(id uint) (*db.Product, error) {
	var product db.Product
	if err := r.db.Preload("Image").First(&product, "id = ?", id).Error; err != nil {
		log.Error().Msgf("Error fetching product with ID %d: %v", id, err)
		return nil, constants.ErrProductNotFound
	}

	return &product, nil
}
