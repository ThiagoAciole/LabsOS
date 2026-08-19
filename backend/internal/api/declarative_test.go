package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeclarativeStorePersistsValidatedManifest(t *testing.T) {
	store := &declarativeStore{path: filepath.Join(t.TempDir(), "apps.json")}
	item, err := store.add(declarativeImportRequest{ID: "my-app", Name: "My App", Compose: "services:\n  app:\n    image: example/app"})
	if err != nil || item.ID != "my-app" {
		t.Fatalf("add = %#v, %v", item, err)
	}
	data := mustReadFile(t, store.path)
	if len(data) == 0 {
		t.Fatal("manifest was not persisted")
	}
}

func mustReadFile(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}
