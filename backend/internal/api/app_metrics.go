package api

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
)

type appMetrics struct {
	ID               string  `json:"id"`
	CPUPercent       float64 `json:"cpuPercent"`
	MemoryUsedBytes  uint64  `json:"memoryUsedBytes"`
	MemoryLimitBytes uint64  `json:"memoryLimitBytes"`
	MemoryPercent    float64 `json:"memoryPercent"`
	NetworkRXBytes   uint64  `json:"networkRXBytes"`
	NetworkTXBytes   uint64  `json:"networkTXBytes"`
	BlockReadBytes   uint64  `json:"blockReadBytes"`
	BlockWriteBytes  uint64  `json:"blockWriteBytes"`
	Available        bool    `json:"available"`
	Message          string  `json:"message"`
}

func (s *server) appMetrics(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if !safeAppMetricID(id) {
		writeError(w, http.StatusBadRequest, "INVALID_APP_ID", "invalid app id")
		return
	}
	result := appMetrics{ID: id, Message: "container metrics unavailable"}
	output, err := exec.CommandContext(r.Context(), "docker", "stats", "--no-stream", "--format", "{{json .}}", id).Output()
	if err != nil {
		respond(w, result, nil, http.StatusOK)
		return
	}
	var raw struct {
		CPUPerc  string `json:"CPUPerc"`
		MemUsage string `json:"MemUsage"`
		MemPerc  string `json:"MemPerc"`
		NetIO    string `json:"NetIO"`
		BlockIO  string `json:"BlockIO"`
	}
	if err := json.Unmarshal([]byte(strings.TrimSpace(string(output))), &raw); err != nil {
		respond(w, result, nil, http.StatusOK)
		return
	}
	result.CPUPercent = parsePercent(raw.CPUPerc)
	result.MemoryUsedBytes, result.MemoryLimitBytes = parsePairBytes(raw.MemUsage)
	result.MemoryPercent = parsePercent(raw.MemPerc)
	result.NetworkRXBytes, result.NetworkTXBytes = parsePairBytes(raw.NetIO)
	result.BlockReadBytes, result.BlockWriteBytes = parsePairBytes(raw.BlockIO)
	result.Available, result.Message = true, "container metrics collected read-only"
	respond(w, result, nil, http.StatusOK)
}

func safeAppMetricID(value string) bool {
	if value == "" || len(value) > 63 {
		return false
	}
	for _, r := range value {
		if !(r == '-' || r == '_' || r == '.' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9') {
			return false
		}
	}
	return true
}
func parsePercent(value string) float64 {
	value = strings.TrimSpace(strings.TrimSuffix(value, "%"))
	result, _ := strconv.ParseFloat(value, 64)
	return result
}
func parsePairBytes(value string) (uint64, uint64) {
	parts := strings.Split(value, "/")
	if len(parts) != 2 {
		return 0, 0
	}
	return parseHumanBytes(parts[0]), parseHumanBytes(parts[1])
}
func parseHumanBytes(value string) uint64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0
	}
	split := 0
	for split < len(value) && (value[split] >= '0' && value[split] <= '9' || value[split] == '.') {
		split++
	}
	number, _ := strconv.ParseFloat(value[:split], 64)
	unit := strings.ToUpper(strings.TrimSpace(value[split:]))
	multiplier := float64(1)
	switch unit {
	case "K", "KB", "KI", "KIB":
		multiplier = 1024
	case "M", "MB", "MI", "MIB":
		multiplier = 1024 * 1024
	case "G", "GB", "GI", "GIB":
		multiplier = 1024 * 1024 * 1024
	case "T", "TB", "TI", "TIB":
		multiplier = 1024 * 1024 * 1024 * 1024
	}
	return uint64(number * multiplier)
}
