package server

import (
	"net/http"
	"sync"

	ds "github.com/MonarchStore/monarchs/docstore"
	log "github.com/sirupsen/logrus"
)

type HTTPServer interface {
	HandleSigterm(srv *http.Server)
	Listen(addr string, stopChan <-chan struct{}) error
}

func NewHttpServer() HTTPServer {
	log.Trace("Creating new HTTP Server")
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
	log.Debug("Initialized HTTP server")

	// Set routes
	http.HandleFunc("/", s.dataHandler)
	http.HandleFunc("/healthz", s.healthCheck)
	http.HandleFunc("/metricz", s.doMetrics)
	log.Debug("Registered routes")

	// Shutdown the server when this function (s.Run) returns
	defer s.HandleSigterm(srv)

	// Go ListenAndServe asynchronously
	go func() {
		log.Println("[ListenAndServe]")
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("HTTPServer.ListenAndServe Error: %s", err)
		}
	}()

	// Wait for a stopChan message
	for {
		select {
		case <-stopChan:
			log.Println("Caught SIGTERM. Exiting...")
			return
		}
	}
	return
}
