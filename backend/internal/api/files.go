package api

import (
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type fileEntry struct {
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	Type       string    `json:"type"`
	Size       int64     `json:"size"`
	ModifiedAt time.Time `json:"modifiedAt"`
}

func filesRoot() string {
	if value := os.Getenv("LABSOS_FILES_ROOT"); value != "" {
		return filepath.Clean(value)
	}
	return "/DATA"
}
func (s *server) files(w http.ResponseWriter, r *http.Request) {
	root := filesRoot()
	requested := strings.TrimSpace(r.URL.Query().Get("path"))
	if requested == "" {
		requested = "."
	}
	clean := filepath.Clean(requested)
	if filepath.IsAbs(clean) {
		clean = strings.TrimPrefix(clean, string(filepath.Separator))
	}
	target := filepath.Join(root, clean)
	relative, err := filepath.Rel(root, target)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
		writeError(w, http.StatusBadRequest, "INVALID_PATH", "path is outside the files root")
		return
	}
	entries, err := os.ReadDir(target)
	if err != nil {
		if os.IsNotExist(err) {
			writeError(w, http.StatusNotFound, "NOT_FOUND", "directory not found")
			return
		}
		writeError(w, http.StatusForbidden, "FILES_UNAVAILABLE", "directory cannot be read")
		return
	}
	result := make([]fileEntry, 0, len(entries))
	for _, entry := range entries {
		info, infoErr := entry.Info()
		if infoErr != nil {
			continue
		}
		itemPath := filepath.ToSlash(filepath.Join(relative, entry.Name()))
		if itemPath == "." {
			itemPath = entry.Name()
		}
		kind := "file"
		if entry.IsDir() {
			kind = "directory"
		}
		result = append(result, fileEntry{Name: entry.Name(), Path: itemPath, Type: kind, Size: info.Size(), ModifiedAt: info.ModTime().UTC()})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Type != result[j].Type {
			return result[i].Type == "directory"
		}
		return strings.ToLower(result[i].Name) < strings.ToLower(result[j].Name)
	})
	respond(w, map[string]any{"root": root, "path": filepath.ToSlash(relative), "entries": result, "readOnly": true}, nil, http.StatusOK)
}
