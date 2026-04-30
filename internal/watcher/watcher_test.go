package watcher_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/envchain/internal/differ"
	"github.com/user/envchain/internal/snapshot"
	"github.com/user/envchain/internal/watcher"
)

func makeSnap(t *testing.T, name string, vals map[string]string) *snapshot.Snapshot {
	t.Helper()
	s, err := snapshot.New(name, vals)
	if err != nil {
		t.Fatalf("snapshot.New: %v", err)
	}
	return s
}

func TestNew_InvalidInterval(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"K": "v"})
	_, err := watcher.New(0, base, func() (*snapshot.Snapshot, error) { return base, nil }, func(_ *differ.Diff) {})
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNew_NilBaseline(t *testing.T) {
	_, err := watcher.New(time.Second, nil, func() (*snapshot.Snapshot, error) { return nil, nil }, func(_ *differ.Diff) {})
	if err == nil {
		t.Fatal("expected error for nil baseline")
	}
}

func TestNew_NilResolver(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"K": "v"})
	_, err := watcher.New(time.Second, base, nil, func(_ *differ.Diff) {})
	if err == nil {
		t.Fatal("expected error for nil resolver")
	}
}

func TestNew_NilHandler(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"K": "v"})
	_, err := watcher.New(time.Second, base, func() (*snapshot.Snapshot, error) { return base, nil }, nil)
	if err == nil {
		t.Fatal("expected error for nil handler")
	}
}

func TestNew_Valid(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"K": "v"})
	_, err := watcher.New(time.Millisecond*50, base, func() (*snapshot.Snapshot, error) { return base, nil }, func(_ *differ.Diff) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestStart_DetectsDrift(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"KEY": "original"})
	updated := makeSnap(t, "updated", map[string]string{"KEY": "changed"})

	var mu sync.Mutex
	var detected bool

	w, err := watcher.New(
		time.Millisecond*20,
		base,
		func() (*snapshot.Snapshot, error) { return updated, nil },
		func(_ *differ.Diff) {
			mu.Lock()
			detected = true
			mu.Unlock()
		},
	)
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}
	w.Start()
	time.Sleep(time.Millisecond * 80)
	w.Stop()

	mu.Lock()
	defer mu.Unlock()
	if !detected {
		t.Error("expected drift to be detected")
	}
}

func TestStart_NoDriftWhenEqual(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"KEY": "same"})

	var mu sync.Mutex
	var count int

	w, err := watcher.New(
		time.Millisecond*20,
		base,
		func() (*snapshot.Snapshot, error) { return base, nil },
		func(_ *differ.Diff) {
			mu.Lock()
			count++
			mu.Unlock()
		},
	)
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}
	w.Start()
	time.Sleep(time.Millisecond * 80)
	w.Stop()

	mu.Lock()
	defer mu.Unlock()
	if count != 0 {
		t.Errorf("expected no drift calls, got %d", count)
	}
}

func TestStart_ResolveErrorSilent(t *testing.T) {
	base := makeSnap(t, "base", map[string]string{"K": "v"})
	var called bool
	w, err := watcher.New(
		time.Millisecond*20,
		base,
		func() (*snapshot.Snapshot, error) { return nil, errors.New("boom") },
		func(_ *differ.Diff) { called = true },
	)
	if err != nil {
		t.Fatalf("watcher.New: %v", err)
	}
	w.Start()
	time.Sleep(time.Millisecond * 60)
	w.Stop()
	if called {
		t.Error("handler should not be called on resolver error")
	}
}
