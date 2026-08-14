package providers_test

import (
	"context"
	"testing"

	"labsos/backend/internal/platform"
	"labsos/backend/providers"
)

func TestMockAppsAreUserApps(t *testing.T) {
	provider, err := providers.New("mock")
	if err != nil {
		t.Fatal(err)
	}
	apps, err := provider.Apps(context.Background(), true)
	if err != nil {
		t.Fatal(err)
	}
	for _, app := range apps {
		if app.Kind != platform.AppKindUser {
			t.Fatalf("app %q kind = %q", app.ID, app.Kind)
		}
	}
}

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
