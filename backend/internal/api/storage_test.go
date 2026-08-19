package api

import "testing"

func TestParseDF(t *testing.T) {
	usage, ok := parseDF("Mounted on Used Avail Size Use%\n/DATA 100 900 1000 10%\n")
	if !ok || usage.UsedBytes != 100 || usage.AvailableBytes != 900 || usage.TotalBytes != 1000 || usage.UsePercent != 10 {
		t.Fatalf("unexpected usage: %#v, ok=%v", usage, ok)
	}
}
