package snapshot_test

import (
	"testing"

	"github.com/your-org/envchain/internal/snapshot"
)

func TestNew_ValidSnapshot(t *testing.T) {
	vals := map[string]string{"FOO": "bar", "BAZ": "qux"}
	s, err := snapshot.New("prod", vals)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Name != "prod" {
		t.Errorf("expected name 'prod', got %q", s.Name)
	}
	if s.Values["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", s.Values["FOO"])
	}
}

func TestNew_EmptyNameReturnsError(t *testing.T) {
	_, err := snapshot.New("", map[string]string{})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestNew_NilValuesReturnsError(t *testing.T) {
	_, err := snapshot.New("test", nil)
	if err == nil {
		t.Fatal("expected error for nil values, got nil")
	}
}

func TestNew_MutationIsolation(t *testing.T) {
	vals := map[string]string{"KEY": "original"}
	s, _ := snapshot.New("test", vals)
	vals["KEY"] = "mutated"
	if s.Values["KEY"] != "original" {
		t.Error("snapshot values should be isolated from source map mutations")
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	s, _ := snapshot.New("test", map[string]string{"Z": "1", "A": "2", "M": "3"})
	keys := s.Keys()
	expected := []string{"A", "M", "Z"}
	for i, k := range keys {
		if k != expected[i] {
			t.Errorf("expected key[%d]=%q, got %q", i, expected[i], k)
		}
	}
}

func TestCompare_DetectsAdded(t *testing.T) {
	old, _ := snapshot.New("old", map[string]string{"A": "1"})
	new, _ := snapshot.New("new", map[string]string{"A": "1", "B": "2"})
	diff := old.Compare(new)
	if _, ok := diff.Added["B"]; !ok {
		t.Error("expected B to be in Added")
	}
}

func TestCompare_DetectsRemoved(t *testing.T) {
	old, _ := snapshot.New("old", map[string]string{"A": "1", "B": "2"})
	new, _ := snapshot.New("new", map[string]string{"A": "1"})
	diff := old.Compare(new)
	if _, ok := diff.Removed["B"]; !ok {
		t.Error("expected B to be in Removed")
	}
}

func TestCompare_DetectsChanged(t *testing.T) {
	old, _ := snapshot.New("old", map[string]string{"A": "old_val"})
	new, _ := snapshot.New("new", map[string]string{"A": "new_val"})
	diff := old.Compare(new)
	if pair, ok := diff.Changed["A"]; !ok || pair[0] != "old_val" || pair[1] != "new_val" {
		t.Errorf("expected Changed[A]=[old_val, new_val], got %v", diff.Changed["A"])
	}
}

func TestDiff_HasChanges_False(t *testing.T) {
	old, _ := snapshot.New("old", map[string]string{"A": "1"})
	new, _ := snapshot.New("new", map[string]string{"A": "1"})
	diff := old.Compare(new)
	if diff.HasChanges() {
		t.Error("expected no changes between identical snapshots")
	}
}
