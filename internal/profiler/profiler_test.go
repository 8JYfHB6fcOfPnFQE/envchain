package profiler

import (
	"testing"
)

func TestDefine_ValidProfile(t *testing.T) {
	p := New()
	err := p.Define("production", []string{"DB_URL", "API_KEY"}, map[string]string{"env": "prod"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefine_EmptyNameReturnsError(t *testing.T) {
	p := New()
	err := p.Define("", []string{"DB_URL"}, nil)
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestDefine_EmptyKeysReturnsError(t *testing.T) {
	p := New()
	err := p.Define("staging", []string{}, nil)
	if err == nil {
		t.Fatal("expected error for empty keys")
	}
}

func TestDefine_DuplicateReturnsError(t *testing.T) {
	p := New()
	_ = p.Define("production", []string{"DB_URL"}, nil)
	err := p.Define("production", []string{"API_KEY"}, nil)
	if err == nil {
		t.Fatal("expected error for duplicate profile")
	}
}

func TestGet_ReturnsProfile(t *testing.T) {
	p := New()
	_ = p.Define("staging", []string{"DB_URL", "SECRET"}, map[string]string{"tier": "low"})
	pr, err := p.Get("staging")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pr.Name() != "staging" {
		t.Errorf("expected name 'staging', got %q", pr.Name())
	}
	keys := pr.Keys()
	if len(keys) != 2 {
		t.Errorf("expected 2 keys, got %d", len(keys))
	}
	tags := pr.Tags()
	if tags["tier"] != "low" {
		t.Errorf("expected tag tier=low, got %q", tags["tier"])
	}
}

func TestGet_NotFoundReturnsError(t *testing.T) {
	p := New()
	_, err := p.Get("unknown")
	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}

func TestList_ReturnsSortedNames(t *testing.T) {
	p := New()
	_ = p.Define("z-env", []string{"A"}, nil)
	_ = p.Define("a-env", []string{"B"}, nil)
	_ = p.Define("m-env", []string{"C"}, nil)
	names := p.List()
	expected := []string{"a-env", "m-env", "z-env"}
	for i, n := range names {
		if n != expected[i] {
			t.Errorf("position %d: expected %q got %q", i, expected[i], n)
		}
	}
}

func TestKeys_MutationIsolation(t *testing.T) {
	p := New()
	_ = p.Define("prod", []string{"DB_URL"}, nil)
	pr, _ := p.Get("prod")
	keys := pr.Keys()
	keys[0] = "MUTATED"
	original := pr.Keys()
	if original[0] == "MUTATED" {
		t.Error("Keys() should return an isolated copy")
	}
}
