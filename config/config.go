package config

import (
	"flag"
	"os"
)

type CLIOptions struct {
	ListenAddress *string
}

func ParseFlags() CLIOptions {
	default_port := ":6789"
	if port, ok := os.LookupEnv("LISTEN_PORT"); ok {
		default_port = port
	}

	addrPtr := flag.String("addr", default_port, "The binding address")
	flag.Parse()

	return CLIOptions{
		ListenAddress: addrPtr,
	}
}
