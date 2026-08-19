package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"labsos/backend/internal/api"
	"labsos/backend/internal/safety"
	"labsos/backend/providers"
)

func main() {
	provider := providers.New()
	addr := listenAddress()
	server := &http.Server{
		Addr:              safety.Address(addr),
		Handler:           api.New(provider),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       2 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Labs API http://%s (real operations enabled: %t)", safety.Address(addr), safety.RealOperationsEnabled())
	log.Fatal(server.ListenAndServe())
}

func listenAddress() string {
	if addr := os.Getenv("LABSOS_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8080"
}
