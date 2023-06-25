package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"voucher_service/repo"
	"voucher_service/service"
)

type VoucherServiceHandler struct {
	voucherService service.Service

	closeFn func() error
}

func NewVoucherServiceHandler() *VoucherServiceHandler {
	r := repo.NewVoucherRepo()
	r.MustInit()

	return &VoucherServiceHandler{
		voucherService: service.NewVoucherService(r),
		closeFn:        r.Close,
	}
}

func (h *VoucherServiceHandler) Close() error {
	return h.closeFn()
}

func (h *VoucherServiceHandler) GetVoucherByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	voucherID, _ := strconv.Atoi(r.URL.Query().Get("id"))

	voucher, err := h.voucherService.GetVoucherByID(int64(voucherID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(voucher)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *VoucherServiceHandler) UseVoucher(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	voucherID, _ := strconv.Atoi(r.URL.Query().Get("voucherid"))
	txnID := r.URL.Query().Get("txnid")
	err := h.voucherService.UseVoucher(int64(voucherID), txnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *VoucherServiceHandler) UseVoucherCancel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid http method", http.StatusMethodNotAllowed)
		return
	}

	txnID := r.URL.Query().Get("txnid")
	err := h.voucherService.ReturnVoucher(txnID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
