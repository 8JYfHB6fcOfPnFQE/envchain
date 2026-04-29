package resolver_test

import (
	"testing"

	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/context"
	"github.com/user/envchain/internal/merger"
	"github.com/user/envchain/internal/resolver"
	"github.com/user/envchain/internal/validator"
)

func makeRegistry(t *testing.T, entries map[string]map[string]string) *context.Registry {
	t.Helper()
	reg := context.New()
	for name, env := range entries {
		if err := reg.Register(name, env); err != nil {
			t.Fatalf("register %q: %v", name, err)
		}
	}
	return reg
}

func TestNew_NilChain(t *testing.T) {
	m := merger.New()
	v := validator.New(nil, nil)
	_, err := resolver.New(nil, m, v)
	if err == nil {
		t.Fatal("expected error for nil chain")
	}
}

func TestNew_NilMerger(t *testing.T) {
	reg := makeRegistry(t, nil)
	c := chain.New(reg)
	v := validator.New(nil, nil)
	_, err := resolver.New(c, nil, v)
	if err == nil {
		t.Fatal("expected error for nil merger")
	}
}

func TestNew_NilValidator(t *testing.T) {
	reg := makeRegistry(t, nil)
	c := chain.New(reg)
	m := merger.New()
	_, err := resolver.New(c, m, nil)
	if err == nil {
		t.Fatal("expected error for nil validator")
	}
}

func TestResolve_MergesContextsInOrder(t *testing.T) {
	reg := makeRegistry(t, map[string]map[string]string{
		"base": {"APP_ENV": "production", "LOG_LEVEL": "info"},
		"override": {"LOG_LEVEL": "debug", "FEATURE_X": "true"},
	})

	c := chain.New(reg)
	_ = c.Add("base")
	_ = c.Add("override")

	m := merger.New()
	v := validator.New(nil, nil)

	res, err := resolver.New(c, m, v)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := res.Resolve()
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}

	if result["LOG_LEVEL"] != "debug" {
		t.Errorf("expected LOG_LEVEL=debug, got %q", result["LOG_LEVEL"])
	}
	if result["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV=production, got %q", result["APP_ENV"])
	}
	if result["FEATURE_X"] != "true" {
		t.Errorf("expected FEATURE_X=true, got %q", result["FEATURE_X"])
	}
}

func TestResolve_ValidationFailure(t *testing.T) {
	reg := makeRegistry(t, map[string]map[string]string{
		"base": {"APP_ENV": "staging"},
	})

	c := chain.New(reg)
	_ = c.Add("base")

	m := merger.New()
	v := validator.New([]string{"MISSING_KEY"}, nil)

	res, err := resolver.New(c, m, v)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	_, err = res.Resolve()
	if err == nil {
		t.Fatal("expected validation error for missing required key")
	}
}
