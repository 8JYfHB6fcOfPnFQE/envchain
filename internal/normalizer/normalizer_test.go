package normalizer_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/normalizer"
)

func TestNew_ValidRules(t *testing.T) {
	rules := normalizer.DefaultRules()
	_, err := normalizer.New(rules)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestNew_BlankRuleNameReturnsError(t *testing.T) {
	rules := []normalizer.Rule{
		{Name: "", Apply: func(k, v string) (string, string) { return k, v }},
	}
	_, err := normalizer.New(rules)
	if err == nil {
		t.Fatal("expected error for blank rule name")
	}
}

func TestNew_NilApplyReturnsError(t *testing.T) {
	rules := []normalizer.Rule{
		{Name: "bad", Apply: nil},
	}
	_, err := normalizer.New(rules)
	if err == nil {
		t.Fatal("expected error for nil Apply")
	}
}

func TestNormalize_NilEnvReturnsError(t *testing.T) {
	n, _ := normalizer.New(normalizer.DefaultRules())
	_, err := n.Normalize(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestNormalize_UppercasesKeys(t *testing.T) {
	n, _ := normalizer.New(normalizer.DefaultRules())
	env := map[string]string{"db_host": "localhost", "api_key": "secret"}
	out, err := n.Normalize(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY to be present")
	}
}

func TestNormalize_TrimsWhitespace(t *testing.T) {
	n, _ := normalizer.New(normalizer.DefaultRules())
	env := map[string]string{" PORT ": "  8080  "}
	out, err := n.Normalize(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v, ok := out["PORT"]; !ok || v != "8080" {
		t.Errorf("expected PORT=8080, got %q", v)
	}
}

func TestNormalize_DoesNotMutateOriginal(t *testing.T) {
	n, _ := normalizer.New(normalizer.DefaultRules())
	env := map[string]string{"host": "  db  "}
	_, err := n.Normalize(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["host"] != "  db  " {
		t.Error("original map was mutated")
	}
}

func TestNormalize_CustomRuleApplied(t *testing.T) {
	rules := []normalizer.Rule{
		{
			Name: "PrefixValues",
			Apply: func(k, v string) (string, string) {
				return k, "prefix_" + v
			},
		},
	}
	n, _ := normalizer.New(rules)
	out, err := n.Normalize(map[string]string{"KEY": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "prefix_val" {
		t.Errorf("expected prefix_val, got %q", out["KEY"])
	}
}
