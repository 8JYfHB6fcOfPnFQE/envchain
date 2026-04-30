package redactor_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/redactor"
)

func TestNew_DefaultMask(t *testing.T) {
	r, err := redactor.New([]string{"SECRET"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Mask() != "********" {
		t.Errorf("expected default mask, got %q", r.Mask())
	}
}

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := redactor.New([]string{"VALID", ""}, "***")
	if err == nil {
		t.Fatal("expected error for empty key name, got nil")
	}
}

func TestRedact_MasksSensitiveKeys(t *testing.T) {
	r, _ := redactor.New([]string{"PASSWORD", "API_KEY"}, "[REDACTED]")
	env := map[string]string{
		"PASSWORD": "supersecret",
		"API_KEY":  "abc123",
		"HOST":     "localhost",
	}
	out := r.Redact(env)
	if out["PASSWORD"] != "[REDACTED]" {
		t.Errorf("PASSWORD not redacted, got %q", out["PASSWORD"])
	}
	if out["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY not redacted, got %q", out["API_KEY"])
	}
	if out["HOST"] != "localhost" {
		t.Errorf("HOST should not be redacted, got %q", out["HOST"])
	}
}

func TestRedact_CaseInsensitiveKeyMatch(t *testing.T) {
	r, _ := redactor.New([]string{"secret_token"}, "***")
	env := map[string]string{"SECRET_TOKEN": "value"}
	out := r.Redact(env)
	if out["SECRET_TOKEN"] != "***" {
		t.Errorf("expected redacted value, got %q", out["SECRET_TOKEN"])
	}
}

func TestRedact_DoesNotMutateOriginal(t *testing.T) {
	r, _ := redactor.New([]string{"TOKEN"}, "***")
	env := map[string]string{"TOKEN": "original"}
	_ = r.Redact(env)
	if env["TOKEN"] != "original" {
		t.Error("original map was mutated")
	}
}

func TestIsSensitive_ReturnsCorrectly(t *testing.T) {
	r, _ := redactor.New([]string{"DB_PASS"}, "***")
	if !r.IsSensitive("db_pass") {
		t.Error("expected db_pass to be sensitive")
	}
	if r.IsSensitive("DB_HOST") {
		t.Error("expected DB_HOST to not be sensitive")
	}
}

func TestRedact_EmptySensitiveKeys(t *testing.T) {
	r, _ := redactor.New([]string{}, "***")
	env := map[string]string{"FOO": "bar"}
	out := r.Redact(env)
	if out["FOO"] != "bar" {
		t.Errorf("expected FOO unchanged, got %q", out["FOO"])
	}
}
