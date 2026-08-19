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

func TestExtractRejectsSymlinkAndLinkEntries(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(t.TempDir(), "outside")
	if err := os.WriteFile(outside, []byte("safe"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, "escape")); err != nil {
		t.Fatal(err)
	}
	archive := makeArchiveWithEntry(t, &tar.Header{Name: "escape/pwned", Mode: 0600, Size: 4}, []byte("oops"))
	if err := extract(archive, root); err == nil {
		t.Fatal("expected extraction through symlink to be rejected")
	}
	if got, err := os.ReadFile(outside); err != nil || string(got) != "safe" {
		t.Fatalf("outside file changed: %q, %v", got, err)
	}

	archive = makeArchiveWithEntry(t, &tar.Header{Name: "link", Typeflag: tar.TypeSymlink, Linkname: outside}, nil)
	if err := extract(archive, t.TempDir()); err == nil {
		t.Fatal("expected symlink archive entry to be rejected")
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

func makeArchiveWithEntry(t *testing.T, header *tar.Header, data []byte) []byte {
	buffer := new(bytes.Buffer)
	gz := gzip.NewWriter(buffer)
	tw := tar.NewWriter(gz)
	header.Size = int64(len(data))
	if err := tw.WriteHeader(header); err != nil {
		t.Fatal(err)
	}
	if len(data) > 0 {
		if _, err := tw.Write(data); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gz.Close(); err != nil {
		t.Fatal(err)
	}
	return buffer.Bytes()
}
