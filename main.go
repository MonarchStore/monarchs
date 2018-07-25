package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/arturom/monarchs/config"
	"github.com/arturom/monarchs/server"
)

func main() {
	opts := config.ParseFlags()

	stopChan := make(chan struct{}, 1)
	go handleSigterm(stopChan)

	srv := server.NewHttpServer()
	log.Printf("Listening for http connections on %s", *opts.ListenAddress)

	err := srv.Listen(*opts.ListenAddress, stopChan)
	if err != nil {
		log.Fatalf("Failed to open connection. %s", err)
	}
}

func handleSigterm(stopChan chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM)
	<-signals
	log.Println("Received SIGTERM. Terminating gracefully...")
	close(stopChan)
}
