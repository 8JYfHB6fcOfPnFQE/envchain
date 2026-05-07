package labeler_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/labeler"
)

func TestAttach_ValidLabel(t *testing.T) {
	l := labeler.New()
	if err := l.Attach("DB_HOST", "env", "production"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	lbls := l.Get("DB_HOST")
	if len(lbls) != 1 {
		t.Fatalf("expected 1 label, got %d", len(lbls))
	}
	if lbls[0].Name != "env" || lbls[0].Value != "production" {
		t.Errorf("unexpected label: %+v", lbls[0])
	}
}

func TestAttach_BlankKeyReturnsError(t *testing.T) {
	l := labeler.New()
	if err := l.Attach("", "env", "production"); err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestAttach_BlankNameReturnsError(t *testing.T) {
	l := labeler.New()
	if err := l.Attach("DB_HOST", "", "production"); err == nil {
		t.Fatal("expected error for blank name")
	}
}

func TestAttach_BlankValueReturnsError(t *testing.T) {
	l := labeler.New()
	if err := l.Attach("DB_HOST", "env", ""); err == nil {
		t.Fatal("expected error for blank value")
	}
}

func TestAttach_CaseInsensitiveKey(t *testing.T) {
	l := labeler.New()
	_ = l.Attach("db_host", "tier", "backend")
	lbls := l.Get("DB_HOST")
	if len(lbls) != 1 {
		t.Fatalf("expected 1 label via normalised key, got %d", len(lbls))
	}
}

func TestGet_UnknownKeyReturnsNil(t *testing.T) {
	l := labeler.New()
	if got := l.Get("UNKNOWN"); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestGet_ReturnsCopy(t *testing.T) {
	l := labeler.New()
	_ = l.Attach("API_KEY", "sensitivity", "high")
	copy1 := l.Get("API_KEY")
	copy1[0].Value = "mutated"
	copy2 := l.Get("API_KEY")
	if copy2[0].Value == "mutated" {
		t.Error("Get should return an isolated copy")
	}
}

func TestFindByLabel_ReturnsMatchingKeys(t *testing.T) {
	l := labeler.New()
	_ = l.Attach("DB_HOST", "env", "production")
	_ = l.Attach("DB_PORT", "env", "production")
	_ = l.Attach("DEBUG", "env", "development")

	keys := l.FindByLabel("env", "production")
	if len(keys) != 2 {
		t.Fatalf("expected 2 keys, got %d: %v", len(keys), keys)
	}
}

func TestSummary_NoLabels(t *testing.T) {
	l := labeler.New()
	s := l.Summary("MISSING_KEY")
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestSummary_WithLabels(t *testing.T) {
	l := labeler.New()
	_ = l.Attach("TOKEN", "sensitivity", "high")
	_ = l.Attach("TOKEN", "env", "production")
	s := l.Summary("TOKEN")
	if s == "" {
		t.Error("expected non-empty summary")
	}
}
