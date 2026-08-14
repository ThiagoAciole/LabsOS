package mock_test

import (
	"context"
	"testing"

	"labsos/backend/providers/mock"
)

func TestAppActionsChangeProductStateAndCreateJob(t *testing.T) {
	p := mock.New()
	job, err := p.AppAction(context.Background(), "jellyfin", "stop")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "success" {
		t.Fatalf("job status = %q", job.Status)
	}
	apps, err := p.Apps(context.Background(), false)
	if err != nil {
		t.Fatal(err)
	}
	if apps[0].ID != "jellyfin" || apps[0].Status != "stopped" {
		t.Fatalf("unexpected app state: %+v", apps[0])
	}
}

func TestSettingsReturnsDefensiveCopy(t *testing.T) {
	p := mock.New()
	settings, err := p.Settings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	settings["hostname"] = "mutated-outside-provider"
	again, err := p.Settings(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if again["hostname"] == "mutated-outside-provider" {
		t.Fatal("settings leaked mutable provider state")
	}
}

func TestPowerIsSimulated(t *testing.T) {
	p := mock.New()
	job, err := p.Power(context.Background(), "reboot")
	if err != nil {
		t.Fatal(err)
	}
	if job.Status != "success" || job.Message != "reboot simulated in mock mode" {
		t.Fatalf("unexpected job: %+v", job)
	}
}
