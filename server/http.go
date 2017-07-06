package server

import (
	"net/http"

	ds "bitbucket.org/enticusa/kingdb/docstore"
)

type HTTPServer interface {
	Listen(addr string) error
}

func NewHttpServer() HTTPServer {
	return &httpServer{
		storeMap: make(ds.StoreMap),
	}
}

type httpServer struct {
	storeMap ds.StoreMap
}

func (s *httpServer) Listen(addr string) error {
	http.HandleFunc("/", s.dataHandler)
	return http.ListenAndServe(addr, Logger(http.DefaultServeMux))
}
