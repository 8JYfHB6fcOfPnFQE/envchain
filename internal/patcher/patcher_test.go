package patcher_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/patcher"
)

func TestNew_BlankKeyReturnsError(t *testing.T) {
	_, err := patcher.New([]patcher.Op{{Kind: patcher.OpSet, Key: "", Value: "v"}})
	if err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestNew_UnknownKindReturnsError(t *testing.T) {
	_, err := patcher.New([]patcher.Op{{Kind: "upsert", Key: "FOO"}})
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestApply_NilEnvReturnsError(t *testing.T) {
	p, _ := patcher.New([]patcher.Op{{Kind: patcher.OpSet, Key: "A", Value: "1"}})
	_, err := p.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestApply_SetsNewKey(t *testing.T) {
	p, err := patcher.New([]patcher.Op{{Kind: patcher.OpSet, Key: "FOO", Value: "bar"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{}
	res, err := p.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", env["FOO"])
	}
	if len(res.Applied) != 1 {
		t.Errorf("expected 1 applied op, got %d", len(res.Applied))
	}
}

func TestApply_DeleteExistingKey(t *testing.T) {
	p, _ := patcher.New([]patcher.Op{{Kind: patcher.OpDelete, Key: "FOO"}})
	env := map[string]string{"FOO": "bar"}
	res, err := p.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := env["FOO"]; ok {
		t.Error("expected FOO to be deleted")
	}
	if len(res.Applied) != 1 || len(res.Skipped) != 0 {
		t.Errorf("unexpected result counts: applied=%d skipped=%d", len(res.Applied), len(res.Skipped))
	}
}

func TestApply_DeleteMissingKeyIsSkipped(t *testing.T) {
	p, _ := patcher.New([]patcher.Op{{Kind: patcher.OpDelete, Key: "MISSING"}})
	env := map[string]string{}
	res, err := p.Apply(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Skipped) != 1 {
		t.Errorf("expected 1 skipped op, got %d", len(res.Skipped))
	}
}

func TestOps_ReturnsCopy(t *testing.T) {
	ops := []patcher.Op{{Kind: patcher.OpSet, Key: "X", Value: "1"}}
	p, _ := patcher.New(ops)
	out := p.Ops()
	out[0].Key = "MUTATED"
	if p.Ops()[0].Key != "X" {
		t.Error("Ops should return an isolated copy")
	}
}
