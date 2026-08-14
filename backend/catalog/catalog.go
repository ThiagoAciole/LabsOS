package catalog

import (
	"errors"
	"strings"
)

var ErrNotFound = errors.New("catalog app not found")

type RemoteApp struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Icon        string   `json:"icon"`
	Category    string   `json:"category"`
	Version     string   `json:"version"`
	Title       string   `json:"title"`
	Tagline     string   `json:"tagline"`
	Categories  []string `json:"categories"`
}

type App struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	Category    string `json:"category"`
	Version     string `json:"version,omitempty"`
	Source      string `json:"source"`
	Installed   bool   `json:"installed"`
	Installable bool   `json:"installable"`
	Compose     string `json:"compose,omitempty"`
}

func Normalize(remote RemoteApp, source string) (App, error) {
	name := remote.Name
	if name == "" {
		name = remote.Title
	}
	if remote.Description == "" {
		remote.Description = remote.Tagline
	}
	category := strings.ToLower(remote.Category)
	if category == "" && len(remote.Categories) > 0 {
		category = strings.ToLower(remote.Categories[0])
	}
	if strings.TrimSpace(remote.ID) == "" || strings.TrimSpace(name) == "" {
		return App{}, errors.New("catalog app requires id and name")
	}
	id := remote.ID
	switch strings.ToLower(name) {
	case "jellyfin":
		id = "jellyfin"
	case "syncthing":
		id = "syncthing"
	}
	return App{ID: id, Name: name, Description: remote.Description, Icon: remote.Icon, Category: category, Version: remote.Version, Source: source, Installable: id == "jellyfin" || id == "syncthing"}, nil
}
