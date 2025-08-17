package response

type Product struct {
	ID       string       `json:"id"`
	Name     string       `gorm:"size:255;not null"`
	Price    float64      `gorm:"not null"`
	Category string       `gorm:"size:255;not null"`
	Image    ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // one-to-one relation
}

type ProductImage struct {
	Thumbnail string `json:"thumbnail"`
	Mobile    string `json:"mobile"`
	Tablet    string `json:"tablet"`
	Desktop   string `json:"desktop"`
}

type Products []Product

type ProductsResponse struct {
	Products
}

type ProductResponse Product
