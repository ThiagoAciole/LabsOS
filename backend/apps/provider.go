package apps

import "context"

type Status string

const (
	AppStatusInstalling Status = "installing"
	AppStatusRunning    Status = "running"
	AppStatusStopped    Status = "stopped"
	AppStatusError      Status = "error"
)

type App struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Status          Status `json:"status"`
	URL             string `json:"url,omitempty"`
	Version         string `json:"version,omitempty"`
	UpdateAvailable bool   `json:"updateAvailable"`
}

type Provider interface {
	List(context.Context) ([]App, error)
	Install(context.Context, string) (string, error)
	Start(context.Context, string) error
	Stop(context.Context, string) error
	Restart(context.Context, string) error
	Remove(context.Context, string) error
}
