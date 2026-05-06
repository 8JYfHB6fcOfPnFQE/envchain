package inspector_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/inspector"
)

func TestNew_ValidRequired(t *testing.T) {
	_, err := inspector.New([]string{"HOST", "PORT"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyKeyReturnsError(t *testing.T) {
	_, err := inspector.New([]string{"HOST", ""})
	if err == nil {
		t.Fatal("expected error for empty required key, got nil")
	}
}

func TestInspect_NilEnvReturnsError(t *testing.T) {
	insp, _ := inspector.New([]string{"HOST"})
	_, err := insp.Inspect(nil)
	if err == nil {
		t.Fatal("expected error for nil env, got nil")
	}
}

func TestInspect_CountsKeys(t *testing.T) {
	insp, _ := inspector.New([]string{"HOST", "PORT"})
	env := map[string]string{"HOST": "localhost", "PORT": "8080", "DEBUG": "true"}
	s, err := insp.Inspect(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.TotalKeys != 3 {
		t.Errorf("expected TotalKeys=3, got %d", s.TotalKeys)
	}
	if s.PresentKeys != 2 {
		t.Errorf("expected PresentKeys=2, got %d", s.PresentKeys)
	}
	if s.MissingKeys != 0 {
		t.Errorf("expected MissingKeys=0, got %d", s.MissingKeys)
	}
}

func TestInspect_DetectsMissingRequired(t *testing.T) {
	insp, _ := inspector.New([]string{"HOST", "PORT", "SECRET"})
	env := map[string]string{"HOST": "localhost"}
	s, err := insp.Inspect(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.MissingKeys != 2 {
		t.Errorf("expected MissingKeys=2, got %d", s.MissingKeys)
	}
}

func TestInspect_KeyNamesSorted(t *testing.T) {
	insp, _ := inspector.New(nil)
	env := map[string]string{"ZEBRA": "1", "ALPHA": "2", "MANGO": "3"}
	s, _ := insp.Inspect(env)
	if len(s.KeyNames) != 3 || s.KeyNames[0] != "ALPHA" || s.KeyNames[2] != "ZEBRA" {
		t.Errorf("expected sorted key names, got %v", s.KeyNames)
	}
}

func TestMissingRequired_ReturnsAbsentKeys(t *testing.T) {
	insp, _ := inspector.New([]string{"HOST", "PORT", "TOKEN"})
	env := map[string]string{"HOST": "localhost", "PORT": ""}
	missing, err := insp.MissingRequired(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(missing) != 2 {
		t.Errorf("expected 2 missing keys, got %d: %v", len(missing), missing)
	}
}

func TestMissingRequired_NilEnvReturnsError(t *testing.T) {
	insp, _ := inspector.New([]string{"HOST"})
	_, err := insp.MissingRequired(nil)
	if err == nil {
		t.Fatal("expected error for nil env, got nil")
	}
}
