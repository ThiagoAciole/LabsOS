package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAuditStorePersistsEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.json")
	t.Setenv("LABSOS_AUDIT_FILE", path)
	store := newAuditStore()
	store.add("test", "target", "success", "details")
	loaded := newAuditStore()
	entries := loaded.list()
	if len(entries) != 1 || entries[0].Action != "test" || entries[0].Status != "success" {
		t.Fatalf("unexpected entries: %#v", entries)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatal(err)
	}
}
