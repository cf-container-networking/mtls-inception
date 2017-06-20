package lib

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
)

type MTLSFlags struct {
	CAPath   string
	KeyPath  string
	CertPath string
}

// AddFlags will add command line flags for CA, Key and Cert to the given flagset
// The role will be used to name the default values and set the usage text
func (c *MTLSFlags) AddFlags(flagSet *flag.FlagSet, role, caName string) {
	flagSet.StringVar(&c.CAPath, "ca", "certs/"+caName+".crt", "path to CA certificate")
	flagSet.StringVar(&c.KeyPath, "key", "certs/"+role+".key", "path to "+role+" key")
	flagSet.StringVar(&c.CertPath, "cert", "certs/"+role+".crt", "path to "+role+" certificate")
}

func (c *MTLSFlags) LoadConfig() (*tls.Config, error) {
	keyPair, err := tls.LoadX509KeyPair(c.CertPath, c.KeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to load cert or key: %s", err)
	}

	caCertBytes, err := ioutil.ReadFile(c.CAPath)
	if err != nil {
		return nil, fmt.Errorf("failed read ca cert file: %s", err.Error())
	}

	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCertBytes); !ok {
		return nil, errors.New("unable to load ca cert")
	}

	tlsConfig := &tls.Config{
		Certificates:             []tls.Certificate{keyPair},
		MinVersion:               tls.VersionTLS12,
		RootCAs:                  caCertPool,
		ClientCAs:                caCertPool,
		ClientAuth:               tls.RequireAndVerifyClientCert,
		PreferServerCipherSuites: true,
		CipherSuites:             []uint16{tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256},
	}

	return tlsConfig, nil
}
