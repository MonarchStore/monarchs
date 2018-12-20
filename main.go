package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/MonarchStore/monarchs/config"
	"github.com/MonarchStore/monarchs/server"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.NewConfigFromArgs(os.Args[1:])

	stopChan := make(chan struct{}, 1)
	go handleSigterm(stopChan)

	srv := server.NewHttpServer()
	log.Printf("Listening for http connections on %s", cfg.GetListenAddress())

	err := srv.Listen(cfg.GetListenPort(), stopChan)
	if err != nil {
		log.Fatalf("Failed to open connection. %s", err)
	}
	log.Trace("(done)")
}

func handleSigterm(stopChan chan struct{}) {
	log.Debug("Awaiting SIGTERM...")
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Println("Received SIGTERM. Terminating gracefully...")
	close(stopChan)
}
