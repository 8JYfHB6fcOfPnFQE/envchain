package tagger_test

import (
	"testing"

	"envchain/internal/tagger"
)

func TestTag_BlankKeyReturnsError(t *testing.T) {
	tr := tagger.New()
	if err := tr.Tag("", "env", "production"); err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestTag_BlankNameReturnsError(t *testing.T) {
	tr := tagger.New()
	if err := tr.Tag("DB_HOST", "", "production"); err == nil {
		t.Fatal("expected error for blank tag name")
	}
}

func TestTag_AttachesTagToKey(t *testing.T) {
	tr := tagger.New()
	if err := tr.Tag("DB_HOST", "env", "production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	tags := tr.GetTags("DB_HOST")
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Name != "env" || tags[0].Value != "production" {
		t.Errorf("unexpected tag: %+v", tags[0])
	}
}

func TestTag_MultipleTagsOnSameKey(t *testing.T) {
	tr := tagger.New()
	_ = tr.Tag("API_KEY", "sensitivity", "high")
	_ = tr.Tag("API_KEY", "env", "staging")
	tags := tr.GetTags("API_KEY")
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
}

func TestGetTags_UnknownKeyReturnsNil(t *testing.T) {
	tr := tagger.New()
	if tags := tr.GetTags("UNKNOWN"); tags != nil {
		t.Errorf("expected nil, got %v", tags)
	}
}

func TestGetTags_ReturnsCopy(t *testing.T) {
	tr := tagger.New()
	_ = tr.Tag("PORT", "type", "network")
	tags := tr.GetTags("PORT")
	tags[0].Name = "mutated"
	original := tr.GetTags("PORT")
	if original[0].Name == "mutated" {
		t.Error("GetTags should return a copy, not a reference")
	}
}

func TestKeysWithTag_ReturnsSortedKeys(t *testing.T) {
	tr := tagger.New()
	_ = tr.Tag("Z_KEY", "group", "a")
	_ = tr.Tag("A_KEY", "group", "b")
	_ = tr.Tag("M_KEY", "group", "c")
	keys := tr.KeysWithTag("group")
	if len(keys) != 3 {
		t.Fatalf("expected 3 keys, got %d", len(keys))
	}
	if keys[0] != "A_KEY" || keys[1] != "M_KEY" || keys[2] != "Z_KEY" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestRemove_ClearsTagsForKey(t *testing.T) {
	tr := tagger.New()
	_ = tr.Tag("SECRET", "sensitivity", "high")
	tr.Remove("SECRET")
	if tags := tr.GetTags("SECRET"); tags != nil {
		t.Errorf("expected nil after removal, got %v", tags)
	}
}
