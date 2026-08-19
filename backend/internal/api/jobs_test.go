package api

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestJobEngineRestoresProgressAndMarksInterruptedJobs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "jobs.json")
	t.Setenv("LABSOS_JOBS_FILE", path)
	engine := newJobEngine()
	job := engine.create("app.install", "demo")
	engine.update(job.ID, func(item *managedJob) {
		item.Status = "running"
		item.Progress = 42
		item.Message = "baixando imagem"
	})
	engine.log(job.ID, "imagem recebida")

	restored := newJobEngine()
	loaded, ok := restored.get(job.ID)
	if !ok {
		t.Fatal("job was not restored")
	}
	if loaded.Status != "error" || loaded.Progress != 100 || loaded.Error == "" {
		t.Fatalf("interrupted job was not closed safely: %+v", loaded)
	}
	if len(loaded.Logs) != 1 || !strings.HasSuffix(loaded.Logs[0], "imagem recebida") {
		t.Fatalf("job logs were not restored: %#v", loaded.Logs)
	}
}
