package services

type CouponCodeRepository interface {
	CountFilesByCode(code string) int64
}

type CouponCodeService struct {
	repo CouponCodeRepository
}

func NewCouponCodeService(r CouponCodeRepository) *CouponCodeService {
	return &CouponCodeService{repo: r}
}

func (s *CouponCodeService) CountFilesByCode(code string) int64 {
	return s.repo.CountFilesByCode(code)
}
