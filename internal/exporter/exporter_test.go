package exporter_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/exporter"
)

func TestNew_ValidFormat(t *testing.T) {
	for _, f := range []exporter.Format{exporter.FormatDotenv, exporter.FormatExport, exporter.FormatJSON} {
		_, err := exporter.New(f, &strings.Builder{})
		if err != nil {
			t.Errorf("expected no error for format %q, got %v", f, err)
		}
	}
}

func TestNew_InvalidFormat(t *testing.T) {
	_, err := exporter.New("yaml", &strings.Builder{})
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestWrite_Dotenv(t *testing.T) {
	var sb strings.Builder
	e, _ := exporter.New(exporter.FormatDotenv, &sb)
	err := e.Write(map[string]string{"FOO": "bar", "BAZ": "qux"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := sb.String()
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("unexpected dotenv output: %q", out)
	}
}

func TestWrite_Export(t *testing.T) {
	var sb strings.Builder
	e, _ := exporter.New(exporter.FormatExport, &sb)
	_ = e.Write(map[string]string{"KEY": "value"})
	out := sb.String()
	if !strings.Contains(out, "export KEY=value") {
		t.Errorf("expected export prefix, got: %q", out)
	}
}

func TestWrite_JSON(t *testing.T) {
	var sb strings.Builder
	e, _ := exporter.New(exporter.FormatJSON, &sb)
	_ = e.Write(map[string]string{"A": "1"})
	out := sb.String()
	if !strings.Contains(out, `"A": "1"`) {
		t.Errorf("unexpected JSON output: %q", out)
	}
}

func TestWrite_DotenvQuotesSpaces(t *testing.T) {
	var sb strings.Builder
	e, _ := exporter.New(exporter.FormatDotenv, &sb)
	_ = e.Write(map[string]string{"MSG": "hello world"})
	out := sb.String()
	if !strings.Contains(out, `"hello world"`) {
		t.Errorf("expected quoted value with spaces, got: %q", out)
	}
}

func TestWrite_SortedOutput(t *testing.T) {
	var sb strings.Builder
	e, _ := exporter.New(exporter.FormatDotenv, &sb)
	_ = e.Write(map[string]string{"Z": "last", "A": "first", "M": "mid"})
	out := sb.String()
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 || !strings.HasPrefix(lines[0], "A=") || !strings.HasPrefix(lines[2], "Z=") {
		t.Errorf("expected sorted output, got: %v", lines)
	}
}
