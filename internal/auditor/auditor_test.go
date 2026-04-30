package auditor_test

import (
	"testing"

	"github.com/envchain/envchain/internal/auditor"
	"github.com/envchain/envchain/internal/snapshot"
)

func makeSnapshot(t *testing.T, name string, vals map[string]string) *snapshot.Snapshot {
	t.Helper()
	snap, err := snapshot.New(name, vals)
	if err != nil {
		t.Fatalf("failed to create snapshot: %v", err)
	}
	return snap
}

func TestRecord_ValidEntry(t *testing.T) {
	a := auditor.New()
	snap := makeSnapshot(t, "prod", map[string]string{"KEY": "val"})
	if err := a.Record("prod", snap, "initial load"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Count() != 1 {
		t.Errorf("expected 1 entry, got %d", a.Count())
	}
}

func TestRecord_EmptyContextReturnsError(t *testing.T) {
	a := auditor.New()
	snap := makeSnapshot(t, "prod", map[string]string{"KEY": "val"})
	if err := a.Record("", snap, ""); err == nil {
		t.Error("expected error for empty context name")
	}
}

func TestRecord_NilSnapshotReturnsError(t *testing.T) {
	a := auditor.New()
	if err := a.Record("prod", nil, ""); err == nil {
		t.Error("expected error for nil snapshot")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	a := auditor.New()
	snap := makeSnapshot(t, "staging", map[string]string{"A": "1", "B": "2"})
	_ = a.Record("staging", snap, "test")
	entries := a.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	entries[0].Note = "mutated"
	if a.Entries()[0].Note == "mutated" {
		t.Error("entries slice should be a copy")
	}
}

func TestEntries_ContainsKeys(t *testing.T) {
	a := auditor.New()
	snap := makeSnapshot(t, "dev", map[string]string{"FOO": "bar", "BAZ": "qux"})
	_ = a.Record("dev", snap, "")
	entries := a.Entries()
	if len(entries[0].Keys) != 2 {
		t.Errorf("expected 2 keys in entry, got %d", len(entries[0].Keys))
	}
}

func TestClear_RemovesAllEntries(t *testing.T) {
	a := auditor.New()
	snap := makeSnapshot(t, "prod", map[string]string{"X": "y"})
	_ = a.Record("prod", snap, "")
	a.Clear()
	if a.Count() != 0 {
		t.Errorf("expected 0 entries after clear, got %d", a.Count())
	}
}
