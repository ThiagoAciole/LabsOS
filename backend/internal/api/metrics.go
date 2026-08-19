package api

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type metricsSnapshot struct {
	CollectedAt time.Time         `json:"collectedAt"`
	CPU         float64           `json:"cpuPercent"`
	Memory      metricValue       `json:"memory"`
	Storage     metricValue       `json:"storage"`
	Load        []float64         `json:"load"`
	Network     map[string]uint64 `json:"networkBytes"`
}
type metricValue struct {
	Used    uint64 `json:"used"`
	Total   uint64 `json:"total"`
	Percent int    `json:"percent"`
}

func (s *server) metrics(w http.ResponseWriter, r *http.Request) {
	summary, _ := s.provider.SystemSummary(r.Context())
	load := []float64{}
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		for _, field := range fields[:minInt(3, len(fields))] {
			if value, err := strconv.ParseFloat(field, 64); err == nil {
				load = append(load, value)
			}
		}
	}
	network := map[string]uint64{}
	if data, err := os.Open("/proc/net/dev"); err == nil {
		defer data.Close()
		scanner := bufio.NewScanner(data)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if !strings.Contains(line, ":") {
				continue
			}
			parts := strings.SplitN(line, ":", 2)
			fields := strings.Fields(parts[1])
			if len(fields) >= 9 {
				rx, _ := strconv.ParseUint(fields[0], 10, 64)
				tx, _ := strconv.ParseUint(fields[8], 10, 64)
				name := strings.TrimSpace(parts[0])
				network[name+"Rx"] = rx
				network[name+"Tx"] = tx
			}
		}
	}
	respond(w, metricsSnapshot{CollectedAt: time.Now().UTC(), CPU: summary.CPUUsage, Memory: metricValue{Used: uint64(summary.MemoryUsed), Total: uint64(summary.MemoryTotal), Percent: percentage(summary.MemoryUsed, summary.MemoryTotal)}, Storage: metricValue{Used: uint64(summary.StorageUsed), Total: uint64(summary.StorageTotal), Percent: percentage(summary.StorageUsed, summary.StorageTotal)}, Load: load, Network: network}, nil, http.StatusOK)
}
func percentage(used, total int64) int {
	if total <= 0 {
		return 0
	}
	return int((used * 100) / total)
}
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
