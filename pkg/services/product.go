package services

import (
	"github.com/malakagl/kart-challenge/pkg/models"
)

type ProductRepository interface {
	FindAll() ([]models.Product, error)
	FindByID(id string) (*models.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(r ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) FindAll() ([]models.Product, error) {
	return s.repo.FindAll()
}

func (s *ProductService) FindByID(id string) (*models.Product, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}
