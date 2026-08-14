package catalog

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComposeImporterNormalizesLocalCompose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "compose.yml")
	if err := os.WriteFile(path, []byte("services:\n  web:\n    image: nginx:alpine\n    ports:\n      - 8080:80\n"), 0600); err != nil {
		t.Fatal(err)
	}
	app, err := ImportCompose(context.Background(), path)
	if err != nil {
		t.Fatal(err)
	}
	if app.ID == "" || app.Name == "" || !strings.Contains(app.Compose, "nginx:alpine") {
		t.Fatalf("unexpected app: %+v", app)
	}
}

func TestComposeImporterRejectsUnsafeCompose(t *testing.T) {
	path := filepath.Join(t.TempDir(), "docker-compose.yml")
	if err := os.WriteFile(path, []byte("services:\n  web:\n    image: nginx\n    privileged: true\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ImportCompose(context.Background(), path); err == nil {
		t.Fatal("unsafe compose accepted")
	}
}
