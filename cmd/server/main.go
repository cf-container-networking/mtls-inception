package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"

	"github.com/cf-container-networking/mtls-inception/lib"
)

func main() {
	var tlsFlags lib.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "server", "ca")

	var listenAddr string
	flag.StringVar(&listenAddr, "listenAddr", "127.0.0.12:7012", "local listen address for the end server")
	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	listener, err := tls.Listen("tcp", listenAddr, tlsConfig)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %s", err)
		}
		go handleConnection(conn.(*tls.Conn))
	}
}

var counter int32

func handleConnection(conn *tls.Conn) {
	c := atomic.AddInt32(&counter, 1)

	msg := strings.NewReader(fmt.Sprintf(`hello from server,
	you are coming from %s
	you are connection %d
`, conn.RemoteAddr(), c))

	err := lib.CopyBoth(conn, msg, os.Stdout)
	if err != nil {
		log.Printf("copy: %s", err)
	}
}
