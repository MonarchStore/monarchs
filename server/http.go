package server

import (
	"net/http"

	ds "bitbucket.org/enticusa/kingdb/docstore"
)

type HTTPServer interface {
	Listen(addr string) error
}

func NewHttpServer(store ds.Store) HTTPServer {
	return &httpServer{
		store: store,
	}
}

type httpServer struct {
	store ds.Store
}

func (s *httpServer) Listen(addr string) error {
	http.HandleFunc("/data/", s.dataHandler)
	http.HandleFunc("/schema/", s.schemaHandler)
	return http.ListenAndServe(addr, Logger(http.DefaultServeMux))
}

func (s *httpServer) schemaHandler(w http.ResponseWriter, r *http.Request) {
}
