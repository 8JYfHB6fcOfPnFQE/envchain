package interpolator_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/interpolator"
)

func TestNew_ValidInterpolator(t *testing.T) {
	i, err := interpolator.New(false, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if i == nil {
		t.Fatal("expected non-nil interpolator")
	}
}

func TestNew_ZeroMaxDepthReturnsError(t *testing.T) {
	_, err := interpolator.New(false, 0)
	if err == nil {
		t.Fatal("expected error for zero maxDepth")
	}
}

func TestExpand_NilEnvReturnsError(t *testing.T) {
	i, _ := interpolator.New(false, 5)
	_, err := i.Expand(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestExpand_SimpleSubstitution(t *testing.T) {
	i, _ := interpolator.New(false, 5)
	env := map[string]string{
		"BASE": "/usr/local",
		"BIN":  "${BASE}/bin",
	}
	out, err := i.Expand(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["BIN"]; got != "/usr/local/bin" {
		t.Errorf("BIN: got %q, want %q", got, "/usr/local/bin")
	}
}

func TestExpand_UnresolvedLeftAsLiteral(t *testing.T) {
	i, _ := interpolator.New(false, 5)
	env := map[string]string{
		"VAL": "${MISSING}_suffix",
	}
	out, err := i.Expand(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := out["VAL"]; got != "$MISSING_suffix" {
		t.Errorf("VAL: got %q", got)
	}
}

func TestExpandValue_NilEnvReturnsError(t *testing.T) {
	i, _ := interpolator.New(false, 5)
	_, err := i.ExpandValue("${FOO}", nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestExpandValue_DirectSubstitution(t *testing.T) {
	i, _ := interpolator.New(false, 5)
	env := map[string]string{"HOST": "localhost"}
	got, err := i.ExpandValue("http://${HOST}:8080", env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "http://localhost:8080" {
		t.Errorf("got %q, want %q", got, "http://localhost:8080")
	}
}

func TestExpand_NoReferencesUnchanged(t *testing.T) {
	i, _ := interpolator.New(false, 5)
	env := map[string]string{"PLAIN": "hello world"}
	out, err := i.Expand(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PLAIN"] != "hello world" {
		t.Errorf("expected unchanged value, got %q", out["PLAIN"])
	}
}
