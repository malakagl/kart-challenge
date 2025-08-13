package models

type Order struct {
	ID         string      `json:"id"`
	Items      []OrderItem `json:"items"`
	CouponCode string      `json:"couponCode,omitempty"`
}

type OrderItem struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderRequest struct {
	CouponCode string      `json:"couponCode,omitempty"`
	Items      []OrderItem `json:"items" validate:"required,dive,required"`
}

type OrderResponse struct {
	ID       string      `json:"id"`
	Items    []OrderItem `json:"items"`
	Products []*Product  `json:"products"`
}
