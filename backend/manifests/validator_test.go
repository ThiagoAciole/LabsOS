package manifests

import "testing"

func TestValidateRejectsPrivilegedAndHostDockerSocket(t *testing.T) {
	for _, manifest := range []Manifest{{Privileged: true}, {HostPID: true}, {HostIPC: true}, {Capabilities: []string{"SYS_ADMIN"}}, {Volumes: []string{"/var/run/docker.sock:/var/run/docker.sock"}}} {
		if err := Validate(manifest); err == nil {
			t.Fatalf("unsafe manifest accepted: %+v", manifest)
		}
	}
}

func TestValidateAcceptsAppDataUnderAppsRoot(t *testing.T) {
	if err := Validate(Manifest{Volumes: []string{"/DATA/Apps/jellyfin:/config"}}); err != nil {
		t.Fatal(err)
	}
}

func TestJellyfinManifestIsSafeByDefault(t *testing.T) {
	if err := Validate(Jellyfin()); err != nil {
		t.Fatal(err)
	}
}

func TestSyncthingManifestIsSafeByDefault(t *testing.T) {
	if err := Validate(Syncthing()); err != nil { t.Fatal(err) }
}
