package db

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID        uuid.UUID       `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Total     float64         `gorm:"not null"`
	Discounts float64         `gorm:"not null"`
	Products  []*OrderProduct `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
}

type OrderProduct struct {
	ID        int       `gorm:"primaryKey;autoIncrement"`
	OrderID   uuid.UUID `gorm:"type:uuid;not null;index"`
	ProductID string    `gorm:"not null" json:"productId" validate:"required"`
	Quantity  int       `gorm:"not null" json:"quantity" validate:"required,min=1"`
}
