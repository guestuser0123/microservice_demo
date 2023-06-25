package adapter

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

func NewHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns, _ = strconv.Atoi(os.Getenv("HTTP_MAX_IDLE_CONNS"))
	t.MaxConnsPerHost, _ = strconv.Atoi(os.Getenv("HTTP_MAX_CONNS_PER_HOST"))
	t.MaxIdleConnsPerHost, _ = strconv.Atoi(os.Getenv("HTTP_MAX_IDLE_CONNS_PER_HOST"))

	return &http.Client{
		Timeout:   5 * time.Second,
		Transport: t,
	}
}
