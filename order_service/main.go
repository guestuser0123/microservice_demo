package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"order_service/handler"
	"os"
)

func main() {
	h := handler.NewOrderServiceHandler()
	defer func() {
		if err := h.Close(); err != nil {
			log.Println("OrderServiceHandler closed with error=", err)
		}
	}()

	http.HandleFunc("/place-order", h.PlaceOrder)
	log.Fatal(http.ListenAndServe(os.Getenv("ORDER_SERVICE_ADDR"), nil))
}
