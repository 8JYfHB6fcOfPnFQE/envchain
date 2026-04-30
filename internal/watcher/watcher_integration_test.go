package watcher_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/envchain/internal/differ"
	"github.com/user/envchain/internal/snapshot"
	"github.com/user/envchain/internal/watcher"
)

// TestIntegration_WatcherDetectsMultipleDrifts verifies that the watcher
// fires the handler on every poll cycle where drift exists.
func TestIntegration_WatcherDetectsMultipleDrifts(t *testing.T) {
	base := makeSnap(t, "prod-base", map[string]string{
		"DB_HOST": "db.internal",
		"API_KEY": "secret",
	})

	drifted := makeSnap(t, "prod-current", map[string]string{
		"DB_HOST": "db.internal",
		"API_KEY": "rotated",
		"NEW_VAR": "added",
	})

	var mu sync.Mutex
	var diffs []*differ.Diff

	w, err := watcher.New(
		time.Millisecond*15,
		base,
		func() (*snapshot.Snapshot, error) { return drifted, nil },
		func(d *differ.Diff) {
			mu.Lock()
			diffs = append(diffs, d)
			mu.Unlock()
		},
	)
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}

	w.Start()
	time.Sleep(time.Millisecond * 70)
	w.Stop()

	mu.Lock()
	defer mu.Unlock()

	if len(diffs) == 0 {
		t.Fatal("expected at least one drift event")
	}

	// Each diff should report changes (modified API_KEY + added NEW_VAR).
	for i, d := range diffs {
		if !d.HasChanges() {
			t.Errorf("diff[%d]: expected HasChanges to be true", i)
		}
	}
}
