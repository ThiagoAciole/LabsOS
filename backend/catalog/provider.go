package catalog

import "fmt"

type Provider interface {
	ListApps() ([]App, error)
	GetApp(string) (App, error)
}

type builtInProvider struct{ apps []App }

func NewBuiltInProvider(remote []RemoteApp) Provider {
	apps := make([]App, 0, len(remote))
	for _, item := range remote {
		if app, err := Normalize(item, "builtin"); err == nil {
			apps = append(apps, app)
		}
	}
	return &builtInProvider{apps: apps}
}

func (p *builtInProvider) ListApps() ([]App, error) {
	return append([]App(nil), p.apps...), nil
}

func (p *builtInProvider) GetApp(id string) (App, error) {
	for _, app := range p.apps {
		if app.ID == id {
			return app, nil
		}
	}
	return App{}, fmt.Errorf("catalog app %q not found", id)
}
