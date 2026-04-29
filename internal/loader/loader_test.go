package loader_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envchain/internal/loader"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, ".env")
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

func TestLoad_BasicKeyValue(t *testing.T) {
	path := writeTempFile(t, "FOO=bar\nBAZ=qux\n")
	l := loader.New(path)
	got, err := l.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestLoad_IgnoresCommentsAndBlanks(t *testing.T) {
	path := writeTempFile(t, "# comment\n\nKEY=value\n")
	l := loader.New(path)
	got, err := l.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got["KEY"] != "value" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestLoad_StripQuotes(t *testing.T) {
	path := writeTempFile(t, `QUOTED="hello world"` + "\nSINGLE='world'\n")
	l := loader.New(path)
	got, err := l.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["QUOTED"] != "hello world" {
		t.Errorf("expected 'hello world', got %q", got["QUOTED"])
	}
	if got["SINGLE"] != "world" {
		t.Errorf("expected 'world', got %q", got["SINGLE"])
	}
}

func TestLoad_InvalidLine(t *testing.T) {
	path := writeTempFile(t, "NOEQUALS\n")
	l := loader.New(path)
	_, err := l.Load()
	if err == nil {
		t.Fatal("expected error for invalid line, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	l := loader.New("/nonexistent/.env")
	_, err := l.Load()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_EmptyValue(t *testing.T) {
	path := writeTempFile(t, "EMPTY=\nALSO_EMPTY=\"\"\n")
	l := loader.New(path)
	got, err := l.Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["EMPTY"] != "" {
		t.Errorf("expected empty string for EMPTY, got %q", got["EMPTY"])
	}
	if got["ALSO_EMPTY"] != "" {
		t.Errorf("expected empty string for ALSO_EMPTY, got %q", got["ALSO_EMPTY"])
	}
}
