package services

import (
	"strconv"

	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
	"github.com/malakagl/kart-challenge/pkg/repositories"
)

type ProductService struct {
	repo repositories.ProductRepo
}

func NewProductService(r repositories.ProductRepo) ProductService {
	return ProductService{repo: r}
}

func (s *ProductService) FindAll() (*response.ProductsResponse, error) {
	res, err := s.repo.FindAll()
	if err != nil {
		log.Error().Msgf("findAll failed with error: %v", err)
		return nil, constants.ErrProductNotFound
	}

	if res == nil {
		log.Error().Msg("findAll returned 0 elements")
		return nil, constants.ErrProductNotFound
	}

	products := make([]response.Product, len(res))
	for i, p := range res {
		products[i] = response.Product{
			ID:       strconv.FormatUint(uint64(p.ID), 10),
			Name:     p.Name,
			Price:    p.Price,
			Category: p.Category,
			Image: response.ProductImage{
				Thumbnail: p.Image.Thumbnail,
				Mobile:    p.Image.Mobile,
				Tablet:    p.Image.Tablet,
				Desktop:   p.Image.Desktop,
			},
		}
	}

	return &response.ProductsResponse{Products: products}, nil
}

func (s *ProductService) FindByID(id uint) (*response.ProductResponse, error) {
	res, err := s.repo.FindByID(id)
	if err != nil {
		log.Error().Msgf("findAll failed with error: %v", err)
		return nil, constants.ErrProductNotFound
	}

	if res == nil {
		log.Error().Msg("findAll returned 0 elements")
		return nil, constants.ErrProductNotFound
	}

	return &response.ProductResponse{
		ID:       strconv.FormatUint(uint64(res.ID), 10),
		Name:     res.Name,
		Price:    res.Price,
		Category: res.Category,
		Image: response.ProductImage{
			Thumbnail: res.Image.Thumbnail,
			Mobile:    res.Image.Mobile,
			Tablet:    res.Image.Tablet,
			Desktop:   res.Image.Desktop,
		},
	}, nil
}
