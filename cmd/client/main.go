package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/cf-container-networking/tls-tunnel-experiments/lib"
)

func main() {
	var tlsFlags lib.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "client", "ca")

	var remoteAddr string
	flag.StringVar(&remoteAddr, "remoteAddr", "127.0.0.21:7021", "remote to connect to (client proxy listener)")
	flag.Parse()

	var clientAddr string
	flag.StringVar(&clientAddr, "clientAddr", "127.0.0.11:0", "local ip address to use when initiating connection to remote server")
	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	localAddr, err := net.ResolveTCPAddr("tcp", clientAddr)
	if err != nil {
		log.Fatalf("parsing client ip as tcp addr: %s", err)
	}

	dialer := &net.Dialer{LocalAddr: localAddr}
	conn, err := tls.DialWithDialer(dialer, "tcp", remoteAddr, tlsConfig)
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
