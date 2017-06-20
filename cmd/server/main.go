package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
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
		go handleConnection(conn.(*tls.Conn))
	}
}

var counter int32

func validateClientIP(conn *tls.Conn) error {
	err := conn.Handshake()
	if err != nil {
		return err
	}

	clientAddr := conn.RemoteAddr().String()
	clientCerts := conn.ConnectionState().PeerCertificates
	for _, clientCert := range clientCerts {
		for _, ipAddr := range clientCert.IPAddresses {
			if strings.HasPrefix(clientAddr, ipAddr.String()) {
				return nil
			}
		}
	}
	return errors.New("failed to validate client IP")
}

func handleConnection(conn *tls.Conn) {
	c := atomic.AddInt32(&counter, 1)

	err := validateClientIP(conn)
	if err != nil {
		log.Printf("%s\n", err)
		conn.Close()
		return
	}

	msg := strings.NewReader(fmt.Sprintf(`hello from server,
	you are coming from %s
	you are connection %d
`, conn.RemoteAddr(), c))
	err = lib.CopyBoth(conn, msg, os.Stdout)
	if err != nil {
		log.Fatalf("copy: %s", err)
	}
}
