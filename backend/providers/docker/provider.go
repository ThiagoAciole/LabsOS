package docker

import (
	"context"
	"errors"

	"labsos/backend/apps"
	"labsos/backend/runtime"
)

var ErrUnavailable = errors.New("docker runtime unavailable")

type Provider struct{ runtime runtime.ContainerRuntime }

func New(containerRuntime runtime.ContainerRuntime) *Provider {
	return &Provider{runtime: containerRuntime}
}

func (p *Provider) ready(ctx context.Context) error {
	available, err := p.runtime.Available(ctx)
	if err != nil || !available {
		return ErrUnavailable
	}
	return nil
}

func (p *Provider) List(ctx context.Context) ([]apps.App, error) {
	if err := p.ready(ctx); err != nil {
		return nil, err
	}
	return []apps.App{}, nil
}

func (p *Provider) Install(ctx context.Context, _ string) (string, error) {
	if err := p.ready(ctx); err != nil {
		return "", err
	}
	return "", errors.New("docker app installation is not enabled")
}

func (p *Provider) Start(ctx context.Context, _ string) error   { return p.ready(ctx) }
func (p *Provider) Stop(ctx context.Context, _ string) error    { return p.ready(ctx) }
func (p *Provider) Restart(ctx context.Context, _ string) error { return p.ready(ctx) }
func (p *Provider) Remove(ctx context.Context, _ string) error  { return p.ready(ctx) }
