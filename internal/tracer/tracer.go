package tracer

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// Event represents a single trace event recorded during environment resolution.
type Event struct {
	Timestamp time.Time
	Context   string
	Key       string
	Action    string
	Detail    string
}

// Tracer records trace events for environment variable resolution steps.
type Tracer struct {
	mu     sync.Mutex
	events []Event
}

// New creates a new Tracer instance.
func New() *Tracer {
	return &Tracer{
		events: make([]Event, 0),
	}
}

// Record appends a new trace event. Returns an error if context or action is blank.
func (t *Tracer) Record(context, key, action, detail string) error {
	if context == "" {
		return errors.New("tracer: context must not be empty")
	}
	if action == "" {
		return errors.New("tracer: action must not be empty")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.events = append(t.events, Event{
		Timestamp: time.Now(),
		Context:   context,
		Key:       key,
		Action:    action,
		Detail:    detail,
	})
	return nil
}

// Events returns a copy of all recorded trace events.
func (t *Tracer) Events() []Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	copy := make([]Event, len(t.events))
	for i, e := range t.events {
		copy[i] = e
	}
	return copy
}

// FilterByContext returns all events matching the given context name.
func (t *Tracer) FilterByContext(context string) []Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	var result []Event
	for _, e := range t.events {
		if e.Context == context {
			result = append(result, e)
		}
	}
	return result
}

// Summary returns a human-readable summary of all recorded events.
func (t *Tracer) Summary() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.events) == 0 {
		return "no trace events recorded"
	}
	summary := fmt.Sprintf("%d trace event(s):\n", len(t.events))
	for _, e := range t.events {
		summary += fmt.Sprintf("  [%s] ctx=%s key=%s action=%s detail=%s\n",
			e.Timestamp.Format(time.RFC3339), e.Context, e.Key, e.Action, e.Detail)
	}
	return summary
}
