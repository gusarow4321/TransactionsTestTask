package server

import (
	"net/http"
	"time"
)

type Service interface {
	AddTx(http.ResponseWriter, *http.Request)
}

func NewServer(addr string, srv Service) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/add-tx", srv.AddTx)

	s := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s
}
