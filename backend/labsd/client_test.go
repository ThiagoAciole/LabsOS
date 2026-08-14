package labsd

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"testing"
)

func TestClientCall(t *testing.T) {
	socket := filepath.Join(t.TempDir(), "labsd.sock")
	server := New(socket, func(context.Context, ...string) error { return nil })
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = server.ListenAndServe(ctx) }()
	var conn net.Conn
	var err error
	for i := 0; i < 20; i++ {
		conn, err = net.Dial("unix", socket)
		if err == nil {
			conn.Close()
			break
		}
		_ = os.ErrClosed
	}
	if err != nil {
		t.Fatal(err)
	}
	if err := (Client{Socket: socket}).Call(context.Background(), "StartApp", "jellyfin"); err != nil {
		t.Fatal(err)
	}
}
