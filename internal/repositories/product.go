package repositories

import "github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/models"

type InMemoryProductRepo struct {
	products []models.Product
}

func NewInMemoryProductRepo(products []models.Product) *InMemoryProductRepo {
	return &InMemoryProductRepo{products: products}
}

func (r *InMemoryProductRepo) FindAll() ([]models.Product, error) {
	return r.products, nil
}
