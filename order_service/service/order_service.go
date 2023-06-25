package service

import (
	"errors"
	"order_service/model"
	"order_service/repo"
	"order_service/util/saga"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"
)

type OrderService interface {
	AddOrder(itemID, itemQty, voucherID int64) (*model.Order, error)
}

type orderService struct {
	itemService    ItemService
	voucherService VoucherService

	orderRepo repo.OrderRepo

	snowflakeGenerator *snowflake.Node
}

func NewOrderService(itemService ItemService, voucherService VoucherService, orderRepo repo.OrderRepo) OrderService {
	nodeID, _ := strconv.Atoi(os.Getenv("node_id"))
	g, _ := snowflake.NewNode(int64(nodeID))

	return &orderService{
		itemService:        itemService,
		voucherService:     voucherService,
		orderRepo:          orderRepo,
		snowflakeGenerator: g,
	}
}

func (s *orderService) AddOrder(itemID, itemQty, voucherID int64) (*model.Order, error) {
	txnID := s.snowflakeGenerator.Generate().String()

	sagaExecutor := saga.NewSaga()
	sagaExecutor.AddStep(
		func() error {
			return s.itemService.DeductItem(itemID, itemQty, txnID)
		},
		func() error {
			return s.itemService.DeductItemCancel(txnID)
		},
	)
	if voucherID > 0 {
		sagaExecutor.AddStep(
			func() error {
				return s.voucherService.UseVoucher(voucherID, txnID)
			},
			func() error {
				return s.voucherService.UseVoucherCancel(txnID)
			},
		)
	}
	sagaExecutor.AddStep(
		func() error {
			_, err := s.orderRepo.AddOrder(txnID, itemID, itemQty, voucherID)
			return err
		},
		func() error {
			return nil
		},
	)

	if isCompleted, err := sagaExecutor.Execute(); err != nil {
		return nil, err
	} else if !isCompleted {
		return nil, errors.New("failed to place order")
	}

	return &model.Order{
		TxnID:     txnID,
		ItemID:    itemID,
		ItemQty:   itemQty,
		VoucherID: voucherID,
	}, nil
}
