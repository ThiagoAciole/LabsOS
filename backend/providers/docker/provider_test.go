package docker

import (
	"context"
	"testing"
)

type unavailableRuntime struct{}

func (unavailableRuntime) Available(context.Context) (bool, error) { return false, nil }
func (unavailableRuntime) Version(context.Context) (string, error) { return "", nil }

func TestProviderFailsClosedWhenDockerIsUnavailable(t *testing.T) {
	provider := New(unavailableRuntime{})
	if _, err := provider.List(context.Background()); err == nil {
		t.Fatal("list succeeded without Docker")
	}
	if _, err := provider.Install(context.Background(), "jellyfin"); err == nil {
		t.Fatal("install succeeded without Docker")
	}
}
