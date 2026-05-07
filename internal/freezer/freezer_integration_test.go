package freezer_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/freezer"
)

// TestIntegration_FreezeThawRoundtrip verifies that a full freeze → mutate
// original → thaw cycle preserves the frozen state correctly.
func TestIntegration_FreezeThawRoundtrip(t *testing.T) {
	f := freezer.New()

	original := map[string]string{
		"DATABASE_URL": "postgres://localhost/dev",
		"APP_ENV":      "staging",
		"PORT":         "5432",
	}

	if err := f.Freeze("baseline", original); err != nil {
		t.Fatalf("Freeze failed: %v", err)
	}

	// Mutate original after freezing — frozen frame must be unaffected.
	original["DATABASE_URL"] = "postgres://remote/prod"
	original["NEW_KEY"] = "injected"

	thawed, err := f.Thaw("baseline")
	if err != nil {
		t.Fatalf("Thaw failed: %v", err)
	}

	if thawed["DATABASE_URL"] != "postgres://localhost/dev" {
		t.Errorf("DATABASE_URL corrupted: got %q", thawed["DATABASE_URL"])
	}
	if _, ok := thawed["NEW_KEY"]; ok {
		t.Error("NEW_KEY should not exist in frozen frame")
	}
	if thawed["APP_ENV"] != "staging" {
		t.Errorf("APP_ENV corrupted: got %q", thawed["APP_ENV"])
	}

	// Drop and confirm it is gone.
	if err := f.Drop("baseline"); err != nil {
		t.Fatalf("Drop failed: %v", err)
	}
	if _, err := f.Thaw("baseline"); err == nil {
		t.Error("expected error after Drop, got nil")
	}
}
