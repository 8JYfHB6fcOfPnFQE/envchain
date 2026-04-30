package auditor

import (
	"fmt"
	"time"

	"github.com/envchain/envchain/internal/snapshot"
)

// Entry represents a single audit log entry recording a resolved environment event.
type Entry struct {
	Timestamp   time.Time
	ContextName string
	Keys        []string
	Note        string
}

// Auditor records audit entries for resolved environment snapshots.
type Auditor struct {
	entries []Entry
}

// New creates a new Auditor instance.
func New() *Auditor {
	return &Auditor{
		entries: make([]Entry, 0),
	}
}

// Record appends an audit entry for the given snapshot.
func (a *Auditor) Record(contextName string, snap *snapshot.Snapshot, note string) error {
	if contextName == "" {
		return fmt.Errorf("auditor: context name must not be empty")
	}
	if snap == nil {
		return fmt.Errorf("auditor: snapshot must not be nil")
	}
	entry := Entry{
		Timestamp:   time.Now().UTC(),
		ContextName: contextName,
		Keys:        snap.Keys(),
		Note:        note,
	}
	a.entries = append(a.entries, entry)
	return nil
}

// Entries returns a copy of all recorded audit entries.
func (a *Auditor) Entries() []Entry {
	result := make([]Entry, len(a.entries))
	copy(result, a.entries)
	return result
}

// Count returns the number of recorded audit entries.
func (a *Auditor) Count() int {
	return len(a.entries)
}

// Clear removes all recorded audit entries.
func (a *Auditor) Clear() {
	a.entries = make([]Entry, 0)
}
