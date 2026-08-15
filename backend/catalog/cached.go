package catalog

import "context"

type CachedFetchProvider struct {
	Fetch   Fetch
	Resolve func(string) (App, error)
	Cache   FileCache
}

func (p CachedFetchProvider) ListApps() ([]App, error) {
	if apps, err := p.Cache.Load(); err == nil && len(apps) > 0 { return apps, nil }
	return LoadWithCache(context.Background(), p.Fetch, p.Cache)
}

func (p CachedFetchProvider) GetApp(id string) (App, error) {
	apps, err := p.ListApps(); if err != nil { return App{}, err }
	for _, app := range apps {
		if app.ID == id {
			if app.ComposeURL == "" && p.Resolve != nil {
				return p.Resolve(id)
			}
			return app, nil
		}
	}
	return App{}, ErrNotFound
}

type CachedProvider struct {
	Remote RemoteProvider
	Cache  FileCache
}

func (p CachedProvider) ListApps() ([]App, error) {
	apps, err := LoadWithCache(context.Background(), p.Remote.ListApps, p.Cache)
	for i := range apps {
		apps[i].Installable = apps[i].ID == "jellyfin" || apps[i].ID == "syncthing"
	}
	return apps, err
}

func (p CachedProvider) GetApp(id string) (App, error) {
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
