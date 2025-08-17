package constants

import "errors"

var (
	ErrProductNotFound     = errors.New("product not found")
	ErrInvalidCouponCode   = errors.New("invalid coupon code")
	ErrInvalidProductID    = errors.New("invalid product ID")
	ErrInternalServerError = errors.New("internal server error")
)
