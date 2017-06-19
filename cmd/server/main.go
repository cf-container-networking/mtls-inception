package main

import (
	"flag"
	"log"

	"github.com/rosenhouse/tls-tunnel-experiments/config"
)

func main() {
	var tlsFlags config.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "server")
	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	log.Printf("tlsConfig: %+v\n", tlsConfig)
}
