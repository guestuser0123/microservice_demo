package handler

import (
	"encoding/json"
	"net/http"
	"order_service/adapter"
	"order_service/repo"
	"order_service/service"
)

type OrderServiceHandler struct {
	orderService service.OrderService

	closeFn func() error
}

func NewOrderServiceHandler() *OrderServiceHandler {
	itemRepo := repo.NewItemRepo(adapter.NewHTTPClient())
	voucherRepo := repo.NewVoucherRepo(adapter.NewHTTPClient())
	orderRepo := repo.NewRepo()
	orderRepo.MustInit()

	itemService := service.NewItemService(itemRepo)
	voucherService := service.NewVoucherService(voucherRepo)
	orderService := service.NewOrderService(itemService, voucherService, orderRepo)

	return &OrderServiceHandler{
		orderService: orderService,
		closeFn:      orderRepo.Close,
	}
}

func (h *OrderServiceHandler) Close() error {
	return h.closeFn()
}

type PlaceOrderRequest struct {
	ItemID    int64 `json:"item_id"`
	ItemQty   int64 `json:"item_qty"`
	VoucherID int64 `json:"voucher_id"`
}

func (h *OrderServiceHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	var req PlaceOrderRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	order, err := h.orderService.AddOrder(req.ItemID, req.ItemQty, req.VoucherID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
