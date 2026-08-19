package labsd

import (
	"context"
	"reflect"
	"testing"
)

func TestDispatchUsesFixedDockerOperations(t *testing.T) {
	var got []string
	s := New("", func(_ context.Context, args ...string) error { got = args; return nil })
	if _, err := s.dispatch(context.Background(), Request{Operation: "InstallApp", App: "jellyfin"}); err != nil {
		t.Fatal(err)
	}
	want := []string{"compose", "-f", "/opt/labsos/apps/jellyfin/compose.yaml", "up", "-d"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("args = %#v, want %#v", got, want)
	}
}

func TestDispatchUpdatesComposeImagesAndRecreatesApp(t *testing.T) {
	var calls [][]string
	s := New("", func(_ context.Context, args ...string) error {
		calls = append(calls, append([]string(nil), args...))
		return nil
	})
	if _, err := s.dispatch(context.Background(), Request{Operation: "UpdateApp", App: "jellyfin"}); err != nil {
		t.Fatal(err)
	}
	want := [][]string{
		{"compose", "-f", "/opt/labsos/apps/jellyfin/compose.yaml", "pull"},
		{"compose", "-f", "/opt/labsos/apps/jellyfin/compose.yaml", "up", "-d", "--remove-orphans"},
	}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestSafeAppIDAllowsCasaOSDots(t *testing.T) {
	if !safeAppID("org.icewhale.2fauth") {
		t.Fatal("CasaOS app id with internal dots was rejected")
	}
	if safeAppID("../escape") || safeAppID("a..b") {
		t.Fatal("unsafe app id was accepted")
	}
}
