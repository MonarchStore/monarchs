package main

import (
	"flag"
	"log"
	"os"
	
	"github.com/arturom/monarchs/server"
)

func main() {
	opts := readCLIOptions()

	log.Printf("Listening for http connections on %s", *opts.addr)
	if err := server.NewHttpServer().Listen(*opts.addr); err != nil {
		log.Printf("Failed to open connection. %s", err)
	}
}

type cliOpts struct {
	addr *string
}

func readCLIOptions() cliOpts {
	default_port := ":6789"
	if port, ok := os.LookupEnv("LISTEN_PORT"); ok {
		default_port = port
	}

	addrPtr := flag.String("addr", default_port, "The binding address")
	flag.Parse()

	return cliOpts{
		addr: addrPtr,
	}
}
