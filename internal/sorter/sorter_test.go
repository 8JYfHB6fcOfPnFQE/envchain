package sorter_test

import (
	"testing"

	"github.com/example/envchain/internal/sorter"
)

func TestNew_ValidAscending(t *testing.T) {
	s, err := sorter.New(sorter.Ascending, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sorter")
	}
}

func TestNew_InvalidOrderReturnsError(t *testing.T) {
	_, err := sorter.New(sorter.Order(99), false)
	if err == nil {
		t.Fatal("expected error for invalid order")
	}
}

func TestSort_NilEnvReturnsError(t *testing.T) {
	s, _ := sorter.New(sorter.Ascending, false)
	_, err := s.Sort(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestSort_AscendingByKey(t *testing.T) {
	s, _ := sorter.New(sorter.Ascending, false)
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	pairs, err := s.Sort(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"APPLE", "MANGO", "ZEBRA"}
	for i, p := range pairs {
		if p.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, p.Key, expected[i])
		}
	}
}

func TestSort_DescendingByKey(t *testing.T) {
	s, _ := sorter.New(sorter.Descending, false)
	env := map[string]string{"ZEBRA": "1", "APPLE": "2", "MANGO": "3"}
	pairs, err := s.Sort(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"ZEBRA", "MANGO", "APPLE"}
	for i, p := range pairs {
		if p.Key != expected[i] {
			t.Errorf("position %d: got %q, want %q", i, p.Key, expected[i])
		}
	}
}

func TestSort_ByValueLength(t *testing.T) {
	s, _ := sorter.New(sorter.Ascending, true)
	env := map[string]string{"A": "hello world", "B": "hi", "C": "hey"}
	pairs, err := s.Sort(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// shortest value first: "hi"(2), "hey"(3), "hello world"(11)
	if pairs[0].Key != "B" {
		t.Errorf("expected B first, got %q", pairs[0].Key)
	}
	if pairs[2].Key != "A" {
		t.Errorf("expected A last, got %q", pairs[2].Key)
	}
}

func TestKeys_ReturnsSortedKeys(t *testing.T) {
	s, _ := sorter.New(sorter.Ascending, false)
	env := map[string]string{"Z": "1", "A": "2", "M": "3"}
	keys, err := s.Keys(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keys[0] != "A" || keys[1] != "M" || keys[2] != "Z" {
		t.Errorf("unexpected key order: %v", keys)
	}
}

func TestKeys_NilEnvReturnsError(t *testing.T) {
	s, _ := sorter.New(sorter.Ascending, false)
	_, err := s.Keys(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}
