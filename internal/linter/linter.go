package linter

import (
	"errors"
	"fmt"
	"strings"
)

// Rule represents a single lint rule applied to an environment map.
type Rule struct {
	Name    string
	Message string
	Check   func(key, value string) bool
}

// Issue describes a lint violation found in an environment map.
type Issue struct {
	Key     string
	Rule    string
	Message string
}

func (i Issue) Error() string {
	return fmt.Sprintf("[%s] %s: %s", i.Rule, i.Key, i.Message)
}

// Linter applies a set of rules to an environment map and collects issues.
type Linter struct {
	rules []Rule
}

// New creates a Linter with the provided rules. Returns an error if any rule
// has a blank name or a nil Check function.
func New(rules []Rule) (*Linter, error) {
	for _, r := range rules {
		if strings.TrimSpace(r.Name) == "" {
			return nil, errors.New("linter: rule name must not be blank")
		}
		if r.Check == nil {
			return nil, fmt.Errorf("linter: rule %q has nil Check function", r.Name)
		}
	}
	return &Linter{rules: rules}, nil
}

// Lint runs all rules against the provided environment map and returns any
// issues found. A nil or empty env is valid and returns no issues.
func (l *Linter) Lint(env map[string]string) ([]Issue, error) {
	if env == nil {
		return nil, errors.New("linter: env must not be nil")
	}
	var issues []Issue
	for key, value := range env {
		for _, rule := range l.rules {
			if rule.Check(key, value) {
				issues = append(issues, Issue{
					Key:     key,
					Rule:    rule.Name,
					Message: rule.Message,
				})
			}
		}
	}
	return issues, nil
}

// DefaultRules returns a baseline set of commonly useful lint rules.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name:    "empty-value",
			Message: "variable is set but has an empty value",
			Check:   func(_, value string) bool { return value == "" },
		},
		{
			Name:    "lowercase-key",
			Message: "variable key contains lowercase letters; prefer ALL_CAPS",
			Check:   func(key, _ string) bool { return key != strings.ToUpper(key) },
		},
		{
			Name:    "whitespace-value",
			Message: "variable value has leading or trailing whitespace",
			Check:   func(_, value string) bool { return value != strings.TrimSpace(value) },
		},
	}
}
