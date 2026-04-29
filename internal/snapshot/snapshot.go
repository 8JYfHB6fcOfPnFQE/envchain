// Package snapshot provides functionality for capturing and comparing
// environment variable sets at a point in time, enabling drift detection
// between expected and actual deployment configurations.
package snapshot

import (
	"fmt"
	"sort"
	"time"
)

// Snapshot holds a captured set of environment variables with metadata.
type Snapshot struct {
	Name      string
	CapturedAt time.Time
	Values    map[string]string
}

// Diff represents the difference between two snapshots.
type Diff struct {
	Added   map[string]string
	Removed map[string]string
	Changed map[string][2]string // key -> [old, new]
}

// New creates a new Snapshot with the given name and values.
func New(name string, values map[string]string) (*Snapshot, error) {
	if name == "" {
		return nil, fmt.Errorf("snapshot name must not be empty")
	}
	if values == nil {
		return nil, fmt.Errorf("snapshot values must not be nil")
	}
	copy := make(map[string]string, len(values))
	for k, v := range values {
		copy[k] = v
	}
	return &Snapshot{
		Name:       name,
		CapturedAt: time.Now(),
		Values:     copy,
	}, nil
}

// Keys returns a sorted list of all keys in the snapshot.
func (s *Snapshot) Keys() []string {
	keys := make([]string, 0, len(s.Values))
	for k := range s.Values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Compare computes the diff between this snapshot (old) and another (new).
func (s *Snapshot) Compare(other *Snapshot) *Diff {
	diff := &Diff{
		Added:   make(map[string]string),
		Removed: make(map[string]string),
		Changed: make(map[string][2]string),
	}
	for k, v := range other.Values {
		oldVal, exists := s.Values[k]
		if !exists {
			diff.Added[k] = v
		} else if oldVal != v {
			diff.Changed[k] = [2]string{oldVal, v}
		}
	}
	for k, v := range s.Values {
		if _, exists := other.Values[k]; !exists {
			diff.Removed[k] = v
		}
	}
	return diff
}

// HasChanges returns true if the diff contains any additions, removals, or changes.
func (d *Diff) HasChanges() bool {
	return len(d.Added) > 0 || len(d.Removed) > 0 || len(d.Changed) > 0
}
