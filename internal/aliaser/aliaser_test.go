package aliaser

import (
	"testing"
)

func TestRegister_ValidAlias(t *testing.T) {
	a := New()
	if err := a.Register("DATABASE_URL", "DB_URL", "db_url"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRegister_BlankCanonicalReturnsError(t *testing.T) {
	a := New()
	if err := a.Register("", "DB_URL"); err == nil {
		t.Fatal("expected error for blank canonical key")
	}
}

func TestRegister_BlankAliasReturnsError(t *testing.T) {
	a := New()
	if err := a.Register("DATABASE_URL", ""); err == nil {
		t.Fatal("expected error for blank alias")
	}
}

func TestResolve_KnownAlias(t *testing.T) {
	a := New()
	_ = a.Register("DATABASE_URL", "DB_URL")

	canonical, ok := a.Resolve("DB_URL")
	if !ok {
		t.Fatal("expected alias to resolve")
	}
	if canonical != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", canonical)
	}
}

func TestResolve_CaseInsensitive(t *testing.T) {
	a := New()
	_ = a.Register("DATABASE_URL", "db_url")

	canonical, ok := a.Resolve("DB_URL")
	if !ok {
		t.Fatal("expected case-insensitive resolution")
	}
	if canonical != "DATABASE_URL" {
		t.Errorf("expected DATABASE_URL, got %s", canonical)
	}
}

func TestResolve_UnknownAliasReturnsFalse(t *testing.T) {
	a := New()
	_, ok := a.Resolve("UNKNOWN_KEY")
	if ok {
		t.Fatal("expected false for unknown alias")
	}
}

func TestLookup_ResolvesAndFetchesValue(t *testing.T) {
	a := New()
	_ = a.Register("DATABASE_URL", "DB_URL")
	env := map[string]string{"DATABASE_URL": "postgres://localhost/mydb"}

	val, ok := a.Lookup("DB_URL", env)
	if !ok {
		t.Fatal("expected value to be found via alias")
	}
	if val != "postgres://localhost/mydb" {
		t.Errorf("unexpected value: %s", val)
	}
}

func TestLookup_UnknownAliasReturnsFalse(t *testing.T) {
	a := New()
	env := map[string]string{"DATABASE_URL": "postgres://localhost/mydb"}
	_, ok := a.Lookup("DB_URL", env)
	if ok {
		t.Fatal("expected false for unregistered alias")
	}
}

func TestAliases_ReturnsRegisteredAliases(t *testing.T) {
	a := New()
	_ = a.Register("DATABASE_URL", "DB_URL", "db_connection")

	aliases := a.Aliases("DATABASE_URL")
	if len(aliases) != 2 {
		t.Errorf("expected 2 aliases, got %d: %v", len(aliases), aliases)
	}
}
