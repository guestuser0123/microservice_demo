package service

import (
	"voucher_service/model"
	"voucher_service/repo"
)

type Service interface {
	GetVoucherByID(voucherID int64) (*model.Voucher, error)
	UseVoucher(itemID int64, txnID string) error
	ReturnVoucher(txnID string) error
}

type voucherService struct {
	voucherRepo repo.VoucherRepo
}

const (
	UseVoucherCount = 1
)

func NewVoucherService(voucherRepo repo.VoucherRepo) Service {
	return &voucherService{
		voucherRepo: voucherRepo,
	}
}

func (s *voucherService) GetVoucherByID(voucherID int64) (*model.Voucher, error) {
	return s.voucherRepo.GetVoucherByID(voucherID)
}

func (s *voucherService) UseVoucher(voucherID int64, txnID string) error {
	return s.voucherRepo.UseVoucher(voucherID, UseVoucherCount, txnID)
}

func (s *voucherService) ReturnVoucher(txnID string) error {
	return s.voucherRepo.ReturnVoucher(txnID)
}
