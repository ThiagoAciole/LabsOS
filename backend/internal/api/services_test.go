package api

import (
	"path/filepath"
	"testing"

	"labsos/backend/internal/platform"
)

func TestServicePortConflicts(t *testing.T) {
	conflicts := servicePortConflicts([]platform.App{
		{ID: "one", Ports: []int{8080, 8443}},
		{ID: "two", Ports: []int{8080}},
		{ID: "three", Ports: []int{9090}},
	})
	if !conflicts["one"] || !conflicts["two"] {
		t.Fatalf("expected shared port owners to conflict: %#v", conflicts)
	}
	if conflicts["three"] {
		t.Fatalf("unexpected conflict for isolated port: %#v", conflicts)
	}
}

func TestExposureStorePersistsState(t *testing.T) {
	t.Setenv("LABSOS_EXPOSURES_FILE", filepath.Join(t.TempDir(), "exposures.json"))
	first := newExposureStore()
	first.set("jellyfin", true)
	second := newExposureStore()
	value, ok := second.get("jellyfin")
	if !ok || !value {
		t.Fatalf("exposure state was not restored: %v %v", value, ok)
	}
}
