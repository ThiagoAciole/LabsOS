package catalog

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadWithCacheKeepsLastValidCatalog(t *testing.T) {
	path := filepath.Join(t.TempDir(), "apps.json")
	cache := FileCache{Path: path}
	valid := []App{{ID: "jellyfin", Name: "Jellyfin"}}
	if err := cache.Save(valid); err != nil {
		t.Fatal(err)
	}
	apps, err := LoadWithCache(context.Background(), func(context.Context) ([]App, error) { return []App{}, nil }, cache)
	if err != nil || len(apps) != 1 || apps[0].ID != "jellyfin" {
		t.Fatalf("apps = %+v, err = %v", apps, err)
	}
	data, err := os.ReadFile(path)
	if err != nil || len(data) == 0 {
		t.Fatalf("cache was lost: %v", err)
	}
}

func TestNormalizeRemoteApp(t *testing.T) {
	app, err := Normalize(RemoteApp{ID: "jellyfin", Name: "Jellyfin", Description: "Media server", Icon: "https://example/icon.png", Category: "media", Version: "10.10"}, "casaos")
	if err != nil {
		t.Fatal(err)
	}
	if app.ID != "jellyfin" || app.Source != "casaos" || app.Installed {
		t.Fatalf("normalized app = %+v", app)
	}
}

func TestNormalizeRejectsMissingIDOrName(t *testing.T) {
	if _, err := Normalize(RemoteApp{Name: "Jellyfin"}, "casaos"); err == nil {
		t.Fatal("missing ID accepted")
	}
	if _, err := Normalize(RemoteApp{ID: "jellyfin"}, "casaos"); err == nil {
		t.Fatal("missing name accepted")
	}
}

func TestBuiltInProviderListsNormalizedApps(t *testing.T) {
	provider := NewBuiltInProvider([]RemoteApp{{ID: "jellyfin", Name: "Jellyfin"}})
	apps, err := provider.ListApps()
	if err != nil || len(apps) != 1 || apps[0].Source != "builtin" {
		t.Fatalf("apps = %+v, err = %v", apps, err)
	}
}

func TestRemoteProviderParsesAndSkipsInvalidApps(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`[{"id":"jellyfin","name":"Jellyfin"},{"name":"invalid"}]`))
	}))
	defer server.Close()

	apps, err := (RemoteProvider{URL: server.URL}).ListApps(context.Background())
	if err != nil || len(apps) != 1 || apps[0].ID != "jellyfin" {
		t.Fatalf("apps = %+v, err = %v", apps, err)
	}
}
