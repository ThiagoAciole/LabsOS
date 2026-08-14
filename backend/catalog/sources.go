package catalog

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

type Source struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URL  string `json:"url"`
	Kind string `json:"kind"`
}

func DefaultSources() []Source {
	return []Source{
		{ID: "bigbear", Name: "BigBearCasaOS", URL: "https://github.com/bigbeartechworld/big-bear-casaos", Kind: "git-compose"},
		{ID: "casaos", Name: "CasaOS/ZimaOS AppStore", URL: "https://github.com/IceWhaleTech/CasaOS-AppStore", Kind: "casaos"},
		{ID: "lissy93", Name: "Lissy93 Portainer Templates", URL: "https://github.com/Lissy93/portainer-templates", Kind: "portainer"},
		{ID: "selfhostedpro", Name: "SelfhostedPro Portainer Templates", URL: "https://github.com/SelfhostedPro/selfhosted_templates", Kind: "portainer"},
	}
}

type CasaOSCatalogAdapter struct {
	Remote     RemoteProvider
	SourceName string
}

func (a CasaOSCatalogAdapter) ListApps() ([]App, error) {
	apps, err := a.Remote.ListApps(context.Background())
	if err != nil {
		return nil, err
	}
	for i := range apps {
		apps[i].Source = a.SourceName
	}
	return apps, nil
}
func (a CasaOSCatalogAdapter) GetApp(id string) (App, error) {
	apps, err := a.ListApps()
	if err != nil {
		return App{}, err
	}
	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return App{}, ErrNotFound
}

type GitComposeProvider struct {
	URL, SourceName string
	Client          *http.Client
}

func (p GitComposeProvider) ListApps() ([]App, error) {
	dir, err := os.MkdirTemp("", "labsos-catalog-")
	if err != nil {
		return nil, err
	}
	if err := exec.Command("git", "clone", "--depth", "1", p.URL, dir).Run(); err != nil {
		return nil, err
	}
	apps := []App{}
	err = filepath.Walk(dir, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil || info.IsDir() || (info.Name() != "compose.yml" && info.Name() != "docker-compose.yml") {
			return walkErr
		}
		app, importErr := ImportCompose(context.Background(), path)
		if importErr != nil {
			return nil
		}
		app.Source = p.SourceName
		apps = append(apps, app)
		return nil
	})
	return apps, err
}
func (p GitComposeProvider) GetApp(id string) (App, error) {
	apps, err := p.ListApps()
	if err != nil {
		return App{}, err
	}
	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return App{}, fmt.Errorf("catalog app %q not found", id)
}

type MergedProvider struct{ Providers []Provider }

func (p MergedProvider) ListApps() ([]App, error) {
	result := []App{}
	for _, provider := range p.Providers {
		apps, err := provider.ListApps()
		if err != nil {
			continue
		}
		result = append(result, apps...)
	}
	return result, nil
}
func (p MergedProvider) GetApp(id string) (App, error) {
	apps, err := p.ListApps()
	if err != nil {
		return App{}, err
	}
	for _, app := range apps {
		if app.ID == id {
			return app, nil
		}
	}
	return App{}, ErrNotFound
}
