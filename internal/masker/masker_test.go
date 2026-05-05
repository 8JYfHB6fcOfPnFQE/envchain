package masker_test

import (
	"testing"

	"github.com/yourusername/envchain/internal/masker"
)

func TestNew_DefaultMask(t *testing.T) {
	m, err := masker.New([]string{"SECRET"}, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Mask() != "********" {
		t.Errorf("expected default mask, got %q", m.Mask())
	}
}

func TestNew_CustomMask(t *testing.T) {
	m, err := masker.New([]string{"TOKEN"}, "REDACTED")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Mask() != "REDACTED" {
		t.Errorf("expected REDACTED, got %q", m.Mask())
	}
}

func TestNew_BlankKeyReturnsError(t *testing.T) {
	_, err := masker.New([]string{"VALID", "  "}, "")
	if err == nil {
		t.Fatal("expected error for blank key, got nil")
	}
}

func TestApply_MasksSensitiveKeys(t *testing.T) {
	m, _ := masker.New([]string{"PASSWORD", "API_KEY"}, "")
	env := map[string]string{
		"PASSWORD": "s3cr3t",
		"API_KEY":  "abc123",
		"HOST":     "localhost",
	}
	result := m.Apply(env)
	if result["PASSWORD"] != "********" {
		t.Errorf("PASSWORD should be masked, got %q", result["PASSWORD"])
	}
	if result["API_KEY"] != "********" {
		t.Errorf("API_KEY should be masked, got %q", result["API_KEY"])
	}
	if result["HOST"] != "localhost" {
		t.Errorf("HOST should be unchanged, got %q", result["HOST"])
	}
}

func TestApply_CaseInsensitiveKeyMatch(t *testing.T) {
	m, _ := masker.New([]string{"secret"}, "")
	env := map[string]string{"SECRET": "topsecret"}
	result := m.Apply(env)
	if result["SECRET"] != "********" {
		t.Errorf("expected masked value, got %q", result["SECRET"])
	}
}

func TestApply_DoesNotMutateOriginal(t *testing.T) {
	m, _ := masker.New([]string{"TOKEN"}, "")
	original := map[string]string{"TOKEN": "mytoken", "APP": "envchain"}
	_ = m.Apply(original)
	if original["TOKEN"] != "mytoken" {
		t.Error("Apply must not mutate the original map")
	}
}

func TestIsSensitive_ReturnsCorrectly(t *testing.T) {
	m, _ := masker.New([]string{"DB_PASS"}, "")
	if !m.IsSensitive("db_pass") {
		t.Error("expected db_pass to be sensitive")
	}
	if m.IsSensitive("DB_HOST") {
		t.Error("expected DB_HOST to not be sensitive")
	}
}
