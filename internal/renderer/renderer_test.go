package renderer_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/exporter"
	"github.com/yourorg/envchain/internal/renderer"
)

func TestRender_Dotenv(t *testing.T) {
	r := renderer.New(map[string]string{"HOST": "localhost", "PORT": "5432"})
	var sb strings.Builder
	if err := r.Render(exporter.FormatDotenv, &sb); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "HOST=localhost") {
		t.Errorf("missing HOST in output: %q", out)
	}
}

func TestRender_MasksKey(t *testing.T) {
	r := renderer.New(map[string]string{"SECRET": "supersecret", "APP": "myapp"})
	r.Mask("SECRET")
	var sb strings.Builder
	_ = r.Render(exporter.FormatDotenv, &sb)
	out := sb.String()
	if strings.Contains(out, "supersecret") {
		t.Errorf("secret value should be masked, got: %q", out)
	}
	if !strings.Contains(out, "SECRET=***") {
		t.Errorf("expected masked placeholder, got: %q", out)
	}
}

func TestRender_MaskCaseInsensitive(t *testing.T) {
	r := renderer.New(map[string]string{"DB_PASS": "hunter2"})
	r.Mask("db_pass")
	var sb strings.Builder
	_ = r.Render(exporter.FormatDotenv, &sb)
	out := sb.String()
	if strings.Contains(out, "hunter2") {
		t.Errorf("expected value to be masked regardless of case")
	}
}

func TestRender_InvalidFormat(t *testing.T) {
	r := renderer.New(map[string]string{"X": "1"})
	var sb strings.Builder
	err := r.Render("toml", &sb)
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestKeys_ReturnsAllKeys(t *testing.T) {
	vars := map[string]string{"A": "1", "B": "2", "C": "3"}
	r := renderer.New(vars)
	keys := r.Keys()
	if len(keys) != 3 {
		t.Errorf("expected 3 keys, got %d", len(keys))
	}
}

func TestRender_DoesNotMutateOriginal(t *testing.T) {
	orig := map[string]string{"TOKEN": "abc123"}
	r := renderer.New(orig)
	r.Mask("TOKEN")
	var sb strings.Builder
	_ = r.Render(exporter.FormatDotenv, &sb)
	if orig["TOKEN"] != "abc123" {
		t.Error("original map should not be mutated by Mask")
	}
}
