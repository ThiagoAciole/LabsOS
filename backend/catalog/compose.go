package catalog

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type composeDocument struct {
	Services map[string]struct {
		Privileged  bool     `yaml:"privileged"`
		NetworkMode string   `yaml:"network_mode"`
		PID         string   `yaml:"pid"`
		IPC         string   `yaml:"ipc"`
		CapAdd      []string `yaml:"cap_add"`
		Volumes     []composeVolume `yaml:"volumes"`
		Devices     []composeVolume `yaml:"devices"`
	} `yaml:"services"`
}

type composeVolume struct {
	Source string
}

func (v *composeVolume) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		v.Source = strings.SplitN(value.Value, ":", 2)[0]
		return nil
	}
	var long struct {
		Source string `yaml:"source"`
	}
	if err := value.Decode(&long); err != nil {
		return err
	}
	v.Source = long.Source
	return nil
}

func ImportCompose(ctx context.Context, source string) (App, error) {
	content, name, err := resolveCompose(ctx, source)
	if err != nil {
		return App{}, err
	}
	var document composeDocument
	if err := yaml.Unmarshal(content, &document); err != nil {
		return App{}, fmt.Errorf("invalid compose: %w", err)
	}
	if len(document.Services) == 0 {
		return App{}, fmt.Errorf("compose has no services")
	}
	for service, value := range document.Services {
		if value.Privileged || value.NetworkMode == "host" || value.PID == "host" || value.IPC == "host" || len(value.CapAdd) > 0 {
			return App{}, fmt.Errorf("unsafe compose service %q", service)
		}
		for _, volume := range append(value.Volumes, value.Devices...) {
			if err := validateComposePath(volume.Source); err != nil {
				return App{}, fmt.Errorf("service %q: %w", service, err)
			}
		}
	}
	id := normalizeComposeID(name)
	return App{ID: id, Name: name, Description: "Docker Compose application", Source: source, Installable: true, Compose: string(content)}, nil
}

func resolveCompose(ctx context.Context, source string) ([]byte, string, error) {
	if info, err := os.Stat(source); err == nil {
		if info.IsDir() {
			source = findCompose(source)
		} else {
			content, readErr := os.ReadFile(source)
			return content, composeName(source), readErr
		}
	}
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		if strings.HasSuffix(source, ".git") || strings.Contains(source, "github.com/") && !strings.Contains(filepath.Base(source), ".yml") && !strings.Contains(filepath.Base(source), ".yaml") {
			dir, err := os.MkdirTemp("", "labsos-compose-")
			if err != nil {
				return nil, "", err
			}
			if err := exec.CommandContext(ctx, "git", "clone", "--depth", "1", source, dir).Run(); err != nil {
				return nil, "", err
			}
			source = findCompose(dir)
		} else {
			request, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
			if err != nil {
				return nil, "", err
			}
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				return nil, "", err
			}
			defer response.Body.Close()
			if response.StatusCode != http.StatusOK {
				return nil, "", fmt.Errorf("compose source returned %s", response.Status)
			}
			content, err := io.ReadAll(io.LimitReader(response.Body, 2<<20))
			return content, composeName(source), err
		}
	}
	content, err := os.ReadFile(source)
	return content, composeName(source), err
}

func findCompose(root string) string {
	var found string
	_ = filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && (entry.Name() == "compose.yml" || entry.Name() == "docker-compose.yml") && found == "" {
			found = path
		}
		return nil
	})
	return found
}

func validateComposePath(value string) error {
	path := strings.SplitN(value, ":", 2)[0]
	if strings.Contains(path, "..") || strings.Contains(path, "docker.sock") {
		return fmt.Errorf("unsafe host path")
	}
	if strings.HasPrefix(path, "/") && !strings.HasPrefix(path, "/DATA/Apps/") && !strings.HasPrefix(path, "/DATA/Media/") && !strings.HasPrefix(path, "/DATA/AppData/") && path != "/DATA/Apps" && path != "/DATA/Media" && path != "/DATA/AppData" {
		return fmt.Errorf("host path outside /DATA")
	}
	return nil
}

func composeName(source string) string {
	base := filepath.Base(source)
	if base == "compose.yml" || base == "docker-compose.yml" { return filepath.Base(filepath.Dir(source)) }
	return strings.TrimSuffix(strings.TrimSuffix(base, ".yml"), ".yaml")
}
func normalizeComposeID(value string) string {
	value = strings.ToLower(value)
	value = strings.NewReplacer(" ", "-", "_", "-", ".", "-").Replace(value)
	return value
}
