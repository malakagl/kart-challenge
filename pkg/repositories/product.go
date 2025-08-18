package repositories

import (
	"context"

	"github.com/malakagl/kart-challenge/pkg/errors"
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

func (r *ProductRepo) FindAll(ctx context.Context) ([]db.Product, error) {
	var products []db.Product
	if err := r.db.WithContext(ctx).Preload("Image").Find(&products).Error; err != nil {
		return nil, errors.ErrProductNotFound
	}

	return products, nil
}

func (r *ProductRepo) FindByID(ctx context.Context, id uint) (*db.Product, error) {
	var product db.Product
	if err := r.db.WithContext(ctx).Preload("Image").First(&product, "id = ?", id).Error; err != nil {
		log.Error().Msgf("Error fetching product with ID %d: %v", id, err)
		if err.Error() == "record not found" {
			return nil, errors.ErrProductNotFound
		}

		return nil, errors.ErrDatabaseError
	}

	return &product, nil
}
