package systemapps

import "context"

type App struct {
	ID  string
	URL string
}

type Lifecycle interface {
	EnsureInstalled(context.Context) error
	Start(context.Context) error
	Stop(context.Context) error
	Restart(context.Context) error
	Healthy(context.Context) (bool, error)
	URL(context.Context) (string, error)
}
