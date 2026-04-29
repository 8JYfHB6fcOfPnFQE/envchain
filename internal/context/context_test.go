package context_test

import (
	"testing"

	"github.com/user/envchain/internal/context"
)

func TestRegister_AndGet(t *testing.T) {
	reg := context.New()
	ctx := &context.Context{
		Name: "production",
		EnvSets: map[string]context.EnvSetRef{
			"db": {SetName: "db", FilePath: "/etc/envchain/db.env"},
		},
	}
	if err := reg.Register(ctx); err != nil {
		t.Fatalf("unexpected error registering context: %v", err)
	}
	got, err := reg.Get("production")
	if err != nil {
		t.Fatalf("unexpected error getting context: %v", err)
	}
	if got.Name != "production" {
		t.Errorf("expected name %q, got %q", "production", got.Name)
	}
}

func TestRegister_NilContext(t *testing.T) {
	reg := context.New()
	if err := reg.Register(nil); err == nil {
		t.Error("expected error for nil context, got nil")
	}
}

func TestRegister_EmptyName(t *testing.T) {
	reg := context.New()
	ctx := &context.Context{Name: ""}
	if err := reg.Register(ctx); err == nil {
		t.Error("expected error for empty context name, got nil")
	}
}

func TestGet_NotFound(t *testing.T) {
	reg := context.New()
	_, err := reg.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing context, got nil")
	}
}

func TestList_ReturnsAllNames(t *testing.T) {
	reg := context.New()
	for _, name := range []string{"staging", "production", "dev"} {
		_ = reg.Register(&context.Context{Name: name})
	}
	names := reg.List()
	if len(names) != 3 {
		t.Errorf("expected 3 contexts, got %d", len(names))
	}
}

func TestRemove_DeletesContext(t *testing.T) {
	reg := context.New()
	_ = reg.Register(&context.Context{Name: "staging"})
	reg.Remove("staging")
	_, err := reg.Get("staging")
	if err == nil {
		t.Error("expected error after removal, got nil")
	}
}
