package response

type OrderResponse struct {
	ID        string    `json:"id"`
	Total     float64   `json:"total"`
	Discounts float64   `json:"discounts,omitempty"`
	Items     []Item    `json:"items"`
	Products  []Product `json:"products"`
}

type Item struct {
	ProductID string `json:"productId"`
	Quantity  int    `json:"quantity"`
}
