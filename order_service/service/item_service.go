package service

import (
	"order_service/repo"
	"strconv"
)

type ItemService interface {
	DeductItem(itemID, itemQty int64, txnID string) error
	DeductItemCancel(txnID string) error
}

type itemService struct {
	itemRepo repo.ItemRepo
}

func NewItemService(r repo.ItemRepo) ItemService {
	return &itemService{
		itemRepo: r,
	}
}

func (s *itemService) DeductItem(itemID, itemQty int64, txnID string) error {
	return s.itemRepo.DeductItem(strconv.FormatInt(itemID, 10), strconv.FormatInt(itemQty, 10), txnID)
}

func (s *itemService) DeductItemCancel(txnID string) error {
	return s.itemRepo.DeductItemCancel(txnID)
}
