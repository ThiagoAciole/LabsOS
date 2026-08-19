package api

import "testing"

func TestValidPublicKey(t *testing.T) {
	if !validPublicKey("ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI example") {
		t.Fatal("expected valid key")
	}
	if validPublicKey("not-a-key") || validPublicKey("ssh-ed25519 bad\nkey") {
		t.Fatal("expected invalid key")
	}
}
