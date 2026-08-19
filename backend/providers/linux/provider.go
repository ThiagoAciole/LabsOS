package linux

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"labsos/backend/catalog"
	"labsos/backend/internal/platform"
	"labsos/backend/internal/safety"
	"labsos/backend/internal/version"
	"labsos/backend/labsd"
	"labsos/backend/update"
)

type Provider struct {
	apps    labsd.Client
	catalog catalog.Provider
	power   func(context.Context, string) error
}

type persistedSettings struct {
	Language     string `json:"language,omitempty"`
	RemoteAccess bool   `json:"remoteAccess"`
}

func New() *Provider {
	url := os.Getenv("LABSOS_CATALOG_URL")
	if url == "" {
		url = "https://raw.githubusercontent.com/IceWhaleTech/CasaOS-AppStore/gh-pages/index.json"
	}
	remote := catalog.CasaOSCatalogAdapter{Remote: catalog.RemoteProvider{URL: url}, SourceName: "CasaOS/ZimaOS AppStore"}
	bigbear := catalog.GitComposeProvider{URL: "https://github.com/bigbeartechworld/big-bear-casaos", SourceName: "BigBearCasaOS"}
	provider := catalog.CachedFetchProvider{Resolve: func(id string) (catalog.App, error) {
		if app, err := remote.Remote.GetApp(id); err == nil {
			return app, nil
		}
		return bigbear.GetApp(id)
	}, Fetch: func(ctx context.Context) ([]catalog.App, error) {
		casa, _ := remote.Remote.ListApps(ctx)
		for i := range casa {
			casa[i].Source = remote.SourceName
		}
		bear, _ := bigbear.ListApps()
		return append(casa, bear...), nil
	}, Cache: catalog.FileCache{Path: "/var/lib/labsos/catalog/apps.json"}}
	return &Provider{apps: labsd.Client{Socket: "/run/labsos/labsd.sock"}, catalog: provider, power: runPower}
}
func NewWithApps(apps labsd.Client) *Provider {
	return &Provider{apps: apps, power: runPower}
}

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
	storageUsed, storageTotal := filesystemUsage("/DATA")
	if storageTotal == 0 {
		storageUsed, storageTotal = filesystemUsage("/")
	}

	summary := platform.SystemSummary{Hostname: hostname, Status: "healthy", UptimeSeconds: uptime, Version: version.Value, CPUUsage: cpuUsage(before, after), MemoryUsed: memoryUsed, MemoryTotal: memoryTotal, StorageUsed: storageUsed, StorageTotal: storageTotal, Temperature: temperature(), NetworkOnline: false}
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
		return platform.SystemHealth{Status: "unavailable"}, platform.ErrUnavailable
	}
	return platform.SystemHealth{Status: "healthy", Components: map[string]string{"labs-api": "healthy", "system-provider": "healthy"}}, nil
}

func (p *Provider) Power(ctx context.Context, action string) (platform.Job, error) {
	if action != "reboot" && action != "shutdown" {
		return platform.Job{}, platform.ErrUnsupported
	}
	runner := p.power
	if runner == nil {
		runner = runPower
	}
	if err := runner(ctx, action); err != nil {
		return platform.Job{}, platform.ErrUnavailable
	}
	return platform.Job{ID: "power-" + action, Status: "accepted", Message: action + " accepted by systemd"}, nil
}

func runPower(ctx context.Context, action string) error {
	return exec.CommandContext(ctx, "systemctl", action).Run()
}
func (p *Provider) Apps(ctx context.Context, catalogApps bool) ([]platform.App, error) {
	if catalogApps {
		if p.catalog == nil {
			return nil, platform.ErrUnavailable
		}
		items, err := p.catalog.ListApps()
		if err != nil {
			return nil, platform.ErrUnavailable
		}
		result := make([]platform.App, 0, len(items))
		for _, item := range items {
			result = append(result, platform.App{ID: item.ID, Kind: platform.AppKindUser, Name: item.Name, Icon: item.Icon, Description: item.Description, Category: item.Category, Source: item.Source, Version: item.Version, Installable: item.Installable, Actions: []string{"install"}, Architecture: item.Architecture, Requirements: item.Requirements})
		}
		if !catalogHas(result, "jellyfin") {
			result = append(result, platform.App{ID: "jellyfin", Kind: platform.AppKindUser, Name: "Jellyfin", Icon: "/app-icons/jellyfin.svg", Description: "Personal media server", Category: "media", Source: "labs", Version: "10.10.7", Installable: true, Actions: []string{"install"}, Architecture: []string{"amd64", "arm64"}})
		}
		return result, nil
	}
	names, err := p.apps.List(ctx)
	if err != nil {
		status, statusErr := p.apps.Status(ctx, "jellyfin")
		if statusErr != nil {
			return nil, platform.ErrUnavailable
		}
		if status != "running" {
			status = "stopped"
		}
		return []platform.App{{ID: "jellyfin", Kind: platform.AppKindUser, Name: "Jellyfin", Icon: "/app-icons/jellyfin.svg", Description: "Personal media server", Status: status, Installed: true}}, nil
	}
	result := make([]platform.App, 0, len(names))
	for _, id := range names {
		if strings.Contains(id, "-") && !labsdHasApp([]string{"jellyfin", "syncthing"}, id) {
			continue
		}
		name, icon, description := id, "", "Installed Docker Compose app"
		switch id {
		case "jellyfin":
			name, icon, description = "Jellyfin", "/app-icons/jellyfin.svg", "Personal media server"
		case "syncthing":
			name, description = "Syncthing", "Continuous file synchronization"
		}
		composePath := appComposePath(id)
		ports := catalog.PublishedPorts(composePath)
		url := appURL(id)
		status, statusErr := p.apps.Status(ctx, id)
		if statusErr != nil || status != "running" {
			status = "stopped"
		}
		result = append(result, platform.App{ID: id, Kind: platform.AppKindUser, Name: name, Icon: icon, Description: description, Status: status, Health: status, URL: url, Ports: ports, Installed: true, Actions: []string{"start", "stop", "restart", "remove"}})
	}
	return result, nil
}

func appURL(id string) string {
	composePath := appComposePath(id)
	if port := catalog.PublishedPort(composePath); port > 0 {
		host := envOr("LABSOS_PUBLIC_HOST", "192.168.0.2")
		return fmt.Sprintf("http://%s:%d", host, port)
	}
	switch id {
	case "jellyfin":
		return "http://labsos.local:8096"
	case "syncthing":
		return "http://labsos.local:8384"
	case "2fauth":
		return "http://labsos.local:8000"
	default:
		return ""
	}
}

func appComposePath(id string) string {
	composePath := filepath.Join("/opt/labsos/apps", id, "compose.yaml")
	if _, err := os.Stat(composePath); err != nil {
		matches, _ := filepath.Glob(filepath.Join("/opt/labsos/apps", "*"+id+"*", "compose.yaml"))
		if len(matches) == 1 {
			composePath = matches[0]
		}
	}
	return composePath
}

func labsdHasApp(names []string, id string) bool {
	for _, name := range names {
		if name == id {
			return true
		}
	}
	return false
}
func catalogHas(apps []platform.App, id string) bool {
	for _, app := range apps {
		if app.ID == id {
			return true
		}
	}
	return false
}
func (p *Provider) AppAction(ctx context.Context, id, action string) (platform.Job, error) {
	if id != "" {
		operation := map[string]string{"install": "InstallApp", "start": "StartApp", "stop": "StopApp", "restart": "RestartApp", "update": "UpdateApp"}[action]
		if operation == "" {
			return platform.Job{}, platform.ErrUnsupported
		}
		var err error
		if action == "install" && p.catalog != nil {
			app, lookupErr := p.catalog.GetApp(id)
			if lookupErr == nil && app.Compose != "" {
				err = p.apps.Install(ctx, id, app.Compose)
			} else {
				err = p.apps.Call(ctx, operation, id)
			}
		} else {
			err = p.apps.Call(ctx, operation, id)
		}
		if err != nil {
			return platform.Job{}, platform.ErrUnavailable
		}
		return platform.Job{ID: id + "-" + action, Status: "success", Message: id + " " + action + " completed"}, nil
	}
	return platform.Job{}, platform.ErrNotFound
}

// AppActionWithCompose installs a locally registered declarative manifest
// without resolving it through a remote catalog.
func (p *Provider) AppActionWithCompose(ctx context.Context, id, action, compose string) (platform.Job, error) {
	if action != "install" || strings.TrimSpace(compose) == "" || id == "" {
		return platform.Job{}, platform.ErrUnsupported
	}
	if err := p.apps.Install(ctx, id, compose); err != nil {
		return platform.Job{}, platform.ErrUnavailable
	}
	return platform.Job{ID: id + "-install", Status: "success", Message: id + " installed from declarative manifest"}, nil
}
func (p *Provider) RemoveApp(ctx context.Context, id string) (platform.Job, error) {
	if id != "" {
		if err := p.apps.Call(ctx, "RemoveApp", id); err != nil {
			return platform.Job{}, platform.ErrUnavailable
		}
		return platform.Job{ID: id + "-remove", Status: "success", Message: id + " removed"}, nil
	}
	return platform.Job{}, platform.ErrNotFound
}

func (p *Provider) AppLogs(ctx context.Context, id string, lines int) (string, error) {
	logs, err := p.apps.Logs(ctx, id, lines)
	if err != nil {
		return "", platform.ErrUnavailable
	}
	return logs, nil
}
func (*Provider) CatalogSources(context.Context) ([]platform.CatalogSource, error) {
	items := catalog.DefaultSources()
	result := make([]platform.CatalogSource, len(items))
	for i, item := range items {
		result[i] = platform.CatalogSource{ID: item.ID, Name: item.Name, URL: item.URL, Kind: item.Kind}
	}
	return result, nil
}
func (p *Provider) UpdateStatus(ctx context.Context) (platform.UpdateStatus, error) {
	manager := update.Manager{Root: "/opt/labsos", ManifestURL: os.Getenv("LABSOS_UPDATE_MANIFEST_URL")}
	status, err := manager.Check(ctx)
	return platform.UpdateStatus{CurrentVersion: status.CurrentVersion, LatestVersion: status.LatestVersion, UpdateAvailable: status.UpdateAvailable, Channel: status.Channel, Changelog: status.Changelog, SHA256: status.SHA256}, err
}
func (p *Provider) ApplyUpdate(ctx context.Context) (platform.UpdateStatus, error) {
	manager := update.Manager{Root: envOr("LABSOS_UPDATE_STAGING", "/var/lib/labsos/update"), ManifestURL: os.Getenv("LABSOS_UPDATE_MANIFEST_URL")}
	status, err := manager.Update(ctx)
	if err != nil {
		return platform.UpdateStatus{}, err
	}
	stagedRelease := filepath.Join(manager.Root, "releases", status.CurrentVersion)
	release := filepath.Join("/opt/labsos/releases", status.CurrentVersion)
	if err := exec.CommandContext(ctx, "sudo", "-n", "mkdir", "-p", release).Run(); err != nil {
		return platform.UpdateStatus{}, err
	}
	if err := exec.CommandContext(ctx, "sudo", "-n", "cp", "-a", stagedRelease+"/.", release+"/").Run(); err != nil {
		return platform.UpdateStatus{}, err
	}
	if err := exec.CommandContext(ctx, "sudo", "-n", "ln", "-sfn", release, "/opt/labsos/current").Run(); err != nil {
		return platform.UpdateStatus{}, err
	}
	for _, binary := range []string{"labs-api", "labsd", "labs-dashboard"} {
		if err := exec.CommandContext(ctx, "sudo", "-n", "install", "-m", "0755", filepath.Join(release, "bin", binary), "/usr/local/bin/"+binary).Run(); err != nil {
			return platform.UpdateStatus{}, err
		}
	}
	if err := exec.CommandContext(ctx, "sudo", "-n", "mkdir", "-p", "/opt/labsos/dashboard.new").Run(); err != nil {
		return platform.UpdateStatus{}, err
	}
	if err := exec.CommandContext(ctx, "sudo", "-n", "cp", "-a", filepath.Join(release, "dashboard", "."), "/opt/labsos/dashboard.new/").Run(); err != nil {
		return platform.UpdateStatus{}, err
	}
	if err := exec.CommandContext(ctx, "sudo", "-n", "mv", "-T", "/opt/labsos/dashboard.new", "/opt/labsos/dashboard").Run(); err != nil {
		return platform.UpdateStatus{}, err
	}
	go exec.Command("sudo", "-n", "systemctl", "restart", "labs-api", "labsd", "labs-dashboard").Run()
	return platform.UpdateStatus{CurrentVersion: status.CurrentVersion, LatestVersion: status.LatestVersion}, nil
}
func (p *Provider) RollbackUpdate(context.Context) (platform.UpdateStatus, error) {
	manager := update.Manager{Root: "/opt/labsos"}
	status, err := manager.Rollback()
	if err != nil {
		return platform.UpdateStatus{}, err
	}
	go exec.Command("sudo", "-n", "systemctl", "restart", "labs-api", "labsd", "labs-dashboard").Run()
	return platform.UpdateStatus{CurrentVersion: status.CurrentVersion, LatestVersion: status.LatestVersion}, nil
}
func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
func (*Provider) Settings(context.Context) (map[string]any, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, platform.ErrUnavailable
	}
	timezone := "UTC"
	if data, readErr := os.ReadFile("/etc/timezone"); readErr == nil && strings.TrimSpace(string(data)) != "" {
		timezone = strings.TrimSpace(string(data))
	} else if link, linkErr := os.Readlink("/etc/localtime"); linkErr == nil {
		const marker = "/zoneinfo/"
		if index := strings.Index(link, marker); index >= 0 {
			timezone = link[index+len(marker):]
		}
	}
	persisted := loadPersistedSettings()
	language := persisted.Language
	if language == "" {
		language = envOr("LABSOS_LANGUAGE", "pt-BR")
	}
	return map[string]any{"hostname": hostname, "timezone": timezone, "language": language, "remoteAccess": persisted.RemoteAccess}, nil
}
func (p *Provider) UpdateSettings(ctx context.Context, section string, patch map[string]any) (map[string]any, error) {
	if section != "system" && section != "network" {
		return nil, platform.ErrUnsupported
	}
	if !safety.RealOperationsEnabled() || section == "network" && !safety.NetworkChangesEnabled() {
		return nil, platform.ErrUnavailable
	}
	persisted := loadPersistedSettings()
	if section == "system" {
		hostname, ok := patch["hostname"].(string)
		if !ok || strings.TrimSpace(hostname) == "" || !validHostname(hostname) {
			return nil, platform.ErrUnsupported
		}
		if err := os.WriteFile("/etc/hostname", []byte(strings.TrimSpace(hostname)+"\n"), 0644); err != nil {
			return nil, platform.ErrUnavailable
		}
		if language, ok := patch["language"].(string); ok {
			language = strings.TrimSpace(language)
			if language != "pt-BR" && language != "en-US" {
				return nil, platform.ErrUnsupported
			}
			persisted.Language = language
		}
	} else {
		remoteAccess, ok := patch["remoteAccess"].(bool)
		if !ok {
			return nil, platform.ErrUnsupported
		}
		persisted.RemoteAccess = remoteAccess
	}
	if err := savePersistedSettings(persisted); err != nil {
		return nil, platform.ErrUnavailable
	}
	return p.Settings(ctx)
}

func settingsStatePath() string {
	if value := os.Getenv("LABSOS_SETTINGS_FILE"); value != "" {
		return value
	}
	return "/var/lib/labsos/state/settings.json"
}

func loadPersistedSettings() persistedSettings {
	data, err := os.ReadFile(settingsStatePath())
	if err != nil {
		return persistedSettings{}
	}
	var value persistedSettings
	if json.Unmarshal(data, &value) != nil {
		return persistedSettings{}
	}
	return value
}

func savePersistedSettings(value persistedSettings) error {
	path := settingsStatePath()
	if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".settings-*.tmp")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = tmp.Write(data); err == nil {
		err = tmp.Chmod(0600)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}
func (*Provider) Events(context.Context) ([]platform.Event, error) {
	return []platform.Event{}, nil
}
func (*Provider) Job(context.Context, string) (platform.Job, error) {
	return platform.Job{}, platform.ErrNotFound
}

func validHostname(value string) bool {
	value = strings.TrimSpace(value)
	if len(value) == 0 || len(value) > 253 || strings.HasPrefix(value, ".") || strings.HasSuffix(value, ".") {
		return false
	}
	for _, part := range strings.Split(value, ".") {
		if len(part) == 0 || len(part) > 63 || part[0] == '-' || part[len(part)-1] == '-' {
			return false
		}
		for _, char := range part {
			if !(char == '-' || char >= '0' && char <= '9' || char >= 'A' && char <= 'Z' || char >= 'a' && char <= 'z') {
				return false
			}
		}
	}
	return true
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
