package catalog

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type RemoteProvider struct {
	URL    string
	Client *http.Client
}

type remoteIndex struct {
	BaseURL string      `json:"base_url"`
	Apps    []RemoteApp `json:"apps"`
}

func (p RemoteProvider) ListApps(ctx context.Context) ([]App, error) {
	client := p.Client
	if client == nil {
		client = http.DefaultClient
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, p.URL, nil)
	if err != nil {
		return nil, err
	}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("catalog source returned %s", response.Status)
	}
	var payload json.RawMessage
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	var remote []RemoteApp
	baseURL := ""
	if len(payload) > 0 && payload[0] == '[' {
		if err := json.Unmarshal(payload, &remote); err != nil {
			return nil, err
		}
	} else {
		var index remoteIndex
		if err := json.Unmarshal(payload, &index); err != nil {
			return nil, err
		}
		remote = index.Apps
		baseURL = index.BaseURL
	}
	if baseURL != "" {
		base, err := url.Parse(baseURL)
		if err != nil {
			return nil, err
		}
		for i := range remote {
			if strings.HasPrefix(remote[i].Icon, "/") {
				if icon, err := base.Parse(remote[i].Icon); err == nil {
					remote[i].Icon = icon.String()
				}
			}
			if strings.HasPrefix(remote[i].ComposeURL, "/") {
				if compose, err := base.Parse(remote[i].ComposeURL); err == nil {
					remote[i].ComposeURL = compose.String()
				}
			}
		}
	}
	apps := make([]App, 0, len(remote))
	for _, item := range remote {
		app, err := Normalize(item, "remote")
		if err != nil {
			continue
		}
		apps = append(apps, app)
	}
	return apps, nil
}

func (p RemoteProvider) GetApp(id string) (App, error) {
	ctx := context.Background()
	apps, err := p.ListApps(ctx)
	if err != nil {
		return App{}, err
	}
	for _, app := range apps {
		if app.ID == id {
			if app.ComposeURL == "" {
				return app, nil
			}
			compose, err := ImportCompose(ctx, app.ComposeURL)
			if err != nil {
				return App{}, err
			}
			compose.ID, compose.Name, compose.Description = app.ID, app.Name, app.Description
			compose.Icon, compose.Category, compose.Version, compose.Source = app.Icon, app.Category, app.Version, app.Source
			return compose, nil
		}
	}
	return App{}, fmt.Errorf("catalog app %q not found", id)
}
