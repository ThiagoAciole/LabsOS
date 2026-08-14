package linux

import (
	"fmt"
	"strconv"
	"strings"
)

type cpuCounters struct{ total, idle uint64 }

func parseCPUStat(input string) (cpuCounters, error) {
	fields := strings.Fields(strings.SplitN(input, "\n", 2)[0])
	if len(fields) < 5 || fields[0] != "cpu" {
		return cpuCounters{}, fmt.Errorf("invalid /proc/stat")
	}
	var values []uint64
	for _, field := range fields[1:] {
		value, err := strconv.ParseUint(field, 10, 64)
		if err != nil {
			return cpuCounters{}, fmt.Errorf("invalid /proc/stat: %w", err)
		}
		values = append(values, value)
	}
	var total uint64
	for _, value := range values {
		total += value
	}
	idle := values[3]
	if len(values) > 4 {
		idle += values[4]
	}
	return cpuCounters{total: total, idle: idle}, nil
}

func cpuUsage(before, after cpuCounters) float64 {
	total := after.total - before.total
	if total == 0 {
		return 0
	}
	idle := after.idle - before.idle
	return float64(total-idle) / float64(total) * 100
}

func parseMemInfo(input string) (used, total int64, err error) {
	values := map[string]int64{}
	for _, line := range strings.Split(input, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		value, parseErr := strconv.ParseInt(fields[1], 10, 64)
		if parseErr == nil {
			values[strings.TrimSuffix(fields[0], ":")] = value * 1024
		}
	}
	total = values["MemTotal"]
	available, ok := values["MemAvailable"]
	if total == 0 || !ok {
		return 0, 0, fmt.Errorf("invalid /proc/meminfo")
	}
	return total - available, total, nil
}

func parseUptime(input string) (int64, error) {
	fields := strings.Fields(input)
	if len(fields) == 0 {
		return 0, fmt.Errorf("invalid /proc/uptime")
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	return int64(seconds), err
}

func parseNetworkCounters(input string) (rx, tx uint64, err error) {
	found := false
	for _, line := range strings.Split(input, "\n") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 || strings.TrimSpace(parts[0]) == "lo" {
			continue
		}
		fields := strings.Fields(parts[1])
		if len(fields) < 9 {
			continue
		}
		received, rxErr := strconv.ParseUint(fields[0], 10, 64)
		transmitted, txErr := strconv.ParseUint(fields[8], 10, 64)
		if rxErr != nil || txErr != nil {
			continue
		}
		rx += received
		tx += transmitted
		found = true
	}
	if !found {
		return 0, 0, fmt.Errorf("no network interface in /proc/net/dev")
	}
	return rx, tx, nil
}
