package models

import "time"

type CouponCode struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Code      string    `gorm:"size:50;not null;index" json:"code"` // index for searching by code
	FileID    uint      `json:"file_id"`
	File      File      `gorm:"constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

type File struct {
	ID          uint         `gorm:"primaryKey;autoIncrement"`
	FileName    string       `gorm:"size:255;not null" json:"file_name"`
	CouponCodes []CouponCode `gorm:"foreignKey:FileID"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
}
