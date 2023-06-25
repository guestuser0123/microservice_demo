package service

import (
	"item_service/model"
	"item_service/repo"
)

type Service interface {
	GetItemByID(itemID int64) (*model.Item, error)
	DeductItem(itemID, qty int64, txnID string) error
	ReturnItem(txnID string) error
}

type itemService struct {
	repo repo.ItemRepo
}

func NewItemService(repo repo.ItemRepo) Service {
	return &itemService{
		repo: repo,
	}
}

func (s *itemService) GetItemByID(itemID int64) (*model.Item, error) {
	return s.repo.GetItemByID(itemID)
}

func (s *itemService) DeductItem(itemID, qty int64, txnID string) error {
	return s.repo.DeductItem(itemID, qty, txnID)
}

func (s *itemService) ReturnItem(txnID string) error {
	return s.repo.ReturnItem(txnID)
}
