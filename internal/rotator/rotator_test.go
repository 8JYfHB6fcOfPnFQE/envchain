package rotator_test

import (
	"testing"

	"github.com/example/envchain/internal/differ"
	"github.com/example/envchain/internal/rotator"
	"github.com/example/envchain/internal/snapshot"
)

func makeSnap(t *testing.T, name string, vals map[string]string) *snapshot.Snapshot {
	t.Helper()
	s, err := snapshot.New(name, vals)
	if err != nil {
		t.Fatalf("snapshot.New: %v", err)
	}
	return s
}

func makeRotator(t *testing.T) *rotator.Rotator {
	t.Helper()
	d, err := differ.New()
	if err != nil {
		t.Fatalf("differ.New: %v", err)
	}
	r, err := rotator.New(d)
	if err != nil {
		t.Fatalf("rotator.New: %v", err)
	}
	return r
}

func TestNew_NilDifferReturnsError(t *testing.T) {
	_, err := rotator.New(nil)
	if err == nil {
		t.Fatal("expected error for nil differ")
	}
}

func TestApply_NilBaseReturnsError(t *testing.T) {
	r := makeRotator(t)
	next := makeSnap(t, "prod", map[string]string{"A": "1"})
	_, err := r.Apply(nil, next)
	if err == nil {
		t.Fatal("expected error for nil base")
	}
}

func TestApply_NilNextReturnsError(t *testing.T) {
	r := makeRotator(t)
	base := makeSnap(t, "prod", map[string]string{"A": "1"})
	_, err := r.Apply(base, nil)
	if err == nil {
		t.Fatal("expected error for nil next")
	}
}

func TestApply_ContextNameMismatchReturnsError(t *testing.T) {
	r := makeRotator(t)
	base := makeSnap(t, "prod", map[string]string{"A": "1"})
	next := makeSnap(t, "staging", map[string]string{"A": "1"})
	_, err := r.Apply(base, next)
	if err == nil {
		t.Fatal("expected error for mismatched context names")
	}
}

func TestApply_NoChanges(t *testing.T) {
	r := makeRotator(t)
	vals := map[string]string{"DB_PASS": "secret"}
	base := makeSnap(t, "prod", vals)
	next := makeSnap(t, "prod", vals)
	res, err := r.Apply(base, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasChanges() {
		t.Errorf("expected no changes, got rotated=%v added=%v removed=%v", res.Rotated, res.Added, res.Removed)
	}
}

func TestApply_DetectsRotatedKey(t *testing.T) {
	r := makeRotator(t)
	base := makeSnap(t, "prod", map[string]string{"DB_PASS": "old"})
	next := makeSnap(t, "prod", map[string]string{"DB_PASS": "new"})
	res, err := r.Apply(base, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Rotated) != 1 || res.Rotated[0] != "DB_PASS" {
		t.Errorf("expected DB_PASS in Rotated, got %v", res.Rotated)
	}
}

func TestApply_DetectsAddedAndRemoved(t *testing.T) {
	r := makeRotator(t)
	base := makeSnap(t, "prod", map[string]string{"OLD_KEY": "v"})
	next := makeSnap(t, "prod", map[string]string{"NEW_KEY": "v"})
	res, err := r.Apply(base, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Added) != 1 || res.Added[0] != "NEW_KEY" {
		t.Errorf("expected NEW_KEY in Added, got %v", res.Added)
	}
	if len(res.Removed) != 1 || res.Removed[0] != "OLD_KEY" {
		t.Errorf("expected OLD_KEY in Removed, got %v", res.Removed)
	}
	if !res.HasChanges() {
		t.Error("expected HasChanges to return true")
	}
}
