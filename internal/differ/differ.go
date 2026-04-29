// Package differ provides functionality for comparing two sets of
// environment variable snapshots and reporting additions, removals,
// and changes between them.
package differ

import (
	"errors"
	"sort"
)

// DiffKind represents the type of change detected for a key.
type DiffKind string

const (
	Added   DiffKind = "added"
	Removed DiffKind = "removed"
	Changed DiffKind = "changed"
)

// Entry represents a single difference between two snapshots.
type Entry struct {
	Key      string
	Kind     DiffKind
	OldValue string
	NewValue string
}

// Result holds all differences between two snapshots.
type Result struct {
	entries []Entry
}

// Entries returns a stable, sorted slice of diff entries.
func (r *Result) Entries() []Entry {
	out := make([]Entry, len(r.entries))
	copy(out, r.entries)
	return out
}

// HasChanges reports whether any differences were found.
func (r *Result) HasChanges() bool {
	return len(r.entries) > 0
}

// Differ compares two maps of environment variables.
type Differ struct{}

// New returns a new Differ instance.
func New() *Differ {
	return &Differ{}
}

// Compare returns a Result describing the differences between base and next.
// Keys present only in base are Removed; keys only in next are Added;
// keys in both with differing values are Changed.
func (d *Differ) Compare(base, next map[string]string) (*Result, error) {
	if base == nil {
		return nil, errors.New("differ: base map must not be nil")
	}
	if next == nil {
		return nil, errors.New("differ: next map must not be nil")
	}

	var entries []Entry

	for k, oldVal := range base {
		if newVal, ok := next[k]; !ok {
			entries = append(entries, Entry{Key: k, Kind: Removed, OldValue: oldVal})
		} else if newVal != oldVal {
			entries = append(entries, Entry{Key: k, Kind: Changed, OldValue: oldVal, NewValue: newVal})
		}
	}

	for k, newVal := range next {
		if _, ok := base[k]; !ok {
			entries = append(entries, Entry{Key: k, Kind: Added, NewValue: newVal})
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Key < entries[j].Key
	})

	return &Result{entries: entries}, nil
}
