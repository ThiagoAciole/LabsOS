package api

import (
	"path/filepath"
	"testing"
)

func TestNotificationStoreDeletePersists(t *testing.T) {
	store := &notificationStore{path: filepath.Join(t.TempDir(), "notifications.json")}
	item := store.add("Title", "Message", "test", "info")
	if !store.delete(item.ID) || len(store.list()) != 0 {
		t.Fatal("notification was not deleted")
	}
}
