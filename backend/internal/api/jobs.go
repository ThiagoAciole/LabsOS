package api

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"labsos/backend/internal/platform"
)

type managedJob struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Target    string    `json:"target,omitempty"`
	Status    string    `json:"status"`
	Progress  int       `json:"progress"`
	Message   string    `json:"message"`
	Logs      []string  `json:"logs,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	cancel    context.CancelFunc
}

type jobEngine struct {
	mu      sync.RWMutex
	jobs    map[string]*managedJob
	seq     atomic.Uint64
	publish func(string)
	path    string
}

func newJobEngine() *jobEngine {
	path := os.Getenv("LABSOS_JOBS_FILE")
	if path == "" {
		path = "/tmp/labsos-jobs.json"
	}
	e := &jobEngine{jobs: make(map[string]*managedJob), path: path}
	e.load()
	return e
}

func (e *jobEngine) load() {
	data, err := os.ReadFile(e.path)
	if err != nil {
		return
	}
	var jobs []*managedJob
	if json.Unmarshal(data, &jobs) != nil {
		return
	}
	for _, job := range jobs {
		if job == nil || job.ID == "" {
			continue
		}
		if job.Status == "queued" || job.Status == "running" {
			job.Status = "error"
			job.Progress = 100
			job.Error = "processo reiniciado durante a operação"
			job.Message = "Operação interrompida após reinício da API"
		}
		job.cancel = nil
		e.jobs[job.ID] = job
		if value, err := strconv.ParseUint(strings.TrimPrefix(job.ID, "job-"), 10, 64); err == nil {
			for {
				current := e.seq.Load()
				if value <= current || e.seq.CompareAndSwap(current, value) {
					break
				}
			}
		}
	}
}

func (e *jobEngine) persistLocked() {
	items := make([]*managedJob, 0, len(e.jobs))
	for _, job := range e.jobs {
		copy := *job
		copy.cancel = nil
		copy.Logs = append([]string(nil), job.Logs...)
		items = append(items, &copy)
	}
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return
	}
	if err := os.MkdirAll(filepath.Dir(e.path), 0750); err != nil {
		return
	}
	tmp, err := os.CreateTemp(filepath.Dir(e.path), ".jobs-*.tmp")
	if err != nil {
		return
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err == nil {
		_, err = tmp.Write(data)
	}
	if closeErr := tmp.Close(); err == nil {
		err = closeErr
	}
	if err == nil {
		_ = os.Chmod(tmpPath, 0600)
		_ = os.Rename(tmpPath, e.path)
	}
}

func (e *jobEngine) setPublisher(publish func(string)) {
	e.mu.Lock()
	e.publish = publish
	e.mu.Unlock()
}

func (e *jobEngine) create(kind, target string) *managedJob {
	now := time.Now().UTC()
	job := &managedJob{ID: "job-" + formatJobID(e.seq.Add(1)), Kind: kind, Target: target, Status: "queued", Message: "Aguardando execução", CreatedAt: now, UpdatedAt: now}
	e.mu.Lock()
	e.jobs[job.ID] = job
	e.persistLocked()
	e.mu.Unlock()
	return job
}

func (e *jobEngine) get(id string) (*managedJob, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	job, ok := e.jobs[id]
	if !ok {
		return nil, false
	}
	copy := *job
	copy.Logs = append([]string(nil), job.Logs...)
	return &copy, true
}
func (e *jobEngine) list() []*managedJob {
	e.mu.RLock()
	defer e.mu.RUnlock()
	result := make([]*managedJob, 0, len(e.jobs))
	for _, job := range e.jobs {
		copy := *job
		copy.Logs = append([]string(nil), job.Logs...)
		result = append(result, &copy)
	}
	return result
}
func (e *jobEngine) update(id string, fn func(*managedJob)) {
	e.mu.Lock()
	var message string
	var publish func(string)
	if job := e.jobs[id]; job != nil {
		fn(job)
		job.UpdatedAt = time.Now().UTC()
		e.persistLocked()
		message = job.ID + " " + job.Status + " " + formatJobID(uint64(job.Progress)) + "%: " + job.Message
		publish = e.publish
	}
	e.mu.Unlock()
	if publish != nil && message != "" {
		publish(message)
	}
}

func (e *jobEngine) log(id, message string) {
	if message == "" {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if job := e.jobs[id]; job != nil {
		job.Logs = append(job.Logs, time.Now().UTC().Format(time.RFC3339)+" "+message)
		if len(job.Logs) > 100 {
			job.Logs = job.Logs[len(job.Logs)-100:]
		}
		job.UpdatedAt = time.Now().UTC()
		e.persistLocked()
	}
}

func (e *jobEngine) run(ctx context.Context, job *managedJob, action func(context.Context) (platform.Job, error)) {
	jobCtx, cancel := context.WithCancel(ctx)
	e.update(job.ID, func(item *managedJob) {
		item.Status = "running"
		item.Progress = 10
		item.Message = "Iniciando operação"
	})
	e.log(job.ID, "Operação iniciada")
	e.mu.Lock()
	job.cancel = cancel
	e.mu.Unlock()
	e.update(job.ID, func(item *managedJob) {
		item.Progress = 25
		item.Message = "Executando operação"
	})
	e.log(job.ID, "Operação em execução")
	result, err := action(jobCtx)
	cancel()
	if err != nil {
		e.log(job.ID, "Erro: "+err.Error())
		e.update(job.ID, func(item *managedJob) {
			item.Status = "error"
			item.Progress = 100
			item.Error = err.Error()
			item.Message = "Operação falhou"
		})
		return
	}
	e.update(job.ID, func(item *managedJob) {
		item.Progress = 90
		item.Message = "Finalizando operação"
	})
	e.log(job.ID, "Operação finalizada; atualizando estado")
	e.update(job.ID, func(item *managedJob) { item.Status = "success"; item.Progress = 100; item.Message = result.Message })
	e.log(job.ID, result.Message)
}

func (e *jobEngine) cancel(id string) bool {
	job, ok := e.get(id)
	if !ok {
		return false
	}
	if job.cancel != nil {
		job.cancel()
		e.update(id, func(item *managedJob) {
			item.Status = "cancelled"
			item.Progress = 100
			item.Message = "Operação cancelada"
		})
		return true
	}
	return false
}

func formatJobID(value uint64) string {
	const digits = "0123456789"
	if value == 0 {
		return "0"
	}
	result := ""
	for value > 0 {
		result = string(digits[value%10]) + result
		value /= 10
	}
	return result
}
