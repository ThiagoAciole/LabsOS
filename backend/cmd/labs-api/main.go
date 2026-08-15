package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"labsos/backend/internal/api"
	"labsos/backend/providers"
)

func main() {
	mode := os.Getenv("LABSOS_MODE")
	if mode == "" {
		mode = "mock"
	}
	provider, err := providers.New(mode)
	if err != nil {
		log.Fatal(err)
	}
	addr := listenAddress()
	server := &http.Server{
		Addr:              addr,
		Handler:           api.New(provider),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       2 * time.Minute,
		WriteTimeout:      2 * time.Minute,
		IdleTimeout:       60 * time.Second,
	}
	log.Printf("Labs API http://%s Mode %s", addr, provider.Mode())
	log.Fatal(server.ListenAndServe())
}

func listenAddress() string {
	if addr := os.Getenv("LABSOS_ADDR"); addr != "" {
		return addr
	}
	return "127.0.0.1:8080"
}
