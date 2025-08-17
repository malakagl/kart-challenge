package request

type OrderRequest struct {
	CouponCode string `json:"couponCode,omitempty"`
	Items      []Item `json:"items" validate:"required,dive,required"`
}

type Item struct {
	ProductID string `json:"productId" validate:"required"`
	Quantity  int    `json:"quantity" validate:"required,min=1"`
}
