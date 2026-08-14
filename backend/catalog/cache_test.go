package catalog

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
)

func TestCatalogCacheRoundTrip(t *testing.T) {
	cache := FileCache{Path: filepath.Join(t.TempDir(), "catalog.json")}
	want := []App{{ID: "jellyfin", Name: "Jellyfin", Source: "remote"}}
	if err := cache.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := cache.Load()
	if err != nil || len(got) != 1 || got[0].ID != want[0].ID {
		t.Fatalf("got = %+v, err = %v", got, err)
	}
}

func TestLoadWithCacheUsesLastValidCatalogWhenRemoteFails(t *testing.T) {
	cache := FileCache{Path: filepath.Join(t.TempDir(), "catalog.json")}
	want := []App{{ID: "jellyfin", Name: "Jellyfin"}}
	if err := cache.Save(want); err != nil {
		t.Fatal(err)
	}
	got, err := LoadWithCache(context.Background(), func(context.Context) ([]App, error) { return nil, errors.New("offline") }, cache)
	if err != nil || len(got) != 1 || got[0].ID != want[0].ID {
		t.Fatalf("got = %+v, err = %v", got, err)
	}
}
