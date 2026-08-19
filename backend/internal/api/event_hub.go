package api

import (
	"sync"
	"time"

	"labsos/backend/internal/platform"
)

type eventHub struct {
	mu          sync.Mutex
	nextID      uint64
	subscribers map[chan platform.Event]struct{}
}

func newEventHub() *eventHub {
	return &eventHub{subscribers: make(map[chan platform.Event]struct{})}
}

func (h *eventHub) subscribe() (<-chan platform.Event, func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	channel := make(chan platform.Event, 16)
	h.subscribers[channel] = struct{}{}
	return channel, func() {
		h.mu.Lock()
		if _, ok := h.subscribers[channel]; ok {
			delete(h.subscribers, channel)
			close(channel)
		}
		h.mu.Unlock()
	}
}

func (h *eventHub) publish(kind, message string) {
	h.mu.Lock()
	h.nextID++
	event := platform.Event{ID: "hub-" + time.Now().UTC().Format("20060102150405.000000000"), Type: kind, Message: message}
	for channel := range h.subscribers {
		select {
		case channel <- event:
		default:
		}
	}
	h.mu.Unlock()
}
