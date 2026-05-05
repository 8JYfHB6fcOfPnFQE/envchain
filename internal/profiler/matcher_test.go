package profiler

import (
	"strings"
	"testing"
)

func makeProfile(t *testing.T, name string, keys []string) *Profile {
	t.Helper()
	p := New()
	if err := p.Define(name, keys, nil); err != nil {
		t.Fatalf("Define: %v", err)
	}
	pr, _ := p.Get(name)
	return pr
}

func TestMatch_NilProfileReturnsError(t *testing.T) {
	_, err := Match(nil, map[string]string{"A": "1"})
	if err == nil {
		t.Fatal("expected error for nil profile")
	}
}

func TestMatch_AllKeysPresent(t *testing.T) {
	pr := makeProfile(t, "prod", []string{"DB_URL", "API_KEY"})
	env := map[string]string{"DB_URL": "x", "API_KEY": "y"}
	res, err := Match(pr, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsComplete() {
		t.Errorf("expected complete match, missing: %v", res.Missing)
	}
	if len(res.Matched) != 2 {
		t.Errorf("expected 2 matched, got %d", len(res.Matched))
	}
}

func TestMatch_DetectsMissingKeys(t *testing.T) {
	pr := makeProfile(t, "prod", []string{"DB_URL", "API_KEY", "SECRET"})
	env := map[string]string{"DB_URL": "x"}
	res, err := Match(pr, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.IsComplete() {
		t.Error("expected incomplete match")
	}
	if len(res.Missing) != 2 {
		t.Errorf("expected 2 missing, got %d", len(res.Missing))
	}
}

func TestMatch_DetectsExtraKeys(t *testing.T) {
	pr := makeProfile(t, "staging", []string{"DB_URL"})
	env := map[string]string{"DB_URL": "x", "EXTRA_KEY": "y"}
	res, err := Match(pr, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Extra) != 1 {
		t.Errorf("expected 1 extra key, got %d", len(res.Extra))
	}
}

func TestMatch_CaseInsensitiveKeyComparison(t *testing.T) {
	pr := makeProfile(t, "dev", []string{"DB_URL"})
	env := map[string]string{"db_url": "value"}
	res, err := Match(pr, env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.IsComplete() {
		t.Errorf("expected case-insensitive match to succeed, missing: %v", res.Missing)
	}
}

func TestMatchResult_Summary(t *testing.T) {
	pr := makeProfile(t, "prod", []string{"A", "B"})
	env := map[string]string{"A": "1"}
	res, _ := Match(pr, env)
	s := res.Summary()
	if !strings.Contains(s, "profile=prod") {
		t.Errorf("summary missing profile name: %q", s)
	}
	if !strings.Contains(s, "missing=1") {
		t.Errorf("summary missing count wrong: %q", s)
	}
}
