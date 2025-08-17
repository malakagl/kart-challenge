package services

import (
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/db"
	"github.com/malakagl/kart-challenge/pkg/models/dto/request"
	"github.com/malakagl/kart-challenge/pkg/models/dto/response"
	"github.com/malakagl/kart-challenge/pkg/repositories"
	"github.com/malakagl/kart-challenge/pkg/util"
)

type OrderService struct {
	orderRepo      repositories.OrderRepo
	couponCodeRepo repositories.CouponCodeRepo
	productRepo    repositories.ProductRepo
}

func NewOrderService(
	r repositories.OrderRepo,
	c repositories.CouponCodeRepo,
	p repositories.ProductRepo,
) OrderService {
	return OrderService{
		orderRepo:      r,
		couponCodeRepo: c,
		productRepo:    p,
	}
}

func (o *OrderService) isCouponCodeValid(code string) bool {
	if len(code) < 8 || len(code) > 10 {
		return false
	}

	if o.couponCodeRepo.CountFilesByCode(code) > 1 {
		log.Error().Msgf("Coupon code %s is valid: found in multiple files", code)
		return true
	}

	return false
}

func (o *OrderService) Create(req *request.OrderRequest) (*response.OrderResponse, error) {
	if !o.isCouponCodeValid(req.CouponCode) { // !h.couponValidator.ValidateCouponCode(orderReq.CouponCode)
		log.Error().Msgf("Invalid coupon code: %s", req.CouponCode)
		return nil, constants.ErrInvalidCouponCode
	}

	order := db.Order{}
	orderProducts := make([]*db.OrderProduct, len(req.Items))
	products := make([]response.Product, len(req.Items))
	items := make([]response.Item, len(req.Items))
	for i, item := range req.Items {
		productId, err := util.StringToUint(item.ProductID)
		if err != nil || productId == 0 {
			log.Error().Msg("Invalid product ID in order request")
			return nil, constants.ErrInvalidProductID
		}

		product, err := o.productRepo.FindByID(productId)
		if err != nil && err.Error() == "not found" {
			log.Error().Msgf("Error fetching product: %v", err)
			return nil, constants.ErrProductNotFound
		}
		if err != nil {
			log.Error().Msgf("Error fetching product: %v", err)
			return nil, constants.ErrInternalServerError
		}

		orderProducts[i] = &db.OrderProduct{ProductID: item.ProductID, Quantity: item.Quantity}
		order.Total += product.Price * float64(item.Quantity)
		products[i] = response.Product{
			ID:       item.ProductID,
			Name:     product.Name,
			Price:    product.Price,
			Category: product.Category,
			Image: response.ProductImage{
				Thumbnail: product.Image.Thumbnail,
				Mobile:    product.Image.Mobile,
				Tablet:    product.Image.Tablet,
				Desktop:   product.Image.Desktop,
			},
		}
		items[i] = response.Item{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}
	order.Products = orderProducts

	err := o.orderRepo.Create(&order)
	if err != nil {
		log.Error().Msgf("Error creating order: %v", err)
		return nil, constants.ErrInternalServerError
	}

	return &response.OrderResponse{
		ID:        order.ID.String(),
		Total:     order.Total,
		Discounts: order.Discounts,
		Items:     items,
		Products:  products,
	}, nil
}
