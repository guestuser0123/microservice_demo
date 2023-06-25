package main

import (
	"item_service/handler"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	h := handler.NewItemServiceHandler()
	defer func() {
		if err := h.Close(); err != nil {
			log.Println("ItemServiceHandler closed with error=", err)
		}
	}()

	http.HandleFunc("/item/get", h.GetItemByID)
	http.HandleFunc("/item/deduct", h.DeductItem)
	http.HandleFunc("/item/cancel", h.DeductItemCancel)
	log.Fatal(http.ListenAndServe(os.Getenv("ITEM_SERVICE_PUBLIC_ADDR"), nil))
}
