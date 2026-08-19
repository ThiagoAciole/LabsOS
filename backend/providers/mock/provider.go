package mock

import (
	"context"
	"fmt"
	"sync"

	"labsos/backend/catalog"
	"labsos/backend/internal/platform"
)

type Provider struct {
	mu       sync.Mutex
	apps     []platform.App
	settings map[string]any
	events   []platform.Event
	jobs     map[string]platform.Job
	catalog  catalog.Provider
	nextID   int
}

func New() *Provider {
	return &Provider{
		apps: []platform.App{
			{ID: "jellyfin", Kind: platform.AppKindUser, Name: "Jellyfin", Icon: "/app-icons/jellyfin.svg", Description: "Personal media server", Status: "running", Version: "10.10.7", URL: "http://labsos.local:8096", Installed: true},
			{ID: "immich", Kind: platform.AppKindUser, Name: "Immich", Icon: "/app-icons/immich.svg", Description: "Photo and video backup", Status: "stopped", Version: "1.0.0", Installed: true},
			{ID: "home-assistant", Kind: platform.AppKindUser, Name: "Home Assistant", Icon: "/app-icons/home-assistant.svg", Description: "Home automation", Status: "stopped"},
		},
		settings: map[string]any{"hostname": "labsos-dev", "timezone": "America/Sao_Paulo", "language": "pt-BR", "dhcp": true},
		events:   []platform.Event{{ID: "event-1", Type: "system", Message: "Labs API started"}},
		jobs:     make(map[string]platform.Job),
		catalog: catalog.NewBuiltInProvider([]catalog.RemoteApp{
			{ID: "jellyfin", Name: "Jellyfin", Description: "Personal media server", Category: "media", Version: "10.10.7", Icon: "/app-icons/jellyfin.svg"},
			{ID: "immich", Name: "Immich", Description: "Photo and video backup", Category: "storage", Version: "1.0.0", Icon: "/app-icons/immich.svg"},
			{ID: "home-assistant", Name: "Home Assistant", Description: "Home automation", Category: "utilities", Icon: "/app-icons/home-assistant.svg"},
		}),
	}
}

func (*Provider) SystemSummary(context.Context) (platform.SystemSummary, error) {
	return platform.SystemSummary{Hostname: "labsos-dev", Status: "healthy", UptimeSeconds: 483012, Version: "0.1.0", CPUUsage: 18.2, MemoryUsed: 3221225472, MemoryTotal: 8589934592, Temperature: 43, StorageUsed: 459561500672, StorageTotal: 999653638144}, nil
}

func (*Provider) SystemHealth(context.Context) (platform.SystemHealth, error) {
	return platform.SystemHealth{Status: "healthy", Components: map[string]string{"labs-api": "healthy", "provider": "healthy"}}, nil
}

func (p *Provider) Power(_ context.Context, action string) (platform.Job, error) {
	if action != "reboot" && action != "shutdown" {
		return platform.Job{}, platform.ErrUnsupported
	}
	return p.completeJob(action + " simulated"), nil
}

func (p *Provider) Apps(_ context.Context, catalog bool) ([]platform.App, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if catalog {
		items, err := p.catalog.ListApps()
		if err != nil {
			return nil, err
		}
		apps := make([]platform.App, 0, len(items))
		for _, item := range items {
			apps = append(apps, platform.App{ID: item.ID, Kind: platform.AppKindUser, Name: item.Name, Icon: item.Icon, Description: item.Description, Category: item.Category, Source: item.Source, Version: item.Version})
		}
		return apps, nil
	}
	apps := make([]platform.App, 0, len(p.apps))
	for _, app := range p.apps {
		if catalog || app.Installed {
			apps = append(apps, app)
		}
	}
	return apps, nil
}

func (p *Provider) AppAction(_ context.Context, id, action string) (platform.Job, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.apps {
		if p.apps[i].ID != id {
			continue
		}
		switch action {
		case "install", "start", "restart":
			p.apps[i].Installed = true
			p.apps[i].Status = "running"
		case "stop":
			p.apps[i].Status = "stopped"
		default:
			return platform.Job{}, platform.ErrUnsupported
		}
		return p.completeJobLocked(fmt.Sprintf("%s %s completed", id, action)), nil
	}
	return platform.Job{}, platform.ErrNotFound
}

func (p *Provider) RemoveApp(_ context.Context, id string) (platform.Job, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for i := range p.apps {
		if p.apps[i].ID == id && p.apps[i].Installed {
			p.apps[i].Installed = false
			p.apps[i].Status = "stopped"
			return p.completeJobLocked(id + " removed"), nil
		}
	}
	return platform.Job{}, platform.ErrNotFound
}

func (*Provider) AppLogs(context.Context, string, int) (string, error) { return "mock logs", nil }
func (*Provider) CatalogSources(context.Context) ([]platform.CatalogSource, error) {
	return []platform.CatalogSource{{ID: "builtin", Name: "Labs", Kind: "builtin"}}, nil
}
func (*Provider) UpdateStatus(context.Context) (platform.UpdateStatus, error) {
	return platform.UpdateStatus{CurrentVersion: "0.1.0", LatestVersion: "0.1.0"}, nil
}
func (*Provider) ApplyUpdate(context.Context) (platform.UpdateStatus, error) {
	return platform.UpdateStatus{CurrentVersion: "0.1.0", LatestVersion: "0.1.0"}, nil
}
func (*Provider) RollbackUpdate(context.Context) (platform.UpdateStatus, error) {
	return platform.UpdateStatus{CurrentVersion: "0.1.0", LatestVersion: "0.1.0"}, nil
}

func (p *Provider) Settings(context.Context) (map[string]any, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return clone(p.settings), nil
}

func (p *Provider) UpdateSettings(_ context.Context, section string, patch map[string]any) (map[string]any, error) {
	if section != "system" && section != "network" {
		return nil, platform.ErrUnsupported
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for key, value := range patch {
		p.settings[key] = value
	}
	return clone(p.settings), nil
}

func (p *Provider) Events(context.Context) ([]platform.Event, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]platform.Event(nil), p.events...), nil
}

func (p *Provider) Job(_ context.Context, id string) (platform.Job, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	job, ok := p.jobs[id]
	if !ok {
		return platform.Job{}, platform.ErrNotFound
	}
	return job, nil
}

func (p *Provider) completeJob(message string) platform.Job {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.completeJobLocked(message)
}

func (p *Provider) completeJobLocked(message string) platform.Job {
	p.nextID++
	id := fmt.Sprintf("job-%d", p.nextID)
	job := platform.Job{ID: id, Status: "success", Message: message}
	p.jobs[id] = job
	p.events = append(p.events, platform.Event{ID: fmt.Sprintf("event-%d", len(p.events)+1), Type: "job", Message: message})
	return job
}

func clone(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}
