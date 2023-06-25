package handler

import (
	"encoding/json"
	"item_service/repo"
	"item_service/service"
	"net/http"
	"strconv"
)

type ItemServiceHandler struct {
	itemService service.Service

	closeFn func() error
}

func NewItemServiceHandler() *ItemServiceHandler {
	r := repo.NewRepo()
	r.MustInit()

	s := service.NewItemService(r)

	return &ItemServiceHandler{
		itemService: s,
		closeFn:     r.Close,
	}
}

func (h *ItemServiceHandler) Close() error {
	return h.closeFn()
}

func (h *ItemServiceHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	itemID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	item, err := h.itemService.GetItemByID(int64(itemID))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ItemServiceHandler) DeductItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	itemID, _ := strconv.Atoi(r.URL.Query().Get("itemid"))
	itemQty, _ := strconv.Atoi(r.URL.Query().Get("itemqty"))
	txnID := r.URL.Query().Get("txnid")

	err := h.itemService.DeductItem(int64(itemID), int64(itemQty), txnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *ItemServiceHandler) DeductItemCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	txnID := r.URL.Query().Get("txnid")
	err := h.itemService.ReturnItem(txnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
