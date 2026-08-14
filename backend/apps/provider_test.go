package apps

import "testing"

func TestProductAppStatusesStayDockerAgnostic(t *testing.T) {
	if AppStatusRunning != "running" || AppStatusInstalling != "installing" {
		t.Fatal("unexpected product statuses")
	}
}
