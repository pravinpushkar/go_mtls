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
)

func main() {
	sslPort := 8443

	handler := http.NewServeMux()
	handler.HandleFunc("/hello", helloHandler)
	handler.HandleFunc("/bye", byeHandler)
	startHTTPSServer(sslPort, handler)
}

func startHTTPSServer(sslPort int, handler http.Handler) {
	caCertFile, err := os.ReadFile(filepath.Join("certs", "ca.pem"))
	if err != nil {
		log.Fatal("Failed to read ca certificate ca.pem")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertFile)

	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}
	tlsConfig.BuildNameToCertificate()
	server := http.Server{
		Addr:      fmt.Sprintf(":%d", sslPort),
		Handler:   handler,
		TLSConfig: tlsConfig,
	}

	fmt.Println("HTTPS server starting on port ", sslPort)
	if err := server.ListenAndServeTLS(filepath.Join("certs", "server.pem"), filepath.Join("certs", "server.key")); err != nil {
		log.Fatal("HTTPS server err on starting %w ", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		fmt.Println("TLS is enabled")
		printTLSConnectionInfo(r.TLS)
	}
	io.WriteString(w, "Hello, world!")
}

func byeHandler(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		fmt.Println("TLS is enabled")
		printTLSConnectionInfo(r.TLS)
	}
	io.WriteString(w, "Bye bye Hello, world!")
}

func printTLSConnectionInfo(tlsState *tls.ConnectionState) {
	fmt.Println(">>>>>>>>>>>>>>>>> TLS connectionState <<<<<<<<<<<<<<<<<<")
	fmt.Println("ServerName: ", tlsState.ServerName)
	log.Printf("HandshakeComplete: %t", tlsState.HandshakeComplete)

	fmt.Println("Certificate data:")
	for _, cert := range tlsState.PeerCertificates {
		fmt.Println("Subject: ", cert.Subject)
		fmt.Println("Issuer: ", cert.Issuer)
	}
}
