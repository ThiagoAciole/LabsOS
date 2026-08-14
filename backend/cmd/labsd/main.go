package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"labsos/backend/labsd"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	if err := labsd.New("", nil).ListenAndServe(ctx); err != nil { log.Fatal(err) }
}
