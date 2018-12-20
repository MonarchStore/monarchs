package server

import (
	"log"
	"net/http"
	"sync"

	ds "github.com/MonarchStore/monarchs/docstore"
	"github.com/sirupsen/logrus"
)

type HTTPServer interface {
	Listen(addr string, stopChan <-chan struct{}) error
}

func NewHttpServer() HTTPServer {
	logrus.Trace("Creating new HTTP Server")
	return &httpServer{
		storeMap: make(ds.StoreMap),
		mutex:    sync.RWMutex{},
	}
}

type httpServer struct {
	storeMap ds.StoreMap
	mutex    sync.RWMutex
}

func (s *httpServer) Listen(addr string, stopChan <-chan struct{}) (err error) {
	// Init http.Server
	logger := logrus.New()
	lw := logger.Writer()
	defer lw.Close()

	http_logger := log.New(lw, "", 0)

	srv := &http.Server{
		Addr:     addr, // ':6789'
		ErrorLog: http_logger,
	}
	logrus.Debug("Initialized HTTP server")

	// Set routes
	http.HandleFunc("/", s.dataHandler)
	http.HandleFunc("/healthz", s.healthCheck)
	http.HandleFunc("/metricz", s.doMetrics)
	logrus.Debug("Registered routes")

	// Go ListenAndServe asynchronously
	go func() {
		logrus.Println("[ListenAndServe]")
		if err := srv.ListenAndServe(); err != nil {
			logrus.Printf("HTTPServer.ListenAndServe Error: %s", err)
		}
	}()

	// Wait for a stopChan message
	for {
		select {
		case <-stopChan:
			logrus.Println("Caught SIGTERM. Exiting...")
			return
		}
	}
	return
}
