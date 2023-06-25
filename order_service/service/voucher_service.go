package service

import (
	"order_service/repo"
	"strconv"
)

type VoucherService interface {
	UseVoucher(voucherID int64, txnID string) error
	UseVoucherCancel(txnID string) error
}

type voucherService struct {
	voucherRepo repo.VoucherRepo
}

func NewVoucherService(r repo.VoucherRepo) VoucherService {
	return &voucherService{
		voucherRepo: r,
	}
}

func (s *voucherService) UseVoucher(voucherID int64, txnID string) error {
	return s.voucherRepo.UseVoucher(strconv.FormatInt(voucherID, 10), txnID)
}

func (s *voucherService) UseVoucherCancel(txnID string) error {
	return s.voucherRepo.UseVoucherCancel(txnID)
}
