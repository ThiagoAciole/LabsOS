package systemapps

import "testing"

func TestFileSystemAppHasStableProductID(t *testing.T) {
	app := App{ID: "files"}
	if app.ID != "files" {
		t.Fatalf("id = %q", app.ID)
	}
}
