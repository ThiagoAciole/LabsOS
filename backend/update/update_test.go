package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestManagerDownloadsChecksAndActivatesRelease(t *testing.T) {
	archive := makeArchive(t, "0.1.1")
	sum := fmt.Sprintf("%x", sha256.Sum256(archive))
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/manifest.json" {
			fmt.Fprintf(w, `{"version":"0.1.1","artifactUrl":"%s/artifact.tgz","sha256":"%s"}`, "http://"+r.Host, sum)
			return
		}
		w.Write(archive)
	}))
	defer server.Close()
	manager := Manager{Root: t.TempDir(), Client: server.Client(), ManifestURL: server.URL + "/manifest.json"}
	status, err := manager.Update(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if status.CurrentVersion != "0.1.1" {
		t.Fatalf("version = %s", status.CurrentVersion)
	}
	if _, err := os.Stat(filepath.Join(manager.Root, "current", "version.json")); err != nil {
		t.Fatal(err)
	}
}

func makeArchive(t *testing.T, version string) []byte {
	var output []byte
	buffer := new(bytes.Buffer)
	gz := gzip.NewWriter(buffer)
	tw := tar.NewWriter(gz)
	data := []byte(`{"version":"` + version + `"}`)
	if err := tw.WriteHeader(&tar.Header{Name: "version.json", Mode: 0600, Size: int64(len(data))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(data); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()
	output = buffer.Bytes()
	return output
}
