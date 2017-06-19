package main

import (
	"crypto/tls"
	"flag"
	"log"
	"os"
	"time"

	"github.com/rosenhouse/tls-tunnel-experiments/lib"
)

func main() {
	var tlsFlags lib.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "client")

	var socketFlags lib.SocketFlags
	socketFlags.AddFlags(flag.CommandLine, "remote")
	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	conn, err := tls.Dial("tcp", socketFlags.Address, tlsConfig)
	if err != nil {
		log.Fatalf("dial: %s", err)
	}
	defer conn.Close()

	go func() {
		err = lib.CopyBoth(conn, os.Stdin, os.Stdout)
		if err != nil {
			log.Fatalf("copy: %s", err)
		}
	}()
	time.Sleep(100 * time.Millisecond)
}
