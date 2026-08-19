package api

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"labsos/backend/internal/platform"
	"labsos/backend/internal/safety"
)

type scheduledTask struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Action    string    `json:"action"`
	Target    string    `json:"target"`
	Interval  int       `json:"intervalSeconds"`
	Enabled   bool      `json:"enabled"`
	LastRun   time.Time `json:"lastRun,omitempty"`
	NextRun   time.Time `json:"nextRun"`
	LastState string    `json:"lastState,omitempty"`
}

type schedulerStore struct {
	mu    sync.Mutex
	tasks map[string]scheduledTask
	seq   int
}

func newSchedulerStore() *schedulerStore {
	s := &schedulerStore{tasks: map[string]scheduledTask{}}
	if data, err := os.ReadFile(schedulerPath()); err == nil {
		_ = json.Unmarshal(data, &s.tasks)
	}
	for id := range s.tasks {
		if n, err := strconv.Atoi(id); err == nil && n >= s.seq {
			s.seq = n + 1
		}
	}
	return s
}
func schedulerPath() string {
	if value := os.Getenv("LABSOS_SCHEDULER_FILE"); value != "" {
		return value
	}
	return "/tmp/labsos-scheduler.json"
}
func (s *schedulerStore) save() {
	data, _ := json.MarshalIndent(s.tasks, "", "  ")
	_ = os.WriteFile(schedulerPath(), data, 0600)
}
func (s *schedulerStore) list() []scheduledTask {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]scheduledTask, 0, len(s.tasks))
	for _, task := range s.tasks {
		result = append(result, task)
	}
	return result
}
func (s *schedulerStore) tick(ctx context.Context, server *server) {
	s.mu.Lock()
	now := time.Now().UTC()
	for id, task := range s.tasks {
		if !task.Enabled || task.NextRun.After(now) {
			continue
		}
		task.LastRun, task.NextRun = now, now.Add(time.Duration(task.Interval)*time.Second)
		task.LastState = "queued"
		s.tasks[id] = task
		job := server.jobs.create("scheduled-"+task.Action, task.Target)
		if !safety.RealOperationsEnabled() {
			server.jobs.update(job.ID, func(item *managedJob) {
				item.Status = "success"
				item.Progress = 100
				item.Message = "Ação agendada; operações reais estão desativadas"
			})
			server.jobs.log(job.ID, "Execução bloqueada pela política de segurança")
			task.LastState = "planned"
			s.tasks[id] = task
			continue
		}
		go func(task scheduledTask, scheduledJob *managedJob) {
			server.jobs.run(ctx, scheduledJob, func(actionCtx context.Context) (platform.Job, error) {
				return server.provider.AppAction(actionCtx, task.Target, task.Action)
			})
			s.mu.Lock()
			defer s.mu.Unlock()
			current := s.tasks[task.ID]
			result, ok := server.jobs.get(scheduledJob.ID)
			if !ok || result.Status == "error" {
				current.LastState = "error"
			} else {
				current.LastState = "success"
			}
			s.tasks[task.ID] = current
			s.save()
		}(task, job)
	}
	s.save()
	s.mu.Unlock()
}

func (s *server) schedulerList(w http.ResponseWriter, _ *http.Request) {
	respond(w, s.scheduler.list(), nil, http.StatusOK)
}
func (s *server) schedulerCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name, Action, Target string
		Interval             int
		Enabled              *bool
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Action == "" || input.Target == "" || input.Interval < 5 {
		writeError(w, http.StatusBadRequest, "INVALID_SCHEDULE", "action, target and intervalSeconds >= 5 are required")
		return
	}
	s.scheduler.mu.Lock()
	id := strconv.Itoa(s.scheduler.seq)
	s.scheduler.seq++
	enabled := true
	if input.Enabled != nil {
		enabled = *input.Enabled
	}
	task := scheduledTask{ID: id, Name: input.Name, Action: input.Action, Target: input.Target, Interval: input.Interval, Enabled: enabled, NextRun: time.Now().UTC().Add(time.Duration(input.Interval) * time.Second), LastState: "scheduled"}
	s.scheduler.tasks[id] = task
	s.scheduler.save()
	s.scheduler.mu.Unlock()
	respond(w, task, nil, http.StatusCreated)
	s.auditStore.add("scheduler.create", task.ID, "success", task.Action+" "+task.Target)
}
func (s *server) schedulerDelete(w http.ResponseWriter, r *http.Request) {
	s.scheduler.mu.Lock()
	defer s.scheduler.mu.Unlock()
	if _, ok := s.scheduler.tasks[r.PathValue("id")]; !ok {
		writeError(w, http.StatusNotFound, "NOT_FOUND", "scheduled task not found")
		return
	}
	delete(s.scheduler.tasks, r.PathValue("id"))
	s.scheduler.save()
	respond(w, map[string]any{"deleted": true, "id": r.PathValue("id")}, nil, http.StatusOK)
	s.auditStore.add("scheduler.delete", r.PathValue("id"), "success", "")
}
