package api

import (
	"archive/tar"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"labsos/backend/internal/safety"
)

type backupRecord struct {
	ID        string    `json:"id"`
	Path      string    `json:"path"`
	Size      int64     `json:"size"`
	SHA256    string    `json:"sha256"`
	CreatedAt time.Time `json:"createdAt"`
	Simulated bool      `json:"simulated"`
}

func backupDir() string {
	if value := os.Getenv("LABSOS_BACKUP_DIR"); value != "" {
		return value
	}
	return "/tmp/labsos-backups"
}
func (s *server) backups(w http.ResponseWriter, _ *http.Request) {
	entries, _ := os.ReadDir(backupDir())
	result := []backupRecord{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tar") {
			continue
		}
		path := filepath.Join(backupDir(), entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}
		checksum, _ := sha256File(path)
		result = append(result, backupRecord{ID: strings.TrimSuffix(entry.Name(), ".tar"), Path: path, Size: info.Size(), SHA256: checksum, CreatedAt: info.ModTime()})
	}
	respond(w, result, nil, http.StatusOK)
}
func (s *server) verifyBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !safeBackupID(id) {
		writeError(w, http.StatusBadRequest, "INVALID_BACKUP_ID", "invalid backup id")
		return
	}
	path := filepath.Join(backupDir(), id+".tar")
	info, err := os.Stat(path)
	if err != nil || info.IsDir() {
		writeError(w, http.StatusNotFound, "BACKUP_NOT_FOUND", "backup not found")
		return
	}
	checksum, err := sha256File(path)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "BACKUP_VERIFY_FAILED", err.Error())
		return
	}
	respond(w, map[string]any{"id": id, "path": path, "size": info.Size(), "sha256": checksum, "integrity": "verified"}, nil, http.StatusOK)
}
func (s *server) restoreBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !safeBackupID(id) {
		writeError(w, http.StatusBadRequest, "INVALID_BACKUP_ID", "invalid backup id")
		return
	}
	path := filepath.Join(backupDir(), id+".tar")
	if _, err := os.Stat(path); err != nil {
		writeError(w, http.StatusNotFound, "BACKUP_NOT_FOUND", "backup not found")
		return
	}
	entries, err := inspectArchive(path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_BACKUP", err.Error())
		return
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("backup.restore", id, "planned", "real operations disabled")
		respond(w, map[string]any{"id": id, "status": "planned", "simulated": true, "entries": entries, "message": "restore planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	if err := extractArchive(r.Context(), path, "/"); err != nil {
		writeError(w, http.StatusInternalServerError, "RESTORE_FAILED", err.Error())
		return
	}
	respond(w, map[string]any{"id": id, "status": "success", "simulated": false, "entries": entries}, nil, http.StatusAccepted)
	s.auditStore.add("backup.restore", id, "success", "")
}
func (s *server) deleteBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !safeBackupID(id) {
		writeError(w, http.StatusBadRequest, "INVALID_BACKUP_ID", "invalid backup id")
		return
	}
	path := filepath.Join(backupDir(), id+".tar")
	if _, err := os.Stat(path); err != nil {
		writeError(w, http.StatusNotFound, "BACKUP_NOT_FOUND", "backup not found")
		return
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("backup.delete", id, "planned", "real operations disabled")
		respond(w, map[string]any{"id": id, "status": "planned", "simulated": true, "message": "delete planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	if err := os.Remove(path); err != nil {
		writeError(w, http.StatusInternalServerError, "BACKUP_DELETE_FAILED", err.Error())
		return
	}
	respond(w, map[string]any{"id": id, "status": "deleted", "simulated": false}, nil, http.StatusOK)
	s.auditStore.add("backup.delete", id, "success", "")
}
func (s *server) createBackup(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Apps []string `json:"apps"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "backup request must be JSON")
		return
	}
	id := fmt.Sprintf("backup-%d", time.Now().UnixNano())
	paths := []string{"/DATA"}
	for _, app := range input.Apps {
		if safeBackupID(app) {
			paths = append(paths, filepath.Join("/opt/labsos/apps", app))
		}
	}
	if !safety.RealOperationsEnabled() {
		s.auditStore.add("backup.create", id, "planned", "real operations disabled")
		respond(w, map[string]any{"id": id, "status": "planned", "simulated": true, "paths": paths, "message": "backup is planned; real operations are disabled"}, nil, http.StatusAccepted)
		return
	}
	if err := os.MkdirAll(backupDir(), 0750); err != nil {
		writeError(w, http.StatusInternalServerError, "BACKUP_DIRECTORY_FAILED", err.Error())
		return
	}
	path := filepath.Join(backupDir(), id+".tar")
	if err := createTar(r.Context(), path, paths); err != nil {
		writeError(w, http.StatusInternalServerError, "BACKUP_FAILED", err.Error())
		return
	}
	info, _ := os.Stat(path)
	respond(w, backupRecord{ID: id, Path: path, Size: info.Size(), CreatedAt: time.Now().UTC()}, nil, http.StatusAccepted)
	s.auditStore.add("backup.create", id, "success", path)
}
func safeBackupID(id string) bool {
	return id != "" && id != "." && id != ".." && !strings.Contains(id, "..") && !strings.ContainsAny(id, "/\\")
}
func createTar(ctx context.Context, target string, paths []string) error {
	file, err := os.Create(target)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := tar.NewWriter(file)
	defer writer.Close()
	for _, root := range paths {
		if err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
			if walkErr != nil || info == nil {
				return walkErr
			}
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}
			header, err := tar.FileInfoHeader(info, "")
			if err != nil {
				return err
			}
			header.Name = strings.TrimPrefix(path, "/")
			if err := writer.WriteHeader(header); err != nil {
				return err
			}
			if info.Mode().IsRegular() {
				source, err := os.Open(path)
				if err != nil {
					return err
				}
				_, copyErr := io.Copy(writer, source)
				source.Close()
				if copyErr != nil {
					return copyErr
				}
			}
			return nil
		}); err != nil {
			return err
		}
	}
	return nil
}
func sha256File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := sha256.New()
	_, err = io.Copy(hash, file)
	return hex.EncodeToString(hash.Sum(nil)), err
}

func inspectArchive(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := tar.NewReader(file)
	entries := []string{}
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}
		if !safeArchivePath(header.Name) {
			return nil, fmt.Errorf("archive contains unsafe path %q", header.Name)
		}
		entries = append(entries, header.Name)
	}
}

func extractArchive(ctx context.Context, path, destination string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	reader := tar.NewReader(file)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if !safeArchivePath(header.Name) {
			return fmt.Errorf("archive contains unsafe path %q", header.Name)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		target := filepath.Join(destination, filepath.Clean(header.Name))
		if err := rejectSymlinkPath(destination, target); err != nil {
			return err
		}
		if header.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0750); err != nil {
				return err
			}
			continue
		}
		if !header.FileInfo().Mode().IsRegular() {
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0750); err != nil {
			return err
		}
		output, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, header.FileInfo().Mode().Perm())
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(output, reader)
		closeErr := output.Close()
		if copyErr != nil {
			return copyErr
		}
		if closeErr != nil {
			return closeErr
		}
	}
}

func rejectSymlinkPath(destination, target string) error {
	root, err := filepath.Abs(destination)
	if err != nil {
		return err
	}
	cleanTarget, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(root, cleanTarget)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		return fmt.Errorf("archive target escapes destination")
	}
	current := root
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		if part == "." || part == "" {
			continue
		}
		current = filepath.Join(current, part)
		info, statErr := os.Lstat(current)
		if statErr != nil {
			if os.IsNotExist(statErr) {
				continue
			}
			return statErr
		}
		if info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("archive target crosses a symbolic link: %q", current)
		}
	}
	return nil
}

func safeArchivePath(name string) bool {
	clean := filepath.Clean(name)
	return clean != "." && !filepath.IsAbs(name) && clean != ".." && !strings.HasPrefix(clean, "../") && (clean == "DATA" || strings.HasPrefix(clean, "DATA/") || clean == "opt/labsos/apps" || strings.HasPrefix(clean, "opt/labsos/apps/"))
}
