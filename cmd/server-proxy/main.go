package main

import (
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/rosenhouse/tls-tunnel-experiments/lib"
)

func main() {
	var tlsFlags lib.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "server-proxy", "proxy-ca")

	var tunnelListenAddr string
	var localServerAddr string
	flag.StringVar(&tunnelListenAddr, "tunnelListenAddr", "127.0.0.22:7022", "listen for incoming tunnel connections from the client proxy")
	flag.StringVar(&localServerAddr, "localServerAddr", "127.0.0.12:7012", "local address of the end server")

	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	listener, err := tls.Listen("tcp", tunnelListenAddr, tlsConfig)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %s", err)
		}
		go handleConnection(conn.(*tls.Conn), localServerAddr)
	}
}

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

func handleConnection(tunnelConn *tls.Conn, localServerAddr string) error {
	defer tunnelConn.Close()

	err := validateClientIP(tunnelConn)
	if err != nil {
		return err
	}

	localConn, err := net.Dial("tcp", localServerAddr)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	defer localConn.Close()

	err = lib.CopyBoth(localConn, tunnelConn, tunnelConn)
	if err != nil {
		return fmt.Errorf("connection copy: %s", err)
	}
	return nil
}
