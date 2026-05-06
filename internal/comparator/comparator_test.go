package comparator_test

import (
	"testing"

	"github.com/example/envchain/internal/comparator"
)

func TestNew_ValidNames(t *testing.T) {
	c, err := comparator.New("staging", "production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.LeftName() != "staging" || c.RightName() != "production" {
		t.Errorf("unexpected names: %s / %s", c.LeftName(), c.RightName())
	}
}

func TestNew_EmptyLeftReturnsError(t *testing.T) {
	_, err := comparator.New("", "production")
	if err == nil {
		t.Fatal("expected error for empty leftName")
	}
}

func TestNew_EmptyRightReturnsError(t *testing.T) {
	_, err := comparator.New("staging", "")
	if err == nil {
		t.Fatal("expected error for empty rightName")
	}
}

func TestCompare_NilLeftReturnsError(t *testing.T) {
	c, _ := comparator.New("a", "b")
	_, err := c.Compare(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil left")
	}
}

func TestCompare_NilRightReturnsError(t *testing.T) {
	c, _ := comparator.New("a", "b")
	_, err := c.Compare(map[string]string{}, nil)
	if err == nil {
		t.Fatal("expected error for nil right")
	}
}

func TestCompare_DetectsMatching(t *testing.T) {
	c, _ := comparator.New("a", "b")
	left := map[string]string{"FOO": "bar", "BAZ": "qux"}
	right := map[string]string{"FOO": "bar", "BAZ": "qux"}
	res, err := c.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Matching) != 2 {
		t.Errorf("expected 2 matching, got %d", len(res.Matching))
	}
	if len(res.Conflicts) != 0 || len(res.LeftOnly) != 0 || len(res.RightOnly) != 0 {
		t.Error("expected no conflicts or exclusive keys")
	}
}

func TestCompare_DetectsConflicts(t *testing.T) {
	c, _ := comparator.New("a", "b")
	left := map[string]string{"FOO": "old"}
	right := map[string]string{"FOO": "new"}
	res, err := c.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := res.Conflicts["FOO"]
	if !ok {
		t.Fatal("expected conflict for FOO")
	}
	if v[0] != "old" || v[1] != "new" {
		t.Errorf("unexpected conflict values: %v", v)
	}
}

func TestCompare_DetectsLeftAndRightOnly(t *testing.T) {
	c, _ := comparator.New("a", "b")
	left := map[string]string{"ONLY_LEFT": "1"}
	right := map[string]string{"ONLY_RIGHT": "2"}
	res, err := c.Compare(left, right)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.LeftOnly) != 1 || res.LeftOnly[0] != "ONLY_LEFT" {
		t.Errorf("unexpected LeftOnly: %v", res.LeftOnly)
	}
	if len(res.RightOnly) != 1 || res.RightOnly[0] != "ONLY_RIGHT" {
		t.Errorf("unexpected RightOnly: %v", res.RightOnly)
	}
}
