package catalog

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
)

type Fetch func(context.Context) ([]App, error)

type FileCache struct{ Path string }

func (c FileCache) Save(apps []App) error {
	data, err := json.Marshal(apps)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(c.Path), 0o750); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(c.Path), ".catalog-*")
	if err != nil {
		return err
	}
	name := tmp.Name()
	defer os.Remove(name)
	if err := tmp.Chmod(0o600); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(name, c.Path)
}

func (c FileCache) Load() ([]App, error) {
	data, err := os.ReadFile(c.Path)
	if err != nil {
		return nil, err
	}
	var apps []App
	if err := json.Unmarshal(data, &apps); err != nil {
		return nil, err
	}
	return apps, nil
}

func LoadWithCache(ctx context.Context, fetch Fetch, cache FileCache) ([]App, error) {
	apps, err := fetch(ctx)
	if err == nil && len(apps) > 0 {
		if saveErr := cache.Save(apps); saveErr == nil {
			return apps, nil
		}
	}
	return cache.Load()
}
