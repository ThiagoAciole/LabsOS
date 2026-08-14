package linux

import "testing"

func TestParseCPUUsage(t *testing.T) {
	before, err := parseCPUStat("cpu  100 0 50 850 0 0 0 0 0 0\n")
	if err != nil {
		t.Fatal(err)
	}
	after, err := parseCPUStat("cpu  130 0 70 900 0 0 0 0 0 0\n")
	if err != nil {
		t.Fatal(err)
	}
	if got := cpuUsage(before, after); got != 50 {
		t.Fatalf("cpu usage = %v", got)
	}
}

func TestParseMemoryAndUptime(t *testing.T) {
	used, total, err := parseMemInfo("MemTotal: 8000000 kB\nMemAvailable: 5000000 kB\n")
	if err != nil {
		t.Fatal(err)
	}
	if used != 3_072_000_000 || total != 8_192_000_000 {
		t.Fatalf("memory = %d/%d", used, total)
	}
	uptime, err := parseUptime("483012.91 0.00\n")
	if err != nil || uptime != 483012 {
		t.Fatalf("uptime = %d, %v", uptime, err)
	}
}

func TestParseNetworkCounters(t *testing.T) {
	input := "Inter-| Receive | Transmit\n face |bytes packets|bytes packets\n eth0: 1200 1 0 0 0 0 0 0 3400 2 0 0 0 0 0 0\n lo: 9000 1 0 0 0 0 0 0 9000 1 0 0 0 0 0 0\n"
	rx, tx, err := parseNetworkCounters(input)
	if err != nil {
		t.Fatal(err)
	}
	if rx != 1200 || tx != 3400 {
		t.Fatalf("network = %d/%d", rx, tx)
	}
}
