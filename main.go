package main

import (
	"flag"
	"log"

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
	addrPtr := flag.String("addr", ":6789", "The binding address")
	flag.Parse()

	return cliOpts{
		addr: addrPtr,
	}
}
