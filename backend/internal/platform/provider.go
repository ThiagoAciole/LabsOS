package platform

import (
	"context"
	"errors"
)

var (
	ErrNotFound    = errors.New("not found")
	ErrUnsupported = errors.New("unsupported operation")
	ErrUnavailable = errors.New("provider unavailable")
)

type SystemSummary struct {
	Hostname      string  `json:"hostname"`
	Status        string  `json:"status"`
	UptimeSeconds int64   `json:"uptimeSeconds"`
	Version       string  `json:"version"`
	CPUUsage      float64 `json:"cpuUsage"`
	MemoryUsed    int64   `json:"memoryUsedBytes"`
	MemoryTotal   int64   `json:"memoryTotalBytes"`
	Temperature   float64 `json:"temperatureCelsius"`
	StorageUsed   int64   `json:"storageUsedBytes"`
	StorageTotal  int64   `json:"storageTotalBytes"`
	IPAddress     string  `json:"ipAddress,omitempty"`
	NetworkOnline bool    `json:"networkOnline"`
	NetworkRX     float64 `json:"networkDownloadBytesPerSecond"`
	NetworkTX     float64 `json:"networkUploadBytesPerSecond"`
}

type SystemHealth struct {
	Status     string            `json:"status"`
	Mode       string            `json:"mode"`
	Components map[string]string `json:"components"`
}

type App struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Icon            string `json:"icon"`
	Description     string `json:"description"`
	Status          string `json:"status"`
	Version         string `json:"version,omitempty"`
	URL             string `json:"url,omitempty"`
	UpdateAvailable bool   `json:"updateAvailable"`
	Installed       bool   `json:"installed"`
}

type Event struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Job struct {
	ID      string `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Provider interface {
	Mode() string
	SystemSummary(context.Context) (SystemSummary, error)
	SystemHealth(context.Context) (SystemHealth, error)
	Power(context.Context, string) (Job, error)
	Apps(context.Context, bool) ([]App, error)
	AppAction(context.Context, string, string) (Job, error)
	RemoveApp(context.Context, string) (Job, error)
	Settings(context.Context) (map[string]any, error)
	UpdateSettings(context.Context, string, map[string]any) (map[string]any, error)
	Events(context.Context) ([]Event, error)
	Job(context.Context, string) (Job, error)
}
