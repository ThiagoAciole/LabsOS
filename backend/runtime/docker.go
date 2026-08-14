package runtime

import (
	"context"
	"encoding/json"
	"os/exec"
)

type CommandRunner func(context.Context, ...string) ([]byte, error)

type ContainerRuntime interface {
	Available(context.Context) (bool, error)
	Version(context.Context) (string, error)
}

type dockerRuntime struct{ run CommandRunner }

func NewDockerRuntime(run ...CommandRunner) ContainerRuntime {
	if len(run) > 0 {
		return &dockerRuntime{run: run[0]}
	}
	return &dockerRuntime{run: func(ctx context.Context, args ...string) ([]byte, error) {
		return exec.CommandContext(ctx, "docker", args...).Output()
	}}
}

func (d *dockerRuntime) Available(ctx context.Context) (bool, error) {
	_, err := d.run(ctx, "version", "--format", `{{json .}}`)
	return err == nil, nil
}

func (d *dockerRuntime) Version(ctx context.Context) (string, error) {
	output, err := d.run(ctx, "version", "--format", `{{json .}}`)
	if err != nil {
		return "", err
	}
	var value struct {
		Client struct {
			Version string `json:"Version"`
		} `json:"Client"`
	}
	if err := json.Unmarshal(output, &value); err != nil {
		return "", err
	}
	return value.Client.Version, nil
}
