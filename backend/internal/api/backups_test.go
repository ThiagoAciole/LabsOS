package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSafeArchivePath(t *testing.T) {
	for _, path := range []string{"DATA", "DATA/config/app.json", "opt/labsos/apps/jellyfin/compose.yaml"} {
		if !safeArchivePath(path) {
			t.Errorf("expected safe archive path: %q", path)
		}
	}
	for _, path := range []string{"../../etc/passwd", "/etc/passwd", "var/lib/labsos/secret", "opt/labsos/apps/../../etc/passwd"} {
		if safeArchivePath(path) {
			t.Errorf("expected unsafe archive path: %q", path)
		}
	}
}

func TestSafeBackupID(t *testing.T) {
	for _, id := range []string{"backup-1", "app_data"} {
		if !safeBackupID(id) {
			t.Errorf("expected safe id: %q", id)
		}
	}
	for _, id := range []string{"", "..", "../backup", "a/b", `a\\b`} {
		if safeBackupID(id) {
			t.Errorf("expected unsafe id: %q", id)
		}
	}
}

func TestRejectSymlinkPath(t *testing.T) {
	root := t.TempDir()
	outside := t.TempDir()
	link := filepath.Join(root, "DATA")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := rejectSymlinkPath(root, filepath.Join(root, "DATA", "secret")); err == nil {
		t.Fatal("expected symlink traversal to be rejected")
	}
}
