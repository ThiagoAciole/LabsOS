package main

import "testing"

func TestListenAddressUsesEnvironmentWithStableDefault(t *testing.T) {
	t.Setenv("LABSOS_ADDR", "")
	if got := listenAddress(); got != "127.0.0.1:8080" {
		t.Fatalf("default address = %q", got)
	}
	t.Setenv("LABSOS_ADDR", "127.0.0.1:18080")
	if got := listenAddress(); got != "127.0.0.1:18080" {
		t.Fatalf("configured address = %q", got)
	}
}
