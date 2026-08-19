package api

import "testing"

func TestParseHumanBytes(t *testing.T) {
	if got := parseHumanBytes("1.5GiB"); got != 1610612736 {
		t.Fatalf("unexpected parser result for compact unit: %d", got)
	}
	if got := parseHumanBytes("2 MiB"); got != 2*1024*1024 {
		t.Fatalf("got %d", got)
	}
}
