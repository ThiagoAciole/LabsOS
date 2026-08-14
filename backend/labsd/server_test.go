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
