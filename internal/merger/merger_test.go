package merger

import (
	"testing"
)

func TestNew_EmptyMerger(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("expected non-nil Merger")
	}
	if len(m.Sources()) != 0 {
		t.Errorf("expected 0 sources, got %d", len(m.Sources()))
	}
}

func TestAdd_ValidSource(t *testing.T) {
	m := New()
	err := m.Add("base", map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srcs := m.Sources(); len(srcs) != 1 || srcs[0] != "base" {
		t.Errorf("expected sources [base], got %v", srcs)
	}
}

func TestAdd_EmptyNameReturnsError(t *testing.T) {
	m := New()
	err := m.Add("", map[string]string{"FOO": "bar"})
	if err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestAdd_NilMapReturnsError(t *testing.T) {
	m := New()
	err := m.Add("base", nil)
	if err == nil {
		t.Fatal("expected error for nil map, got nil")
	}
}

func TestMerge_NoSourcesReturnsError(t *testing.T) {
	m := New()
	_, err := m.Merge()
	if err == nil {
		t.Fatal("expected error when merging with no sources")
	}
}

func TestMerge_SingleSource(t *testing.T) {
	m := New()
	_ = m.Add("base", map[string]string{"FOO": "bar", "BAZ": "qux"})
	result, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "bar" || result["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", result)
	}
}

func TestMerge_LaterSourceOverrides(t *testing.T) {
	m := New()
	_ = m.Add("base", map[string]string{"FOO": "original", "KEEP": "yes"})
	_ = m.Add("override", map[string]string{"FOO": "overridden", "EXTRA": "added"})
	result, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["FOO"] != "overridden" {
		t.Errorf("expected FOO=overridden, got %q", result["FOO"])
	}
	if result["KEEP"] != "yes" {
		t.Errorf("expected KEEP=yes, got %q", result["KEEP"])
	}
	if result["EXTRA"] != "added" {
		t.Errorf("expected EXTRA=added, got %q", result["EXTRA"])
	}
}
