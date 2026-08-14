package manifests

import (
	"fmt"
	"strings"
)

type Manifest struct {
	Privileged  bool
	HostNetwork bool
	HostPID     bool
	HostIPC     bool
	Capabilities []string
	Volumes     []string
	Devices     []string
}

func Jellyfin() Manifest {
	return Manifest{Volumes: []string{"/DATA/Apps/jellyfin:/config", "/DATA/Media:/media:ro"}}
}

func Syncthing() Manifest {
	return Manifest{Volumes: []string{"/DATA/Apps/syncthing:/var/syncthing"}}
}

func Validate(manifest Manifest) error {
	if manifest.Privileged {
		return fmt.Errorf("privileged containers are not supported")
	}
	if manifest.HostNetwork {
		return fmt.Errorf("host network is not supported")
	}
	if manifest.HostPID || manifest.HostIPC || len(manifest.Capabilities) > 0 {
		return fmt.Errorf("host namespaces and extra capabilities are not supported")
	}
	for _, volume := range append(append([]string{}, manifest.Volumes...), manifest.Devices...) {
		parts := strings.SplitN(volume, ":", 2)
		allowedDataPath := parts[0] == "/DATA/Apps" || parts[0] == "/DATA/Media" || strings.HasPrefix(parts[0], "/DATA/Apps/") || strings.HasPrefix(parts[0], "/DATA/Media/")
		if strings.Contains(parts[0], "docker.sock") || strings.Contains(parts[0], "..") || (strings.HasPrefix(parts[0], "/") && !allowedDataPath) {
			return fmt.Errorf("host path is not allowed: %s", parts[0])
		}
	}
	return nil
}
