package api

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilesRootAndTraversal(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "DATA"), 0750); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LABSOS_FILES_ROOT", filepath.Join(root, "DATA"))
	if filesRoot() != filepath.Join(root, "DATA") {
		t.Fatal("unexpected files root")
	}
}
