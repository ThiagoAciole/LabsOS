package providers_test

import (
	"labsos/backend/providers"
	"testing"
)

func TestLinuxProviderIsAvailable(t *testing.T) {
	if providers.New() == nil {
		t.Fatal("provider is nil")
	}
}
