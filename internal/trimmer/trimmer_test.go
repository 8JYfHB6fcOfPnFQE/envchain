package trimmer_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/trimmer"
)

func TestNew_ValidTrimmer(t *testing.T) {
	_, err := trimmer.New([]string{"SECRET"}, []string{"INTERNAL_"}, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNew_EmptyExcludeKeyReturnsError(t *testing.T) {
	_, err := trimmer.New([]string{""}, nil, false)
	if err == nil {
		t.Fatal("expected error for empty exclude key")
	}
}

func TestNew_EmptyExcludePrefixReturnsError(t *testing.T) {
	_, err := trimmer.New(nil, []string{""}, false)
	if err == nil {
		t.Fatal("expected error for empty exclude prefix")
	}
}

func TestTrim_NilEnvReturnsError(t *testing.T) {
	tr, _ := trimmer.New(nil, nil, false)
	_, err := tr.Trim(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestTrim_RemovesExcludedKey(t *testing.T) {
	tr, _ := trimmer.New([]string{"SECRET"}, nil, false)
	result, err := tr.Trim(map[string]string{"SECRET": "abc", "APP": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := result["SECRET"]; ok {
		t.Error("expected SECRET to be removed")
	}
	if result["APP"] != "1" {
		t.Error("expected APP to be retained")
	}
}

func TestTrim_CaseInsensitiveKeyMatch(t *testing.T) {
	tr, _ := trimmer.New([]string{"secret"}, nil, false)
	result, _ := tr.Trim(map[string]string{"SECRET": "val"})
	if _, ok := result["SECRET"]; ok {
		t.Error("expected SECRET to be removed via case-insensitive match")
	}
}

func TestTrim_RemovesByPrefix(t *testing.T) {
	tr, _ := trimmer.New(nil, []string{"INTERNAL_"}, false)
	result, _ := tr.Trim(map[string]string{"INTERNAL_KEY": "x", "PUBLIC": "y"})
	if _, ok := result["INTERNAL_KEY"]; ok {
		t.Error("expected INTERNAL_KEY to be removed")
	}
	if result["PUBLIC"] != "y" {
		t.Error("expected PUBLIC to be retained")
	}
}

func TestTrim_DropBlankValues(t *testing.T) {
	tr, _ := trimmer.New(nil, nil, true)
	result, _ := tr.Trim(map[string]string{"EMPTY": "", "SPACES": "   ", "OK": "value"})
	if _, ok := result["EMPTY"]; ok {
		t.Error("expected EMPTY to be removed")
	}
	if _, ok := result["SPACES"]; ok {
		t.Error("expected SPACES to be removed")
	}
	if result["OK"] != "value" {
		t.Error("expected OK to be retained")
	}
}

func TestTrim_DoesNotMutateOriginal(t *testing.T) {
	tr, _ := trimmer.New([]string{"DROP"}, nil, false)
	orig := map[string]string{"DROP": "gone", "KEEP": "here"}
	tr.Trim(orig)
	if _, ok := orig["DROP"]; !ok {
		t.Error("original map must not be mutated")
	}
}
