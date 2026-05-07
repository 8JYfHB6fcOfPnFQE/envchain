package scoper

import (
	"testing"
)

func TestNew_BlankPrefixReturnsError(t *testing.T) {
	_, err := New("", false)
	if err == nil {
		t.Fatal("expected error for blank prefix")
	}
}

func TestNew_ValidPrefix(t *testing.T) {
	s, err := New("APP_", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil scoper")
	}
}

func TestScope_NilEnvReturnsError(t *testing.T) {
	s, _ := New("APP_", false)
	_, err := s.Scope(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestScope_FiltersToPrefix(t *testing.T) {
	s, _ := New("APP_", false)
	env := map[string]string{
		"APP_HOST": "localhost",
		"APP_PORT": "8080",
		"DB_HOST":  "db",
	}
	out, err := s.Scope(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if out["APP_HOST"] != "localhost" {
		t.Errorf("expected APP_HOST=localhost, got %q", out["APP_HOST"])
	}
}

func TestScope_StripPrefix(t *testing.T) {
	s, _ := New("APP_", true)
	env := map[string]string{"APP_HOST": "localhost", "APP_PORT": "9000"}
	out, err := s.Scope(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["HOST"]; !ok {
		t.Errorf("expected stripped key HOST")
	}
	if _, ok := out["APP_HOST"]; ok {
		t.Errorf("expected APP_HOST to be stripped")
	}
}

func TestScope_ExplicitInclude(t *testing.T) {
	s, _ := New("APP_", false)
	if err := s.Include("GLOBAL_ID"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{"APP_HOST": "h", "GLOBAL_ID": "42", "OTHER": "x"}
	out, err := s.Scope(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["GLOBAL_ID"] != "42" {
		t.Errorf("expected GLOBAL_ID in output")
	}
	if _, ok := out["OTHER"]; ok {
		t.Errorf("OTHER should not be in output")
	}
}

func TestInclude_BlankKeyReturnsError(t *testing.T) {
	s, _ := New("APP_", false)
	if err := s.Include(""); err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestScope_CaseInsensitivePrefix(t *testing.T) {
	s, _ := New("app_", false)
	env := map[string]string{"APP_HOST": "h", "other": "v"}
	out, _ := s.Scope(env)
	if _, ok := out["APP_HOST"]; !ok {
		t.Errorf("expected APP_HOST matched case-insensitively")
	}
}
