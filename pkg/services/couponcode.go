package services

import "github.com/malakagl/kart-challenge/pkg/repositories"

type CouponCodeService struct {
	repo repositories.CouponCodeRepo
}

func NewCouponCodeService(r repositories.CouponCodeRepo) *CouponCodeService {
	return &CouponCodeService{repo: r}
}

func (s *CouponCodeService) CountFilesByCode(code string) int64 {
	return s.repo.CountFilesByCode(code)
}
