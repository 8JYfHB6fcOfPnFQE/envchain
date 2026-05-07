package tracer_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/tracer"
)

// TestIntegration_TraceMultipleContextsAndFilter simulates recording resolution
// events across multiple contexts and verifies filtering and summary output.
func TestIntegration_TraceMultipleContextsAndFilter(t *testing.T) {
	tr := tracer.New()

	contexts := []struct {
		ctx, key, action, detail string
	}{
		{"base", "DB_HOST", "set", "localhost"},
		{"base", "PORT", "set", "5432"},
		{"staging", "DB_HOST", "overridden", "staging.db.internal"},
		{"staging", "LOG_LEVEL", "set", "debug"},
		{"production", "DB_HOST", "overridden", "prod.db.internal"},
		{"production", "PORT", "inherited", "from base"},
	}

	for _, ev := range contexts {
		if err := tr.Record(ev.ctx, ev.key, ev.action, ev.detail); err != nil {
			t.Fatalf("unexpected error recording event: %v", err)
		}
	}

	all := tr.Events()
	if len(all) != len(contexts) {
		t.Fatalf("expected %d events, got %d", len(contexts), len(all))
	}

	baseEvents := tr.FilterByContext("base")
	if len(baseEvents) != 2 {
		t.Errorf("expected 2 base events, got %d", len(baseEvents))
	}

	prodEvents := tr.FilterByContext("production")
	if len(prodEvents) != 2 {
		t.Errorf("expected 2 production events, got %d", len(prodEvents))
	}

	summary := tr.Summary()
	for _, ctx := range []string{"base", "staging", "production"} {
		if !strings.Contains(summary, ctx) {
			t.Errorf("summary missing context %q", ctx)
		}
	}

	if !strings.Contains(summary, "6 trace event(s)") {
		t.Errorf("summary should report 6 events, got:\n%s", summary)
	}
}
