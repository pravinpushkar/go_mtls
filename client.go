package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	cert, err := os.ReadFile(filepath.Join("certs", "ca.pem"))
	if err != nil {
		log.Fatalf("could not open ca ertificate file: %v", err)
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(cert)

	clientCert := filepath.Join("certs", "client.pem")
	clientKey := filepath.Join("certs", "client.key")
	certificate, err := tls.LoadX509KeyPair(clientCert, clientKey)
	if err != nil {
		log.Fatalf("could not load certificate and key file: %v", err)
	}

	client := http.Client{
		Timeout: time.Minute * 3,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      caCertPool,
				Certificates: []tls.Certificate{certificate},
			},
		},
	}

	resp, err := client.Get("https://127.0.0.1:8443/hello")
	if err != nil {
		log.Fatalf("error making get request: %v", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v", err)
	}
	fmt.Printf("%s\n", body)
}
