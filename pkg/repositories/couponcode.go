package repositories

import (
	"github.com/malakagl/kart-challenge/pkg/constants"
	"github.com/malakagl/kart-challenge/pkg/log"
	"github.com/malakagl/kart-challenge/pkg/models/db"
	"gorm.io/gorm"
)

type CouponCodeRepo struct {
	db *gorm.DB
}

func NewCouponCodeRepository(db *gorm.DB) CouponCodeRepo {
	return CouponCodeRepo{db: db}
}

func (r *CouponCodeRepo) CountFilesByCode(code string) (int64, error) {
	var count int64
	err := r.db.Model(&db.CouponCode{}).
		Where("code = ?", code).
		Distinct("file_id").
		Count(&count).Error
	if err != nil {
		log.Error().Msgf("error counting coupon code: %v", err)
		return 0, constants.ErrDatabaseError
	}

	return count, nil
}
