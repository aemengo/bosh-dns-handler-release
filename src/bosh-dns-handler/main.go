package main

import (
	cfg "github.com/aemengo/bosh-dns-handler/config"
	"github.com/aemengo/bosh-dns-handler/handler"
	"github.com/miekg/dns"
	"log"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("[USAGE] %s <config-path>", os.Args[0])
	}

	config, err := cfg.New(os.Args[1])
	expectNoError(err)

	srv := &dns.Server{Addr: ":" + config.Port, Net: "udp"}
	srv.Handler = handler.New(config)

	log.Printf("Launching bosh-dns-handler on :%s...\n", config.Port)
	err = srv.ListenAndServe()
	expectNoError(err)
}

func expectNoError(err error) {
	if err != nil {
		log.Fatalf("Failed to initialize: %s\n", err)
	}
}