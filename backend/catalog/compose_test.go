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

func TestComposeImporterAcceptsLongVolumeSyntax(t *testing.T) {
	path := filepath.Join(t.TempDir(), "compose.yml")
	content := "services:\n  app:\n    image: example/app\n    volumes:\n      - type: bind\n        source: /DATA/Apps/demo\n        target: /var/lib/demo\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ImportCompose(context.Background(), path); err != nil {
		t.Fatalf("long volume syntax rejected: %v", err)
	}
}

func TestComposeImporterAllowsCasaOSAppData(t *testing.T) {
	path := filepath.Join(t.TempDir(), "compose.yml")
	content := "services:\n  app:\n    image: example/app\n    volumes:\n      - /DATA/AppData/demo:/var/lib/demo\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ImportCompose(context.Background(), path); err != nil {
		t.Fatalf("CasaOS AppData path rejected: %v", err)
	}
}

func TestComposeImporterAcceptsCasaOSShape(t *testing.T) {
	path := filepath.Join(t.TempDir(), "compose.yml")
	content := "name: 2fauth\nservices:\n  2fauth:\n    image: 2fauth/2fauth:6.1.3\n    deploy:\n      resources:\n        reservations:\n          memory: 64M\n    network_mode: bridge\n    ports:\n      - target: 8000\n        published: '8000'\n        protocol: tcp\n    restart: always\n    volumes:\n      - type: bind\n        source: /DATA/AppData/$AppID\n        target: /2fauth\n    environment:\n      APP_KEY: SomeRandomStringOf32CharsExactly\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := ImportCompose(context.Background(), path); err != nil {
		t.Fatalf("CasaOS compose rejected: %v", err)
	}
}
