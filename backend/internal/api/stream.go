package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"labsos/backend/internal/platform"
)

// eventsStream provides a lightweight server-sent event channel for dashboard
// consumers. It is intentionally read-only and stays local with the API.
func (s *server) eventsStream(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "STREAM_UNAVAILABLE", "SSE is unavailable")
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	subscriber, unsubscribe := s.eventHub.subscribe()
	defer unsubscribe()
	if events, err := s.provider.Events(r.Context()); err == nil {
		data, _ := json.Marshal(events)
		_, _ = fmt.Fprintf(w, "event: snapshot\ndata: %s\n\n", data)
		flusher.Flush()
	}
	heartbeat := time.NewTicker(15 * time.Second)
	defer heartbeat.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case event, ok := <-subscriber:
			if !ok {
				return
			}
			data, _ := json.Marshal([]platform.Event{event})
			_, _ = fmt.Fprintf(w, "event: events\ndata: %s\n\n", data)
			flusher.Flush()
		case <-heartbeat.C:
			_, _ = fmt.Fprint(w, ": keep-alive\n\n")
			flusher.Flush()
		}
	}
}
