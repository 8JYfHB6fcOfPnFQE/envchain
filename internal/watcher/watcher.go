// Package watcher monitors environment variable snapshots for drift
// between a recorded baseline and the current resolved state.
package watcher

import (
	"errors"
	"fmt"
	"time"

	"github.com/user/envchain/internal/differ"
	"github.com/user/envchain/internal/snapshot"
)

// DriftHandler is called when drift is detected between snapshots.
type DriftHandler func(diff *differ.Diff)

// Watcher polls a snapshot source and reports drift via a handler.
type Watcher struct {
	interval time.Duration
	baseline *snapshot.Snapshot
	resolve  func() (*snapshot.Snapshot, error)
	handler  DriftHandler
	diff     *differ.Differ
	stop     chan struct{}
}

// New creates a new Watcher with the given poll interval, resolver function,
// baseline snapshot, and drift handler.
func New(interval time.Duration, baseline *snapshot.Snapshot, resolve func() (*snapshot.Snapshot, error), handler DriftHandler) (*Watcher, error) {
	if interval <= 0 {
		return nil, errors.New("watcher: interval must be positive")
	}
	if baseline == nil {
		return nil, errors.New("watcher: baseline snapshot must not be nil")
	}
	if resolve == nil {
		return nil, errors.New("watcher: resolve function must not be nil")
	}
	if handler == nil {
		return nil, errors.New("watcher: drift handler must not be nil")
	}
	return &Watcher{
		interval: interval,
		baseline: baseline,
		resolve:  resolve,
		handler:  handler,
		diff:     differ.New(),
		stop:     make(chan struct{}),
	}, nil
}

// Start begins polling in a background goroutine.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				w.check()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts the watcher's polling loop.
func (w *Watcher) Stop() {
	close(w.stop)
}

// check performs a single drift evaluation.
func (w *Watcher) check() {
	current, err := w.resolve()
	if err != nil {
		_ = fmt.Errorf("watcher: resolve error: %w", err)
		return
	}
	diff, err := w.diff.Compare(w.baseline, current)
	if err != nil {
		return
	}
	if diff.HasChanges() {
		w.handler(diff)
	}
}
