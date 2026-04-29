package chain_test

import (
	"testing"

	"github.com/envchain/envchain/internal/chain"
	"github.com/envchain/envchain/internal/context"
)

func newRegistry(t *testing.T) *context.Registry {
	t.Helper()
	reg := context.NewRegistry()
	return reg
}

func TestNew_EmptyChain(t *testing.T) {
	reg := newRegistry(t)
	c := chain.New(reg)
	if names := c.Names(); len(names) != 0 {
		t.Errorf("expected empty names, got %v", names)
	}
}

func TestAdd_ValidName(t *testing.T) {
	reg := newRegistry(t)
	c := chain.New(reg)
	if err := c.Add("production"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if names := c.Names(); len(names) != 1 || names[0] != "production" {
		t.Errorf("expected [production], got %v", names)
	}
}

func TestAdd_EmptyNameReturnsError(t *testing.T) {
	reg := newRegistry(t)
	c := chain.New(reg)
	if err := c.Add(""); err == nil {
		t.Fatal("expected error for empty name, got nil")
	}
}

func TestResolve_EmptyChainReturnsError(t *testing.T) {
	reg := newRegistry(t)
	c := chain.New(reg)
	_, err := c.Resolve()
	if err == nil {
		t.Fatal("expected ErrEmptyChain, got nil")
	}
	if err != chain.ErrEmptyChain {
		t.Errorf("expected ErrEmptyChain, got %v", err)
	}
}

func TestResolve_UnknownContextReturnsError(t *testing.T) {
	reg := newRegistry(t)
	c := chain.New(reg)
	_ = c.Add("nonexistent")
	_, err := c.Resolve()
	if err == nil {
		t.Fatal("expected error for unknown context, got nil")
	}
}

func TestNames_ReturnsCopy(t *testing.T) {
	reg := newRegistry(t)
	c := chain.New(reg)
	_ = c.Add("base")
	_ = c.Add("override")
	names := c.Names()
	names[0] = "mutated"
	if c.Names()[0] != "base" {
		t.Error("Names() should return a copy, not a reference")
	}
}
