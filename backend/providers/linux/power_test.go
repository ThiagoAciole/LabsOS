package linux

import (
	"context"
	"testing"

	"labsos/backend/internal/platform"
)

func TestPowerUsesSystemAction(t *testing.T) {
	var got string
	p := &Provider{power: func(_ context.Context, action string) error { got = action; return nil }}
	job, err := p.Power(context.Background(), "reboot")
	if err != nil {
		t.Fatal(err)
	}
	if got != "reboot" || job.Status != "accepted" {
		t.Fatalf("power = %#v, action=%q", job, got)
	}
}

func TestPowerRejectsUnknownAction(t *testing.T) {
	p := &Provider{power: func(context.Context, string) error { t.Fatal("runner called"); return nil }}
	if _, err := p.Power(context.Background(), "format"); err != platform.ErrUnsupported {
		t.Fatalf("err = %v, want %v", err, platform.ErrUnsupported)
	}
}
