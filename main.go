package main

import (
	"flag"
	"log"
	"strings"

	ds "bitbucket.org/enticusa/kingdb/docstore"
	"bitbucket.org/enticusa/kingdb/server"
)

func main() {
	opts := readCLIOptions()
	store := ds.NewStore(opts.labels)
	log.Printf("Schema defined with %d entities: %s", len(opts.labels), opts.labels)

	log.Printf("Listening for http connections on %s", *opts.addr)
	if err := server.NewHttpServer(store).Listen(*opts.addr); err != nil {
		log.Printf("Failed to open connection. %s", err)
	}
}

type cliOpts struct {
	labels ds.Labels
	addr   *string
}

func readCLIOptions() cliOpts {
	joinedLabelsPtr := flag.String("labels", "", "A comma-separated, ordered list of the labels of elements in the hieararchy")
	addrPtr := flag.String("addr", ":6789", "The binding address")
	flag.Parse()

	labels := strings.Split(*joinedLabelsPtr, ",")
	hierarchyLabels := make(ds.Labels, len(labels))
	for i, label := range labels {
		hierarchyLabels[i] = ds.Label(label)
	}

	return cliOpts{
		labels: hierarchyLabels,
		addr:   addrPtr,
	}
}

func createStore(opts cliOpts) ds.Store {
	return ds.NewStore(opts.labels)
}
