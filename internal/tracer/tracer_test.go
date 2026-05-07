package tracer

import (
	"strings"
	"testing"
)

func TestRecord_ValidEvent(t *testing.T) {
	tr := New()
	if err := tr.Record("production", "DB_URL", "resolved", "from base"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	events := tr.Events()
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].Context != "production" {
		t.Errorf("expected context 'production', got %q", events[0].Context)
	}
	if events[0].Key != "DB_URL" {
		t.Errorf("expected key 'DB_URL', got %q", events[0].Key)
	}
}

func TestRecord_EmptyContextReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Record("", "KEY", "resolved", ""); err == nil {
		t.Fatal("expected error for empty context")
	}
}

func TestRecord_EmptyActionReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Record("staging", "KEY", "", ""); err == nil {
		t.Fatal("expected error for empty action")
	}
}

func TestEvents_ReturnsCopy(t *testing.T) {
	tr := New()
	_ = tr.Record("ctx", "K", "set", "")
	events := tr.Events()
	events[0].Context = "mutated"
	original := tr.Events()
	if original[0].Context == "mutated" {
		t.Error("Events() should return a copy, not a reference")
	}
}

func TestFilterByContext_ReturnsMatchingEvents(t *testing.T) {
	tr := New()
	_ = tr.Record("prod", "A", "set", "")
	_ = tr.Record("staging", "B", "set", "")
	_ = tr.Record("prod", "C", "override", "")

	prodEvents := tr.FilterByContext("prod")
	if len(prodEvents) != 2 {
		t.Fatalf("expected 2 prod events, got %d", len(prodEvents))
	}
	for _, e := range prodEvents {
		if e.Context != "prod" {
			t.Errorf("unexpected context %q in filtered results", e.Context)
		}
	}
}

func TestFilterByContext_NoMatch(t *testing.T) {
	tr := New()
	_ = tr.Record("prod", "A", "set", "")
	result := tr.FilterByContext("nonexistent")
	if len(result) != 0 {
		t.Errorf("expected 0 events, got %d", len(result))
	}
}

func TestSummary_EmptyTracer(t *testing.T) {
	tr := New()
	summary := tr.Summary()
	if summary != "no trace events recorded" {
		t.Errorf("unexpected summary: %q", summary)
	}
}

func TestSummary_ContainsEventInfo(t *testing.T) {
	tr := New()
	_ = tr.Record("prod", "SECRET_KEY", "resolved", "from overlay")
	summary := tr.Summary()
	if !strings.Contains(summary, "prod") {
		t.Error("summary should contain context name")
	}
	if !strings.Contains(summary, "SECRET_KEY") {
		t.Error("summary should contain key name")
	}
	if !strings.Contains(summary, "resolved") {
		t.Error("summary should contain action")
	}
}
