package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type storageUsage struct {
	Path           string `json:"path"`
	UsedBytes      uint64 `json:"usedBytes"`
	AvailableBytes uint64 `json:"availableBytes"`
	TotalBytes     uint64 `json:"totalBytes"`
	UsePercent     int    `json:"usePercent"`
}

type storageDevice struct {
	Name        string          `json:"name"`
	Path        string          `json:"path"`
	Size        uint64          `json:"size"`
	FSType      string          `json:"fstype,omitempty"`
	UUID        string          `json:"uuid,omitempty"`
	Mountpoints []string        `json:"mountpoints,omitempty"`
	Type        string          `json:"type"`
	ReadOnly    bool            `json:"readOnly"`
	Transport   string          `json:"transport,omitempty"`
	Model       string          `json:"model,omitempty"`
	Health      string          `json:"health,omitempty"`
	Children    []storageDevice `json:"children,omitempty"`
}

func (s *server) storage(w http.ResponseWriter, r *http.Request) {
	devices := []storageDevice{}
	if output, err := exec.CommandContext(r.Context(), "lsblk", "-J", "-b", "-o", "NAME,PATH,SIZE,FSTYPE,UUID,MOUNTPOINTS,TYPE,RO,TRAN,MODEL").Output(); err == nil {
		var payload struct {
			Blockdevices []storageDevice `json:"blockdevices"`
		}
		if json.Unmarshal(output, &payload) == nil {
			devices = payload.Blockdevices
			for index := range devices {
				populateStorageHealth(r.Context(), &devices[index])
			}
		}
	}
	usage := []storageUsage{}
	for _, path := range []string{"/DATA", "/"} {
		if output, err := exec.CommandContext(r.Context(), "df", "-B1", "--output=target,used,avail,size,pcent", path).Output(); err == nil {
			if item, ok := parseDF(string(output)); ok {
				item.Path = path
				usage = append(usage, item)
			}
		}
	}
	dataMounted := false
	for _, item := range usage {
		if item.Path == "/DATA" {
			dataMounted = true
		}
	}
	respond(w, map[string]any{"devices": devices, "usage": usage, "dataMounted": dataMounted, "readOnly": true, "message": "storage discovery is read-only"}, nil, http.StatusOK)
}

func populateStorageHealth(ctx context.Context, device *storageDevice) {
	if device.Type == "disk" && device.Path != "" {
		device.Health = smartHealth(ctx, device.Path)
	}
	for index := range device.Children {
		populateStorageHealth(ctx, &device.Children[index])
	}
}

func smartHealth(ctx context.Context, path string) string {
	output, err := exec.CommandContext(ctx, "smartctl", "-H", path).CombinedOutput()
	text := strings.ToLower(string(output))
	if err != nil && text == "" {
		return "unavailable"
	}
	if strings.Contains(text, "passed") {
		return "healthy"
	}
	if strings.Contains(text, "failed") {
		return "failed"
	}
	return "unknown"
}

func parseDF(output string) (storageUsage, bool) {
	lines := strings.Fields(output)
	if len(lines) < 10 {
		return storageUsage{}, false
	}
	// Fields are: target used avail size pcent; the header contributes five fields.
	fields := lines[len(lines)-5:]
	used, e1 := strconv.ParseUint(fields[1], 10, 64)
	avail, e2 := strconv.ParseUint(fields[2], 10, 64)
	total, e3 := strconv.ParseUint(fields[3], 10, 64)
	percent, e4 := strconv.Atoi(strings.TrimSuffix(fields[4], "%"))
	if e1 != nil || e2 != nil || e3 != nil || e4 != nil {
		return storageUsage{}, false
	}
	return storageUsage{UsedBytes: used, AvailableBytes: avail, TotalBytes: total, UsePercent: percent}, true
}
