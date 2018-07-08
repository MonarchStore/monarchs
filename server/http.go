package server

import (
	"net/http"
	"sync"

	ds "github.com/arturom/monarchs/docstore"
)

type HTTPServer interface {
	Listen(addr string) error
}

func NewHttpServer() HTTPServer {
	return &httpServer{
		storeMap: make(ds.StoreMap),
		mutex:    sync.RWMutex{},
	}
}

type httpServer struct {
	storeMap ds.StoreMap
	mutex    sync.RWMutex
}

func (s *httpServer) Listen(addr string) error {
	http.HandleFunc("/", s.dataHandler)
	return http.ListenAndServe(addr, Logger(http.DefaultServeMux))
}
