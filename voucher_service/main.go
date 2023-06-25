package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"voucher_service/handler"
)

func main() {
	h := handler.NewVoucherServiceHandler()
	defer func() {
		if err := h.Close(); err != nil {
			log.Println("VoucherServiceHandler closed with error=", err)
		}
	}()

	http.HandleFunc("/voucher/get", h.GetVoucherByID)
	http.HandleFunc("/voucher/use", h.UseVoucher)
	http.HandleFunc("/voucher/cancel", h.UseVoucherCancel)
	log.Fatal(http.ListenAndServe(os.Getenv("VOUCHER_SERVICE_PUBLIC_ADDR"), nil))
}
