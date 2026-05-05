package sanitizer_test

import (
	"testing"

	"github.com/example/envchain/internal/sanitizer"
)

func TestSanitizeKey_Uppercase(t *testing.T) {
	s := sanitizer.New(false)
	got, err := s.SanitizeKey("my_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MY_KEY" {
		t.Errorf("expected MY_KEY, got %s", got)
	}
}

func TestSanitizeKey_ReplacesInvalidChars(t *testing.T) {
	s := sanitizer.New(false)
	got, err := s.SanitizeKey("my-key.name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "MY_KEY_NAME" {
		t.Errorf("expected MY_KEY_NAME, got %s", got)
	}
}

func TestSanitizeKey_EmptyAfterSanitization(t *testing.T) {
	s := sanitizer.New(false)
	_, err := s.SanitizeKey("!!!")
	if err == nil {
		t.Fatal("expected error for key that sanitizes to empty string")
	}
}

func TestSanitizeKey_TrimsWhitespace(t *testing.T) {
	s := sanitizer.New(false)
	got, err := s.SanitizeKey("  APP_ENV  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "APP_ENV" {
		t.Errorf("expected APP_ENV, got %s", got)
	}
}

func TestSanitizeValue_TrimsWhitespace(t *testing.T) {
	s := sanitizer.New(false)
	got := s.SanitizeValue("  hello world  ")
	if got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestSanitizeValue_StripControlChars(t *testing.T) {
	s := sanitizer.New(true)
	got := s.SanitizeValue("hello\tworld\n")
	if got != "helloworld" {
		t.Errorf("expected 'helloworld', got %q", got)
	}
}

func TestSanitizeValue_NoStripControlChars(t *testing.T) {
	s := sanitizer.New(false)
	got := s.SanitizeValue("hello\tworld")
	if got != "hello\tworld" {
		t.Errorf("expected tab preserved, got %q", got)
	}
}

func TestSanitizeMap_AllValid(t *testing.T) {
	s := sanitizer.New(true)
	input := map[string]string{
		"db_host": "  localhost  ",
		"db-port": "5432\n",
	}
	out, err := s.SanitizeMap(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" {
		t.Errorf("expected localhost, got %q", out["DB_HOST"])
	}
	if out["DB_PORT"] != "5432" {
		t.Errorf("expected 5432, got %q", out["DB_PORT"])
	}
}

func TestSanitizeMap_SkipsInvalidKey(t *testing.T) {
	s := sanitizer.New(false)
	input := map[string]string{
		"!!!": "value",
		"valid_key": "ok",
	}
	out, err := s.SanitizeMap(input)
	if err == nil {
		t.Fatal("expected error for invalid key")
	}
	if _, ok := out["VALID_KEY"]; !ok {
		t.Error("expected VALID_KEY to be present in output")
	}
	if len(out) != 1 {
		t.Errorf("expected 1 entry, got %d", len(out))
	}
}
