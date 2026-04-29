package validator_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/validator"
)

func TestValidate_RequiredMissing(t *testing.T) {
	v := validator.New(map[string]validator.Rule{
		"DATABASE_URL": {Required: true},
	})

	errs := v.Validate(map[string]string{})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
}

func TestValidate_RequiredPresent(t *testing.T) {
	v := validator.New(map[string]validator.Rule{
		"DATABASE_URL": {Required: true},
	})

	errs := v.Validate(map[string]string{"DATABASE_URL": "postgres://localhost/db"})
	if errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	v := validator.New(map[string]validator.Rule{
		"PORT": {Pattern: `^\d+$`},
	})

	if errs := v.Validate(map[string]string{"PORT": "8080"}); errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}

	if errs := v.Validate(map[string]string{"PORT": "abc"}); len(errs) != 1 {
		t.Fatalf("expected 1 error for invalid pattern, got %d", len(errs))
	}
}

func TestValidate_AllowedValues(t *testing.T) {
	v := validator.New(map[string]validator.Rule{
		"ENV": {AllowedValues: []string{"development", "staging", "production"}},
	})

	if errs := v.Validate(map[string]string{"ENV": "staging"}); errs != nil {
		t.Fatalf("expected no errors, got %v", errs)
	}

	if errs := v.Validate(map[string]string{"ENV": "local"}); len(errs) != 1 {
		t.Fatalf("expected 1 error for disallowed value, got %d", len(errs))
	}
}

func TestValidate_OptionalMissingSkipped(t *testing.T) {
	v := validator.New(map[string]validator.Rule{
		"LOG_LEVEL": {Pattern: `^(debug|info|warn|error)$`},
	})

	// Optional variable not present — should not trigger pattern error
	if errs := v.Validate(map[string]string{}); errs != nil {
		t.Fatalf("expected no errors for missing optional var, got %v", errs)
	}
}
