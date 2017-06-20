package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/rosenhouse/tls-tunnel-experiments/lib"
)

func main() {
	var tlsFlags lib.MTLSFlags
	tlsFlags.AddFlags(flag.CommandLine, "client-proxy", "proxy-ca")

	var localListenAddr string
	var remoteTunnelAddr string
	flag.StringVar(&localListenAddr, "localProxyListenAddr", "127.0.0.21:7021", "local listen addr for the client proxy")
	flag.StringVar(&remoteTunnelAddr, "remoteTunnelAddr", "127.0.0.22:7022", "address of tunnel listener on server-side proxy")

	flag.Parse()

	tlsConfig, err := tlsFlags.LoadConfig()
	if err != nil {
		log.Fatalf("load tls config: %s", err)
	}

	listener, err := net.Listen("tcp", localListenAddr)
	if err != nil {
		log.Fatalf("listen: %s", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept conn: %s", err)
		}
		go func() {
			err := handleConnection(conn, tlsConfig, localListenAddr, remoteTunnelAddr)
			if err != nil {
				log.Printf("handle conn: %s", err)
			}
		}()
	}
}

func handleConnection(localConn net.Conn, tlsConfig *tls.Config, localListenAddr, remoteTunnelAddr string) error {
	defer localConn.Close()

	localAddr, err := net.ResolveTCPAddr("tcp", localListenAddr)
	if err != nil {
		return fmt.Errorf("parsing client ip as tcp addr: %s", err)
	}
	localAddr.Port = 0

	proxyDialer := &net.Dialer{LocalAddr: localAddr}
	proxyConn, err := tls.DialWithDialer(proxyDialer, "tcp", remoteTunnelAddr, tlsConfig)
	if err != nil {
		return fmt.Errorf("dial: %s", err)
	}

	defer proxyConn.Close()

	err = lib.CopyBoth(proxyConn, localConn, localConn)
	if err != nil {
		return fmt.Errorf("connection copy: %s", err)
	}

	return nil
}
