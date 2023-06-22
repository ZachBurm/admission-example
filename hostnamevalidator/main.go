package main

// credit to https://github.com/giantswarm/grumpy

import (
	"context"
	"crypto/tls"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	"hostnamevalidator/pkg/handler"
)

var (
	tlscert, tlskey string
)

func main() {

	flag.StringVar(&tlscert, "tlsCertFile", "/etc/certs/cert.pem", "File containing the x509 Certificate for HTTPS.")
	flag.StringVar(&tlskey, "tlsKeyFile", "/etc/certs/key.pem", "File containing the x509 private key to --tlsCertFile.")

	flag.Parse()

	certs, err := tls.LoadX509KeyPair(tlscert, tlskey)
	if err != nil {
		log.Err(err).Msg("failed to load key pair")
	}

	mux := http.NewServeMux()

	vh := handler.VerifyHandler{}
	mux.HandleFunc("/validate", vh.Serve)

	srv := &http.Server{
		Addr:    ":443",
		Handler: mux,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
			Certificates:       []tls.Certificate{certs},
		},
	}

	// start webhook server in new routine
	go func() {
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			log.Err(err).Msg("Failed to listen and serve webhook server")
		}
	}()

	log.Info().Msgf("Server running listening in port 443")

	// listening shutdown singal
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	<-signalChan

	log.Info().Msgf("Got shutdown signal, shutting down webhook server gracefully...")
	srv.Shutdown(context.Background())
}
