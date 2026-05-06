package linter

import (
	"testing"
)

func TestNew_ValidRules(t *testing.T) {
	rules := DefaultRules()
	l, err := New(rules)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l == nil {
		t.Fatal("expected non-nil Linter")
	}
}

func TestNew_BlankRuleNameReturnsError(t *testing.T) {
	_, err := New([]Rule{{Name: "", Check: func(_, _ string) bool { return false }}})
	if err == nil {
		t.Fatal("expected error for blank rule name")
	}
}

func TestNew_NilCheckReturnsError(t *testing.T) {
	_, err := New([]Rule{{Name: "bad-rule", Check: nil}})
	if err == nil {
		t.Fatal("expected error for nil Check")
	}
}

func TestLint_NilEnvReturnsError(t *testing.T) {
	l, _ := New(DefaultRules())
	_, err := l.Lint(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestLint_EmptyEnvReturnsNoIssues(t *testing.T) {
	l, _ := New(DefaultRules())
	issues, err := l.Lint(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected 0 issues, got %d", len(issues))
	}
}

func TestLint_DetectsEmptyValue(t *testing.T) {
	l, _ := New(DefaultRules())
	issues, err := l.Lint(map[string]string{"MY_VAR": ""})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(issues, "empty-value") {
		t.Error("expected empty-value issue")
	}
}

func TestLint_DetectsLowercaseKey(t *testing.T) {
	l, _ := New(DefaultRules())
	issues, err := l.Lint(map[string]string{"my_var": "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(issues, "lowercase-key") {
		t.Error("expected lowercase-key issue")
	}
}

func TestLint_DetectsWhitespaceValue(t *testing.T) {
	l, _ := New(DefaultRules())
	issues, err := l.Lint(map[string]string{"MY_VAR": "  value  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !containsRule(issues, "whitespace-value") {
		t.Error("expected whitespace-value issue")
	}
}

func TestLint_CleanEnvReturnsNoIssues(t *testing.T) {
	l, _ := New(DefaultRules())
	issues, err := l.Lint(map[string]string{"MY_VAR": "clean", "OTHER": "value"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(issues) != 0 {
		t.Fatalf("expected 0 issues, got %d: %v", len(issues), issues)
	}
}

func containsRule(issues []Issue, rule string) bool {
	for _, i := range issues {
		if i.Rule == rule {
			return true
		}
	}
	return false
}
