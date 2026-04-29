package differ_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/differ"
)

func TestCompare_NilBaseReturnsError(t *testing.T) {
	d := differ.New()
	_, err := d.Compare(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil base, got nil")
	}
}

func TestCompare_NilNextReturnsError(t *testing.T) {
	d := differ.New()
	_, err := d.Compare(map[string]string{}, nil)
	if err == nil {
		t.Fatal("expected error for nil next, got nil")
	}
}

func TestCompare_NoChanges(t *testing.T) {
	d := differ.New()
	base := map[string]string{"FOO": "bar", "BAZ": "qux"}
	next := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := d.Compare(base, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasChanges() {
		t.Errorf("expected no changes, got %v", res.Entries())
	}
}

func TestCompare_DetectsAdded(t *testing.T) {
	d := differ.New()
	res, err := d.Compare(map[string]string{}, map[string]string{"NEW_KEY": "val"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := res.Entries()
	if len(entries) != 1 || entries[0].Kind != differ.Added || entries[0].Key != "NEW_KEY" {
		t.Errorf("expected one Added entry for NEW_KEY, got %v", entries)
	}
}

func TestCompare_DetectsRemoved(t *testing.T) {
	d := differ.New()
	res, err := d.Compare(map[string]string{"OLD_KEY": "v"}, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := res.Entries()
	if len(entries) != 1 || entries[0].Kind != differ.Removed || entries[0].Key != "OLD_KEY" {
		t.Errorf("expected one Removed entry for OLD_KEY, got %v", entries)
	}
}

func TestCompare_DetectsChanged(t *testing.T) {
	d := differ.New()
	res, err := d.Compare(map[string]string{"KEY": "old"}, map[string]string{"KEY": "new"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := res.Entries()
	if len(entries) != 1 || entries[0].Kind != differ.Changed {
		t.Errorf("expected one Changed entry, got %v", entries)
	}
	if entries[0].OldValue != "old" || entries[0].NewValue != "new" {
		t.Errorf("unexpected values: %+v", entries[0])
	}
}

func TestCompare_EntriesSortedByKey(t *testing.T) {
	d := differ.New()
	base := map[string]string{"Z_KEY": "1", "A_KEY": "2"}
	next := map[string]string{}
	res, err := d.Compare(base, next)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	entries := res.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Key != "A_KEY" || entries[1].Key != "Z_KEY" {
		t.Errorf("entries not sorted: %v", entries)
	}
}
