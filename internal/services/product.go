package services

import "github.com/MalakaGL/kart-challenge-lahiru/kart-challenge/internal/models"

type ProductRepository interface {
	FindAll() ([]models.Product, error)
}

type ProductService struct {
	repo ProductRepository
}

func NewProductService(r ProductRepository) *ProductService {
	return &ProductService{repo: r}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.repo.FindAll()
}
