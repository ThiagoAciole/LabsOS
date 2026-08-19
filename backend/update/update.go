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
	Channel     string `json:"channel,omitempty"`
	Changelog   string `json:"changelog,omitempty"`
}
type Status struct {
	CurrentVersion  string `json:"currentVersion"`
	LatestVersion   string `json:"latestVersion"`
	UpdateAvailable bool   `json:"updateAvailable"`
	Channel         string `json:"channel,omitempty"`
	Changelog       string `json:"changelog,omitempty"`
	SHA256          string `json:"sha256,omitempty"`
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
	status.Channel, status.Changelog, status.SHA256 = manifest.Channel, manifest.Changelog, manifest.SHA256
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

func (m Manager) Rollback() (Status, error) {
	previous := filepath.Join(m.Root, ".current.previous")
	target, err := filepath.EvalSymlinks(previous)
	if err != nil {
		return Status{}, fmt.Errorf("previous release unavailable: %w", err)
	}
	if err := activate(m.Root, target); err != nil {
		return Status{}, err
	}
	return Status{CurrentVersion: versionFromRelease(target), LatestVersion: versionFromRelease(target)}, nil
}

func versionFromRelease(path string) string {
	data, err := os.ReadFile(filepath.Join(path, "version.json"))
	if err != nil {
		return filepath.Base(path)
	}
	var manifest Manifest
	if json.Unmarshal(data, &manifest) == nil && manifest.Version != "" {
		return manifest.Version
	}
	return filepath.Base(path)
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
		name := filepath.Clean(header.Name)
		if filepath.IsAbs(name) || name == ".." || strings.HasPrefix(name, ".."+string(os.PathSeparator)) {
			return fmt.Errorf("archive path escapes release")
		}
		target := filepath.Join(root, name)
		if err := rejectSymlinkPath(root, target); err != nil {
			return err
		}
		if header.Typeflag == tar.TypeSymlink || header.Typeflag == tar.TypeLink {
			return fmt.Errorf("archive links are not allowed: %s", name)
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

func rejectSymlinkPath(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(os.PathSeparator)) {
		return fmt.Errorf("archive path escapes release")
	}
	current := root
	for _, part := range strings.Split(relative, string(os.PathSeparator)) {
		if part == "." || part == "" {
			continue
		}
		current = filepath.Join(current, part)
		info, statErr := os.Lstat(current)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				return nil
			}
			return statErr
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("archive path contains symlink: %s", current)
		}
	}
	return nil
}
func activate(root, release string) error {
	current := filepath.Join(root, "current")
	if runtime.GOOS != "windows" {
		temporary := filepath.Join(root, ".current.new")
		previous := filepath.Join(root, ".current.previous")
		_ = os.Remove(temporary)
		if err := os.Symlink(release, temporary); err != nil {
			return err
		}
		_ = os.Remove(previous)
		if _, err := os.Lstat(current); err == nil {
			if err := os.Rename(current, previous); err != nil {
				_ = os.Remove(temporary)
				return err
			}
		}
		if err := os.Rename(temporary, current); err != nil {
			_ = os.Rename(previous, current)
			return err
		}
		return nil
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
