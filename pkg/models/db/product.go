package db

import "time"

// Product represents the products table
type Product struct {
	ID        uint         `gorm:"primaryKey;autoIncrement"` // use uint for auto-increment
	Name      string       `gorm:"size:255;not null"`
	Price     float64      `gorm:"not null"`
	Category  string       `gorm:"size:255;not null"`
	Image     ProductImage `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE"` // one-to-one relation
	CreatedAt time.Time    `gorm:"autoCreateTime"`
}

// ProductImage represents the product_images table
type ProductImage struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	ProductID uint      `gorm:"not null;uniqueIndex"` // one-to-one relation with Product
	Thumbnail string    `gorm:"size:255;not null"`
	Mobile    string    `gorm:"size:255;not null"`
	Tablet    string    `gorm:"size:255;not null"`
	Desktop   string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
