package tlsconfig

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

func LoadTLSConfig(
	serverCertPath string,
	serverKeyPath string,
	caPath string,
) (*tls.Config, error) {

	serverCert, err := tls.LoadX509KeyPair(
		serverCertPath,
		serverKeyPath,
	)
	if err != nil {
		return nil, err
	}

	caPem, err := os.ReadFile(caPath)
	if err != nil {
		return nil, err
	}

	clientCAPool := x509.NewCertPool()

	if !clientCAPool.AppendCertsFromPEM(caPem) {
		return nil, fmt.Errorf("failed append ca")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{
			serverCert,
		},

		ClientAuth: tls.RequireAndVerifyClientCert,

		ClientCAs:  clientCAPool,
		ServerName: "gateway-client",
		MinVersion: tls.VersionTLS13,
	}, nil
}
