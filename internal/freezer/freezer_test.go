package freezer

import (
	"testing"
)

func TestFreeze_ValidLabel(t *testing.T) {
	f := New()
	env := map[string]string{"APP_ENV": "production", "PORT": "8080"}
	if err := f.Freeze("prod", env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFreeze_EmptyLabelReturnsError(t *testing.T) {
	f := New()
	err := f.Freeze("", map[string]string{"KEY": "val"})
	if err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestFreeze_NilEnvReturnsError(t *testing.T) {
	f := New()
	err := f.Freeze("label", nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestThaw_ReturnsFrozenValues(t *testing.T) {
	f := New()
	env := map[string]string{"KEY": "value"}
	_ = f.Freeze("snap", env)

	out, err := f.Thaw("snap")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected 'value', got %q", out["KEY"])
	}
}

func TestThaw_MutationIsolation(t *testing.T) {
	f := New()
	env := map[string]string{"KEY": "original"}
	_ = f.Freeze("snap", env)

	out, _ := f.Thaw("snap")
	out["KEY"] = "mutated"

	out2, _ := f.Thaw("snap")
	if out2["KEY"] != "original" {
		t.Errorf("frozen frame was mutated: got %q", out2["KEY"])
	}
}

func TestThaw_UnknownLabelReturnsError(t *testing.T) {
	f := New()
	_, err := f.Thaw("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown label")
	}
}

func TestLabels_ReturnsAllLabels(t *testing.T) {
	f := New()
	_ = f.Freeze("a", map[string]string{"X": "1"})
	_ = f.Freeze("b", map[string]string{"Y": "2"})

	labels := f.Labels()
	if len(labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(labels))
	}
}

func TestDrop_RemovesLabel(t *testing.T) {
	f := New()
	_ = f.Freeze("tmp", map[string]string{"K": "v"})
	if err := f.Drop("tmp"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(f.Labels()) != 0 {
		t.Error("expected no labels after drop")
	}
}

func TestDrop_UnknownLabelReturnsError(t *testing.T) {
	f := New()
	err := f.Drop("ghost")
	if err == nil {
		t.Fatal("expected error for unknown label")
	}
}
