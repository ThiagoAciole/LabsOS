package linux

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"labsos/backend/internal/platform"
	"labsos/backend/providers/mock"
)

type Provider struct{ fallback *mock.Provider }

func New() *Provider           { return &Provider{fallback: mock.New()} }
func (*Provider) Mode() string { return "linux" }

func (*Provider) SystemSummary(ctx context.Context) (platform.SystemSummary, error) {
	statBefore, err := os.ReadFile("/proc/stat")
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	netBefore, _ := os.ReadFile("/proc/net/dev")
	select {
	case <-ctx.Done():
		return platform.SystemSummary{}, ctx.Err()
	case <-time.After(100 * time.Millisecond):
	}
	statAfter, err := os.ReadFile("/proc/stat")
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	mem, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	uptimeRaw, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}

	before, err := parseCPUStat(string(statBefore))
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	after, err := parseCPUStat(string(statAfter))
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	memoryUsed, memoryTotal, err := parseMemInfo(string(mem))
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	uptime, err := parseUptime(string(uptimeRaw))
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}
	hostname, err := os.Hostname()
	if err != nil {
		return platform.SystemSummary{}, platform.ErrUnavailable
	}

	summary := platform.SystemSummary{Hostname: hostname, Status: "healthy", UptimeSeconds: uptime, Version: "0.1.0", CPUUsage: cpuUsage(before, after), MemoryUsed: memoryUsed, MemoryTotal: memoryTotal, Temperature: temperature(), NetworkOnline: false}
	if address := primaryAddress(); address != "" {
		summary.IPAddress, summary.NetworkOnline = address, true
	}
	if netAfter, readErr := os.ReadFile("/proc/net/dev"); readErr == nil {
		rxBefore, txBefore, beforeErr := parseNetworkCounters(string(netBefore))
		rxAfter, txAfter, afterErr := parseNetworkCounters(string(netAfter))
		if beforeErr == nil && afterErr == nil {
			summary.NetworkRX, summary.NetworkTX = float64(rxAfter-rxBefore)*10, float64(txAfter-txBefore)*10
		}
	}
	return summary, nil
}

func (*Provider) SystemHealth(context.Context) (platform.SystemHealth, error) {
	if _, err := os.Stat("/proc/stat"); err != nil {
		return platform.SystemHealth{Status: "unavailable", Mode: "linux"}, platform.ErrUnavailable
	}
	return platform.SystemHealth{Status: "healthy", Mode: "linux", Components: map[string]string{"labs-api": "healthy", "system-provider": "healthy"}}, nil
}

func (*Provider) Power(context.Context, string) (platform.Job, error) {
	return platform.Job{}, platform.ErrUnavailable
}
func (p *Provider) Apps(ctx context.Context, catalog bool) ([]platform.App, error) {
	return p.fallback.Apps(ctx, catalog)
}
func (p *Provider) AppAction(ctx context.Context, id, action string) (platform.Job, error) {
	return p.fallback.AppAction(ctx, id, action)
}
func (p *Provider) RemoveApp(ctx context.Context, id string) (platform.Job, error) {
	return p.fallback.RemoveApp(ctx, id)
}
func (p *Provider) Settings(ctx context.Context) (map[string]any, error) {
	return p.fallback.Settings(ctx)
}
func (p *Provider) UpdateSettings(ctx context.Context, section string, patch map[string]any) (map[string]any, error) {
	return p.fallback.UpdateSettings(ctx, section, patch)
}
func (p *Provider) Events(ctx context.Context) ([]platform.Event, error) {
	return p.fallback.Events(ctx)
}
func (p *Provider) Job(ctx context.Context, id string) (platform.Job, error) {
	return p.fallback.Job(ctx, id)
}

func temperature() float64 {
	paths, _ := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(string(raw)), 64)
		if err == nil {
			if value > 1000 {
				value /= 1000
			}
			return value
		}
	}
	return 0
}

func primaryAddress() string {
	interfaces, _ := net.Interfaces()
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, _ := iface.Addrs()
		for _, address := range addresses {
			ip, _, err := net.ParseCIDR(address.String())
			if err == nil && ip.To4() != nil {
				return ip.String()
			}
		}
	}
	return ""
}
