package server

import (
	"log"
	"net/http"
	"sync"

	ds "github.com/MonarchStore/monarchs/docstore"
)

type HTTPServer interface {
	HandleSigterm(srv *http.Server)
	Listen(addr string, stopChan <-chan struct{}) error
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

// An opportunity to gracefully shutdown
func (s *httpServer) HandleSigterm(srv *http.Server) {
	log.Println("Goodbye")
	srv.Shutdown(nil)
}

func (s *httpServer) Listen(addr string, stopChan <-chan struct{}) (err error) {
	// Init http.Server
	srv := &http.Server{Addr: addr}

	// Set routes
	http.HandleFunc("/", s.dataHandler)
	http.HandleFunc("/healthz", s.healthCheck)
	http.HandleFunc("/metricz", s.doMetrics)

	// Shutdown the server when this function (s.Run) returns
	defer s.HandleSigterm(srv)

	// Go ListenAndServe asynchronously
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("HTTPServer.ListenAndServe Error: %s", err)
		}
	}()

	// Wait for a stopChan message
	for {
		select {
		case <-stopChan:
			log.Println("Caught SIGTERM")
			return
		}
	}
}
