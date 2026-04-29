// Package validator provides validation rules for environment variable sets.
package validator

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule represents a validation rule for an environment variable.
type Rule struct {
	Required bool
	Pattern  string
	AllowedValues []string
}

// Validator holds a map of variable names to their validation rules.
type Validator struct {
	rules map[string]Rule
}

// New creates a new Validator with the given rules.
func New(rules map[string]Rule) *Validator {
	return &Validator{rules: rules}
}

// Validate checks the provided env map against all registered rules.
// It returns a list of validation errors, or nil if all checks pass.
func (v *Validator) Validate(env map[string]string) []error {
	var errs []error

	for name, rule := range v.rules {
		val, exists := env[name]

		if rule.Required && !exists {
			errs = append(errs, fmt.Errorf("required variable %q is missing", name))
			continue
		}

		if !exists {
			continue
		}

		if rule.Pattern != "" {
			matched, err := regexp.MatchString(rule.Pattern, val)
			if err != nil {
				errs = append(errs, fmt.Errorf("invalid pattern for %q: %w", name, err))
			} else if !matched {
				errs = append(errs, fmt.Errorf("variable %q value %q does not match pattern %q", name, val, rule.Pattern))
			}
		}

		if len(rule.AllowedValues) > 0 {
			if !contains(rule.AllowedValues, val) {
				errs = append(errs, fmt.Errorf("variable %q value %q is not in allowed values [%s]", name, val, strings.Join(rule.AllowedValues, ", ")))
			}
		}
	}

	if len(errs) == 0 {
		return nil
	}
	return errs
}

func contains(list []string, val string) bool {
	for _, v := range list {
		if v == val {
			return true
		}
	}
	return false
}
