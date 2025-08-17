package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/models/db"
	"gorm.io/gorm"
)

type CouponCodeRepo struct {
	db *gorm.DB
}

func NewCouponCodeRepository(db *gorm.DB) CouponCodeRepo {
	return CouponCodeRepo{db: db}
}

func (r *CouponCodeRepo) CountFilesByCode(code string) int64 {
	var count int64
	r.db.Model(&db.CouponCode{}).
		Where("code = ?", code).
		Distinct("file_id").
		Count(&count)
	return count
}
