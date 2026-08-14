package runtime

import (
	"context"
	"errors"
	"testing"
)

func TestDockerRuntimeReportsAvailabilityAndVersion(t *testing.T) {
	runtime := NewDockerRuntime(func(context.Context, ...string) ([]byte, error) {
		return []byte(`{"Client":{"Version":"28.0.0"}}`), nil
	})

	available, err := runtime.Available(context.Background())
	if err != nil || !available {
		t.Fatalf("available = %v, %v", available, err)
	}
	version, err := runtime.Version(context.Background())
	if err != nil || version != "28.0.0" {
		t.Fatalf("version = %q, %v", version, err)
	}
}

func TestDockerRuntimeTreatsMissingDockerAsUnavailable(t *testing.T) {
	runtime := NewDockerRuntime(func(context.Context, ...string) ([]byte, error) {
		return nil, errors.New("executable file not found")
	})

	available, err := runtime.Available(context.Background())
	if err != nil || available {
		t.Fatalf("available = %v, %v", available, err)
	}
}
