package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync/atomic"

	"github.com/rosenhouse/tls-tunnel-experiments/lib"
)

func main() {
	var tlsFlags lib.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "server")

	var socketFlags lib.SocketFlags
	socketFlags.AddFlags(flag.CommandLine, "listen")
	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	listener, err := tls.Listen("tcp", socketFlags.Address, tlsConfig)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %s", err)
		}
		go handleConnection(conn)
	}
}

var counter int32

func handleConnection(conn net.Conn) {
	c := atomic.AddInt32(&counter, 1)
	msg := strings.NewReader(fmt.Sprintf("hello from server, you are connection %d\n", c))
	err := lib.CopyBoth(conn, msg, os.Stdout)
	if err != nil {
		log.Fatalf("copy: %s", err)
	}
}
