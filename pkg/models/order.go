package models

type Order struct {
	ID        string      `json:"id"`
	Total     float64     `json:"total"`
	Discounts []string    `json:"discounts,omitempty"`
	Items     []OrderItem `json:"items"`
	Products  []*Product  `json:"products"`
}

type OrderItem struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}

type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items" validate:"required,dive,required"`
}
