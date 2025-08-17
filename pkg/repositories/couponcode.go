package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/models"
	"gorm.io/gorm"
)

type CouponCodeRepository struct {
	db *gorm.DB
}

func NewCouponCodeRepository(db *gorm.DB) *CouponCodeRepository {
	return &CouponCodeRepository{db: db}
}

func (r *CouponCodeRepository) CountFilesByCode(code string) int64 {
	var count int64
	r.db.Model(&models.CouponCode{}).
		Where("code = ?", code).
		Distinct("file_id").
		Count(&count)
	return count
}
