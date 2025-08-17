package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models"
	"gorm.io/gorm"
)

type ProductRepo struct {
	db *gorm.DB
}

func NewProductRepo(db *gorm.DB) *ProductRepo {
	return &ProductRepo{db: db}
}

func (r *ProductRepo) FindAll() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Preload("Products").Find(&products).Error; err != nil {
		return nil, constants.ErrProductNotFound
	}

	return products, nil
}

func (r *ProductRepo) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Image").First(&product, "id = ?", id).Error; err != nil {
		log.Error().Msgf("Error fetching product with ID %d: %v", id, err)
		return nil, constants.ErrProductNotFound
	}

	return &product, nil
}
