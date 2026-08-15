package labsd

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"labsos/backend/catalog"
)

const socketPath = "/run/labsos/labsd.sock"

type Runner func(context.Context, ...string) error
type ListRunner func(context.Context, ...string) ([]string, error)

type Server struct {
	Socket string
	run    Runner
	list   ListRunner
	status func(context.Context, string) (string, error)
}

func New(socket string, run Runner) *Server {
	if socket == "" {
		socket = socketPath
	}
	if run == nil {
		run = func(ctx context.Context, args ...string) error {
			return exec.CommandContext(ctx, "docker", args...).Run()
		}
	}
	list := func(ctx context.Context, args ...string) ([]string, error) {
		output, err := exec.CommandContext(ctx, "docker", args...).Output()
		if err != nil {
			return nil, err
		}
		return strings.Fields(string(output)), nil
	}
	status := func(ctx context.Context, app string) (string, error) {
		output, err := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.State.Status}}", app).Output()
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(output)), nil
	}
	return &Server{Socket: socket, run: run, list: list, status: status}
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	if err := os.MkdirAll(filepath.Dir(s.Socket), 0755); err != nil {
		return err
	}
	_ = os.Remove(s.Socket)
	listener, err := net.Listen("unix", s.Socket)
	if err != nil {
		return err
	}
	defer listener.Close()
	if err := os.Chmod(s.Socket, 0660); err != nil {
		return err
	}
	go func() { <-ctx.Done(); listener.Close() }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handle(ctx, conn)
	}
}

func (s *Server) handle(ctx context.Context, conn net.Conn) {
	defer conn.Close()
	var request Request
	response := Response{OK: false}
	if err := json.NewDecoder(conn).Decode(&request); err != nil {
		response.Message = "invalid request"
	} else if request.Operation == "ListApps" {
		apps, err := s.list(ctx, "ps", "-a", "--format", "{{.Names}}")
		if err != nil {
			response.Message = err.Error()
		} else {
			response.OK = true
			response.Apps = apps
		}
	} else if request.Operation == "LogsApp" {
		logs, err := s.logs(ctx, request)
		if err != nil {
			response.Message = err.Error()
		} else {
			response.OK = true
			response.Logs = logs
		}
	} else if status, err := s.dispatch(ctx, request); err != nil {
		response.Message = err.Error()
	} else {
		response.OK = true
		response.Status = status
	}
	_ = json.NewEncoder(conn).Encode(response)
}

func (s *Server) logs(ctx context.Context, request Request) (string, error) {
	if !safeAppID(request.App) {
		return "", fmt.Errorf("invalid app")
	}
	lines := request.Lines
	if lines < 1 {
		lines = 100
	}
	if lines > 1000 {
		lines = 1000
	}
	output, err := exec.CommandContext(ctx, "docker", "compose", "-f", "/opt/labsos/apps/"+request.App+"/compose.yaml", "logs", "--no-color", "--tail", strconv.Itoa(lines)).Output()
	return string(output), err
}

func (s *Server) dispatch(ctx context.Context, request Request) (string, error) {
	if !safeAppID(request.App) {
		return "", fmt.Errorf("invalid app")
	}
	compose := "/opt/labsos/apps/" + request.App + "/compose.yaml"
	switch request.Operation {
	case "EnsureApp", "InstallApp", "StartApp":
		if request.Operation != "StartApp" {
			if request.Compose != "" {
				if err := s.installCompose(ctx, request.App, request.Compose); err != nil {
					return "", err
				}
			}
			if err := s.run(ctx, "compose", "-f", compose, "up", "-d"); err != nil {
				return "", err
			}
		} else if err := s.run(ctx, "compose", "-f", compose, "start"); err != nil {
			return "", err
		}
	case "StopApp":
		return "", s.run(ctx, "compose", "-f", compose, "stop")
	case "RestartApp":
		return "", s.run(ctx, "compose", "-f", compose, "restart")
	case "RemoveApp":
		return "", s.run(ctx, "compose", "-f", compose, "down")
	case "StatusApp":
		return s.status(ctx, request.App)
	default:
		return "", fmt.Errorf("unsupported operation")
	}
	return "", nil
}

func (s *Server) installCompose(ctx context.Context, app, content string) error {
	dir := "/opt/labsos/apps/" + app
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(dir, ".compose-*.yaml")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err = tmp.WriteString(content); err != nil {
		tmp.Close()
		return err
	}
	if err = tmp.Close(); err != nil {
		return err
	}
	if _, err = catalog.ImportCompose(ctx, tmpPath); err != nil {
		return err
	}
	return os.Rename(tmpPath, filepath.Join(dir, "compose.yaml"))
}

func safeAppID(id string) bool {
	if id == "" || id == "." || id == ".." || strings.Contains(id, "..") {
		return false
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '.' {
			return false
		}
	}
	return true
}
