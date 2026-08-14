package providers_test

import (
	"testing"

	"labsos/backend/providers"
)

func TestNewSelectsKnownModesAndRejectsUnknownMode(t *testing.T) {
	for _, mode := range []string{"mock", "linux"} {
		provider, err := providers.New(mode)
		if err != nil || provider.Mode() != mode {
			t.Fatalf("mode %q: provider=%v err=%v", mode, provider, err)
		}
	}
	if _, err := providers.New("invalid"); err == nil {
		t.Fatal("unknown mode was accepted")
	}
}
