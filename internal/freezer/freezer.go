package freezer

import (
	"errors"
	"sync"
)

// Freezer captures a point-in-time immutable copy of an environment map
// and prevents further mutation through the original reference.
type Freezer struct {
	mu     sync.RWMutex
	frames map[string]map[string]string
}

// New returns a new Freezer instance.
func New() *Freezer {
	return &Freezer{
		frames: make(map[string]map[string]string),
	}
}

// Freeze stores a deep copy of env under the given label.
// Returns an error if the label is empty or env is nil.
func (f *Freezer) Freeze(label string, env map[string]string) error {
	if label == "" {
		return errors.New("freezer: label must not be empty")
	}
	if env == nil {
		return errors.New("freezer: env must not be nil")
	}

	copy := make(map[string]string, len(env))
	for k, v := range env {
		copy[k] = v
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	f.frames[label] = copy
	return nil
}

// Thaw returns the frozen environment for the given label.
// The returned map is a defensive copy. Returns an error if the label
// was never frozen.
func (f *Freezer) Thaw(label string) (map[string]string, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()

	frame, ok := f.frames[label]
	if !ok {
		return nil, errors.New("freezer: no frozen frame for label " + label)
	}

	out := make(map[string]string, len(frame))
	for k, v := range frame {
		out[k] = v
	}
	return out, nil
}

// Labels returns all currently frozen labels in insertion-independent order.
func (f *Freezer) Labels() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()

	labels := make([]string, 0, len(f.frames))
	for l := range f.frames {
		labels = append(labels, l)
	}
	return labels
}

// Drop removes a frozen frame by label. Returns an error if not found.
func (f *Freezer) Drop(label string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if _, ok := f.frames[label]; !ok {
		return errors.New("freezer: no frozen frame for label " + label)
	}
	delete(f.frames, label)
	return nil
}
