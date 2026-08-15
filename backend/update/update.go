package update

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"labsos/backend/internal/version"
)

type Manifest struct {
	Version     string `json:"version"`
	ArtifactURL string `json:"artifactUrl"`
	SHA256      string `json:"sha256"`
}
type Status struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
}
type Manager struct {
	Root        string
	ManifestURL string
	Client      *http.Client
}

func (m Manager) CurrentVersion() string {
	data, err := os.ReadFile(filepath.Join(m.Root, "current", "version.json"))
	if err != nil {
		return version.Value
	}
	var value Manifest
	if json.Unmarshal(data, &value) != nil || value.Version == "" {
		return version.Value
	}
	return value.Version
}
func (m Manager) Check(ctx context.Context) (Status, error) {
	current := m.CurrentVersion()
	status := Status{CurrentVersion: current, LatestVersion: current}
	if m.ManifestURL == "" {
		return status, nil
	}
	manifest, err := m.fetchManifest(ctx)
	if err != nil {
		return status, err
	}
	status.LatestVersion, status.UpdateAvailable = manifest.Version, manifest.Version != current
	return status, nil
}
func (m Manager) Update(ctx context.Context) (Status, error) {
	manifest, err := m.fetchManifest(ctx)
	if err != nil {
		return Status{}, err
	}
	client := m.Client
	if client == nil {
		client = http.DefaultClient
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, manifest.ArtifactURL, nil)
	if err != nil {
		return Status{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return Status{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Status{}, fmt.Errorf("artifact returned %s", response.Status)
	}
	archive, err := io.ReadAll(io.LimitReader(response.Body, 512<<20))
	if err != nil {
		return Status{}, err
	}
	sum := sha256.Sum256(archive)
	if !strings.EqualFold(hex.EncodeToString(sum[:]), manifest.SHA256) {
		return Status{}, fmt.Errorf("artifact checksum mismatch")
	}
	release := filepath.Join(m.Root, "releases", manifest.Version)
	if err := os.MkdirAll(release, 0750); err != nil {
		return Status{}, err
	}
	if err := extract(archive, release); err != nil {
		return Status{}, err
	}
	if err := activate(m.Root, release); err != nil {
		return Status{}, err
	}
	return Status{CurrentVersion: manifest.Version, LatestVersion: manifest.Version}, nil
}
func (m Manager) fetchManifest(ctx context.Context) (Manifest, error) {
	client := m.Client
	if client == nil {
		client = http.DefaultClient
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, m.ManifestURL, nil)
	if err != nil {
		return Manifest{}, err
	}
	response, err := client.Do(request)
	if err != nil {
		return Manifest{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return Manifest{}, fmt.Errorf("manifest returned %s", response.Status)
	}
	var value Manifest
	err = json.NewDecoder(response.Body).Decode(&value)
	if err == nil && (value.Version == "" || value.ArtifactURL == "" || value.SHA256 == "") {
		err = fmt.Errorf("manifest is incomplete")
	}
	return value, err
}
func extract(data []byte, root string) error {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer gz.Close()
	tarReader := tar.NewReader(gz)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target := filepath.Join(root, filepath.Clean("/"+header.Name))
		if !strings.HasPrefix(target, filepath.Clean(root)+string(os.PathSeparator)) {
			return fmt.Errorf("archive path escapes release")
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0750); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
			return err
		}
		file, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(file, io.LimitReader(tarReader, header.Size))
		closeErr := file.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
}
func activate(root, release string) error {
	current := filepath.Join(root, "current")
	_ = os.Remove(current)
	if runtime.GOOS != "windows" {
		return os.Symlink(release, current)
	}
	return copyDir(release, current)
}
func copyDir(source, target string) error {
	return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relative, _ := filepath.Rel(source, path)
		destination := filepath.Join(target, relative)
		if info.IsDir() {
			return os.MkdirAll(destination, 0750)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destination, data, info.Mode())
	})
}
