package flattener

import (
	"testing"
)

func TestNew_ValidStrategy(t *testing.T) {
	f, err := New(CollisionKeepFirst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Flattener")
	}
}

func TestNew_InvalidStrategyReturnsError(t *testing.T) {
	_, err := New(CollisionStrategy(99))
	if err == nil {
		t.Fatal("expected error for unknown strategy")
	}
}

func TestAdd_NilEnvReturnsError(t *testing.T) {
	f, _ := New(CollisionKeepFirst)
	err := f.Add("SVC", nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestFlatten_NoSourcesReturnsError(t *testing.T) {
	f, _ := New(CollisionKeepFirst)
	_, err := f.Flatten()
	if err == nil {
		t.Fatal("expected error when no sources registered")
	}
}

func TestFlatten_SingleSourceNoPrefix(t *testing.T) {
	f, _ := New(CollisionKeepFirst)
	_ = f.Add("", map[string]string{"db_host": "localhost", "db_port": "5432"})

	got, err := f.Flatten()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %q", got["DB_HOST"])
	}
	if got["DB_PORT"] != "5432" {
		t.Errorf("expected DB_PORT=5432, got %q", got["DB_PORT"])
	}
}

func TestFlatten_PrefixNamespacesKeys(t *testing.T) {
	f, _ := New(CollisionKeepFirst)
	_ = f.Add("app", map[string]string{"name": "envchain"})

	got, err := f.Flatten()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["APP_NAME"] != "envchain" {
		t.Errorf("expected APP_NAME=envchain, got %q", got["APP_NAME"])
	}
}

func TestFlatten_CollisionKeepFirst(t *testing.T) {
	f, _ := New(CollisionKeepFirst)
	_ = f.Add("", map[string]string{"HOST": "first"})
	_ = f.Add("", map[string]string{"HOST": "second"})

	got, _ := f.Flatten()
	if got["HOST"] != "first" {
		t.Errorf("expected HOST=first, got %q", got["HOST"])
	}
}

func TestFlatten_CollisionKeepLast(t *testing.T) {
	f, _ := New(CollisionKeepLast)
	_ = f.Add("", map[string]string{"HOST": "first"})
	_ = f.Add("", map[string]string{"HOST": "second"})

	got, _ := f.Flatten()
	if got["HOST"] != "second" {
		t.Errorf("expected HOST=second, got %q", got["HOST"])
	}
}

func TestFlatten_CollisionErrorReturnsError(t *testing.T) {
	f, _ := New(CollisionError)
	_ = f.Add("", map[string]string{"HOST": "first"})
	_ = f.Add("", map[string]string{"HOST": "second"})

	_, err := f.Flatten()
	if err == nil {
		t.Fatal("expected collision error")
	}
}
