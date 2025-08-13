package repositories

import (
	"github.com/malakagl/kart-challenge/internal/models"
	"github.com/malakagl/kart-challenge/pkg/constants"
)

type InMemoryProductRepo struct {
	products []models.Product
}

func NewInMemoryProductRepo(products []models.Product) *InMemoryProductRepo {
	return &InMemoryProductRepo{products: products}
}

func (r *InMemoryProductRepo) FindAll() ([]models.Product, error) {
	return r.products, nil
}

func (r *InMemoryProductRepo) FindByID(id string) (*models.Product, error) {
	for _, p := range r.products {
		if p.ID == id {
			return &p, nil
		}
	}

	return nil, constants.ErrProductNotFound
}
