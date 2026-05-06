// Package rotator provides utilities for detecting and applying environment
// variable rotations across named contexts, producing a before/after diff.
package rotator

import (
	"errors"
	"fmt"

	"github.com/example/envchain/internal/differ"
	"github.com/example/envchain/internal/snapshot"
)

// Rotator applies a new set of values to an existing snapshot and reports
// which keys were rotated (changed), added, or removed.
type Rotator struct {
	differ *differ.Differ
}

// Result holds the outcome of a rotation operation.
type Result struct {
	// Context is the name of the environment context that was rotated.
	Context string
	// Rotated contains keys whose values changed.
	Rotated []string
	// Added contains keys that are new in the next snapshot.
	Added []string
	// Removed contains keys that were present in base but absent in next.
	Removed []string
}

// New creates a new Rotator backed by the provided Differ.
func New(d *differ.Differ) (*Rotator, error) {
	if d == nil {
		return nil, errors.New("rotator: differ must not be nil")
	}
	return &Rotator{differ: d}, nil
}

// Apply compares base against next and returns a Result describing the
// rotation. Both snapshots must belong to the same named context.
func (r *Rotator) Apply(base, next *snapshot.Snapshot) (*Result, error) {
	if base == nil {
		return nil, errors.New("rotator: base snapshot must not be nil")
	}
	if next == nil {
		return nil, errors.New("rotator: next snapshot must not be nil")
	}
	if base.Name() != next.Name() {
		return nil, fmt.Errorf("rotator: context name mismatch: %q vs %q", base.Name(), next.Name())
	}

	changes, err := r.differ.Compare(base, next)
	if err != nil {
		return nil, fmt.Errorf("rotator: diff failed: %w", err)
	}

	result := &Result{Context: base.Name()}
	for _, c := range changes {
		switch c.Kind {
		case differ.Added:
			result.Added = append(result.Added, c.Key)
		case differ.Removed:
			result.Removed = append(result.Removed, c.Key)
		case differ.Changed:
			result.Rotated = append(result.Rotated, c.Key)
		}
	}
	return result, nil
}

// HasChanges returns true when the Result contains at least one rotation,
// addition, or removal.
func (res *Result) HasChanges() bool {
	return len(res.Rotated)+len(res.Added)+len(res.Removed) > 0
}
