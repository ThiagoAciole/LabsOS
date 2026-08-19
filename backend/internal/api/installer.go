package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"labsos/backend/internal/safety"
)

type installerDisk struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Path          string `json:"path"`
	SizeBytes     uint64 `json:"sizeBytes"`
	Size          string `json:"size"`
	Kind          string `json:"kind"`
	Ready         bool   `json:"ready"`
	HasLabsOSData bool   `json:"hasLabsOSData"`
}

type installerJob struct {
	ID         string `json:"id"`
	Status     string `json:"status"`
	Progress   int    `json:"progress"`
	Message    string `json:"message"`
	Error      string `json:"error,omitempty"`
	Simulated  bool   `json:"simulated"`
	operation  string
	disk       string
	preserve   bool
	serverName string
	password   string
	cancel     chan struct{}
}

type installerStartRequest struct {
	Operation  string `json:"operation"`
	Disk       string `json:"disk"`
	Preserve   bool   `json:"preserve"`
	ServerName string `json:"serverName"`
	Password   string `json:"password"`
}

var installerSequence atomic.Uint64

func (s *server) installerStatus(w http.ResponseWriter, _ *http.Request) {
	s.installerMu.Lock()
	job := s.installerJob
	var snapshot *installerJob
	if job != nil {
		value := installerJobSnapshot(job)
		snapshot = &value
	}
	s.installerMu.Unlock()
	status := map[string]any{"status": "needs-install"}
	if snapshot != nil {
		status["status"] = snapshot.Status
		status["job"] = snapshot
	}
	respond(w, status, nil, http.StatusOK)
}

func (s *server) installerDisks(w http.ResponseWriter, _ *http.Request) {
	if disks, err := discoverInstallerDisks(); err == nil {
		respond(w, disks, nil, http.StatusOK)
		return
	}
	writeError(w, http.StatusServiceUnavailable, "DISK_DISCOVERY_FAILED", "lsblk is unavailable")
}

func (s *server) installerValidate(w http.ResponseWriter, r *http.Request) {
	var input installerStartRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must be a JSON object")
		return
	}
	if input.Operation != "fresh" && input.Operation != "restore" && input.Operation != "transfer" {
		writeError(w, http.StatusBadRequest, "INVALID_OPERATION", "unsupported installer operation")
		return
	}
	if input.Disk == "" {
		writeError(w, http.StatusBadRequest, "INVALID_DISK", "selected disk is not available")
		return
	}
	if input.ServerName == "" || len(input.Password) < 4 {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "serverName and a password with at least 4 characters are required")
		return
	}
	if safety.RealOperationsEnabled() {
		if err := validateRealInstallerDisk(input.Disk); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INSTALLER_DISK", err.Error())
			return
		}
	}
	respond(w, map[string]any{"valid": true, "destructive": !input.Preserve, "disk": input.Disk}, nil, http.StatusOK)
}

func (s *server) installerStart(w http.ResponseWriter, r *http.Request) {
	var input installerStartRequest
	if err := decodeJSON(r, &input); err != nil {
		writeError(w, http.StatusBadRequest, "INVALID_JSON", "request body must be a JSON object")
		return
	}
	if input.Operation != "fresh" && input.Operation != "restore" && input.Operation != "transfer" {
		writeError(w, http.StatusBadRequest, "INVALID_OPERATION", "unsupported installer operation")
		return
	}
	if input.Disk == "" || input.ServerName == "" || len(input.Password) < 4 {
		writeError(w, http.StatusBadRequest, "INVALID_REQUEST", "disk, serverName and a password with at least 4 characters are required")
		return
	}
	s.installerMu.Lock()
	if s.installerJob != nil && (s.installerJob.Status == "running" || s.installerJob.Status == "queued") {
		s.installerMu.Unlock()
		writeError(w, http.StatusConflict, "INSTALLER_BUSY", "an installer operation is already running")
		return
	}
	s.installerMu.Unlock()
	if safety.RealOperationsEnabled() {
		if err := validateRealInstallerDisk(input.Disk); err != nil {
			writeError(w, http.StatusBadRequest, "INVALID_INSTALLER_DISK", err.Error())
			return
		}
	}
	id := "install-" + strconv.FormatUint(installerSequence.Add(1), 10)
	job := &installerJob{ID: id, Status: "running", Message: "Preparando instalação", Simulated: !safety.RealOperationsEnabled(), operation: input.Operation, disk: input.Disk, preserve: input.Preserve, serverName: input.ServerName, password: input.Password, cancel: make(chan struct{})}
	s.installerMu.Lock()
	s.installerJob = job
	snapshot := installerJobSnapshot(job)
	s.installerMu.Unlock()
	s.eventHub.publish("installer", id+" started")
	go s.runInstallerJob(job)
	respond(w, snapshot, nil, http.StatusAccepted)
}

func (s *server) runInstallerJob(job *installerJob) {
	if !job.Simulated {
		if err := executeRealInstaller(job); err != nil {
			s.installerMu.Lock()
			job.Status = "error"
			job.Progress = 100
			job.Error = err.Error()
			job.Message = "Instalação real falhou"
			s.installerMu.Unlock()
			s.eventHub.publish("installer", job.ID+" error")
			return
		}
	}
	steps := []string{"Validando disco (somente leitura)", "Simulando preparação de partições", "Simulando instalação do root filesystem", "Configurando serviços", "Preparando first boot"}
	if safety.RealOperationsEnabled() {
		steps = []string{"Validando disco", "Preparando instalação protegida", "Executando operação autorizada", "Configurando serviços", "Preparando first boot"}
	}
	for index, message := range steps {
		timer := time.NewTimer(250 * time.Millisecond)
		select {
		case <-job.cancel:
			timer.Stop()
			s.installerMu.Lock()
			job.Status = "cancelled"
			job.Progress = (index * 20)
			job.Message = "Instalação cancelada"
			s.installerMu.Unlock()
			s.eventHub.publish("installer", job.ID+" cancelled")
			return
		case <-timer.C:
		}
		s.installerMu.Lock()
		job.Progress = (index + 1) * 20
		job.Message = message
		progress := job.Progress
		s.installerMu.Unlock()
		s.eventHub.publish("installer", job.ID+" progress "+strconv.Itoa(progress))
	}
	s.installerMu.Lock()
	job.Status = "complete"
	job.Progress = 100
	job.Message = "Configuração concluída"
	s.installerMu.Unlock()
	s.eventHub.publish("installer", job.ID+" complete")
}

// executeRealInstaller delegates the destructive part to an administrator-owned
// executable. The API never constructs shell commands or performs partitioning
// itself. The executable must be explicitly configured and operations must have
// passed the safety policy before this function can be reached.
func executeRealInstaller(job *installerJob) error {
	script := strings.TrimSpace(os.Getenv("LABSOS_INSTALLER_SCRIPT"))
	if script == "" {
		script = "/usr/lib/labsos/labsos-installer"
	}
	info, err := os.Stat(script)
	if err != nil || info.IsDir() || info.Mode().Perm()&0111 == 0 {
		return fmt.Errorf("installer script is missing or not executable")
	}
	tmpDir, err := os.MkdirTemp(os.TempDir(), "labsos-installer-")
	if err != nil {
		return fmt.Errorf("create protected installer directory: %w", err)
	}
	defer os.RemoveAll(tmpDir)
	passwordFile, err := os.CreateTemp(tmpDir, "password-")
	if err != nil {
		return fmt.Errorf("create protected installer input: %w", err)
	}
	path := passwordFile.Name()
	defer os.Remove(path)
	if err := passwordFile.Chmod(0600); err != nil {
		passwordFile.Close()
		return fmt.Errorf("protect installer input: %w", err)
	}
	if _, err := passwordFile.WriteString(job.password + "\n"); err != nil {
		passwordFile.Close()
		return fmt.Errorf("write installer input: %w", err)
	}
	if err := passwordFile.Close(); err != nil {
		return fmt.Errorf("close installer input: %w", err)
	}
	ctx := context.Background()
	preserve := "false"
	if job.preserve {
		preserve = "true"
	}
	command := exec.CommandContext(ctx, script, "--job", job.ID, "--operation", job.operation, "--disk", job.disk, "--preserve", preserve, "--server-name", job.serverName, "--password-file", path)
	if output, err := command.CombinedOutput(); err != nil {
		return fmt.Errorf("installer script failed: %s", strings.TrimSpace(string(output)))
	}
	return nil
}

func (s *server) installerCancel(w http.ResponseWriter, r *http.Request) {
	s.installerMu.Lock()
	job := s.installerJob
	if job == nil || job.ID != r.PathValue("id") {
		s.installerMu.Unlock()
		writeError(w, http.StatusNotFound, "NOT_FOUND", "installer job not found")
		return
	}
	if job.Status != "running" {
		s.installerMu.Unlock()
		respond(w, map[string]any{"cancelled": false, "status": job.Status}, nil, http.StatusConflict)
		return
	}
	job.Status = "cancelling"
	close(job.cancel)
	s.installerMu.Unlock()
	s.auditStore.add("installer.cancel", job.ID, "accepted", "")
	respond(w, map[string]any{"cancelled": true, "id": job.ID}, nil, http.StatusAccepted)
}

func (s *server) installerJobStatus(w http.ResponseWriter, r *http.Request) {
	s.installerMu.Lock()
	job := s.installerJob
	s.installerMu.Unlock()
	if job == nil || job.ID != r.PathValue("id") {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "installer job not found")
		return
	}
	s.installerMu.Lock()
	snapshot := installerJobSnapshot(job)
	s.installerMu.Unlock()
	respond(w, snapshot, nil, http.StatusOK)
}

func (s *server) installerEvents(w http.ResponseWriter, _ *http.Request) {
	s.installerMu.Lock()
	job := s.installerJob
	s.installerMu.Unlock()
	if job == nil {
		respond(w, []any{}, nil, http.StatusOK)
		return
	}
	s.installerMu.Lock()
	snapshot := installerJobSnapshot(job)
	s.installerMu.Unlock()
	respond(w, []map[string]any{{"type": "progress", "jobId": snapshot.ID, "progress": snapshot.Progress, "message": snapshot.Message}}, nil, http.StatusOK)
}

func installerJobSnapshot(job *installerJob) installerJob {
	snapshot := *job
	snapshot.cancel = nil
	return snapshot
}

func (s *server) installerReboot(w http.ResponseWriter, r *http.Request) {
	if !safety.RealOperationsEnabled() {
		message := "reboot planejado; operações reais estão desabilitadas"
		s.auditStore.add("installer.reboot", "labsos", "planned", message)
		s.eventHub.publish("installer", message)
		respond(w, map[string]any{"accepted": false, "planned": true, "simulated": true, "message": message}, nil, http.StatusAccepted)
		return
	}
	job, err := s.provider.Power(r.Context(), "reboot")
	if err != nil {
		s.auditStore.add("installer.reboot", "labsos", "error", err.Error())
		respond(w, nil, err, http.StatusServiceUnavailable)
		return
	}
	s.auditStore.add("installer.reboot", "labsos", "accepted", job.Message)
	respond(w, map[string]any{"accepted": true, "planned": false, "simulated": false, "message": job.Message}, nil, http.StatusAccepted)
}

func discoverInstallerDisks() ([]installerDisk, error) {
	out, err := exec.Command("lsblk", "-b", "-dn", "-o", "NAME,SIZE,TYPE,MODEL,RO").Output()
	if err != nil {
		return nil, err
	}
	var result []installerDisk
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		fields := strings.Fields(line)
		if len(fields) < 5 || fields[2] != "disk" {
			continue
		}
		size, _ := strconv.ParseUint(fields[1], 10, 64)
		ready := size >= 20*1024*1024*1024 && fields[4] == "0"
		model := strings.TrimSpace(strings.Join(fields[3:len(fields)-1], " "))
		if model == "" {
			model = "Linux disk"
		}
		result = append(result, installerDisk{ID: fields[0], Name: model, Path: "/dev/" + fields[0], SizeBytes: size, Size: formatBytes(size), Kind: "disk", Ready: ready})
	}
	return result, nil
}

func validateRealInstallerDisk(path string) error {
	if !strings.HasPrefix(path, "/dev/") || strings.Contains(path, "..") {
		return fmt.Errorf("disk must be an absolute device path")
	}
	disks, err := discoverInstallerDisks()
	if err != nil {
		return fmt.Errorf("unable to discover disks: %w", err)
	}
	for _, disk := range disks {
		if disk.Path == path {
			if !disk.Ready {
				return fmt.Errorf("disk is read-only, too small, or unavailable")
			}
			return nil
		}
	}
	return fmt.Errorf("disk was not returned by lsblk")
}

func formatBytes(value uint64) string {
	for _, unit := range []string{"Bytes", "KB", "MB", "GB", "TB"} {
		if value < 1024 || unit == "TB" {
			if unit == "Bytes" {
				return strconv.FormatUint(value, 10) + " " + unit
			}
			return strconv.FormatFloat(float64(value), 'g', 3, 64) + " " + unit
		}
		value /= 1024
	}
	return "0 Bytes"
}

func decodeJSON(r *http.Request, target any) error {
	return json.NewDecoder(io.LimitReader(r.Body, 1<<20)).Decode(target)
}
