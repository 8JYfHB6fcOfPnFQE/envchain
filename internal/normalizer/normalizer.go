package normalizer

import (
	"errors"
	"strings"
)

// Normalizer applies a set of normalization rules to environment variable maps.
// Rules are applied in the order they were registered.
type Normalizer struct {
	rules []Rule
}

// Rule defines a named transformation applied to every key-value pair.
type Rule struct {
	Name  string
	Apply func(key, value string) (string, string)
}

// New creates a Normalizer with the provided rules.
// Returns an error if any rule has a blank name or a nil Apply function.
func New(rules []Rule) (*Normalizer, error) {
	for _, r := range rules {
		if strings.TrimSpace(r.Name) == "" {
			return nil, errors.New("normalizer: rule name must not be blank")
		}
		if r.Apply == nil {
			return nil, errors.New("normalizer: rule Apply must not be nil")
		}
	}
	return &Normalizer{rules: rules}, nil
}

// Normalize applies all rules to the given environment map and returns a new
// map with the transformed key-value pairs. The original map is not mutated.
func (n *Normalizer) Normalize(env map[string]string) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("normalizer: env must not be nil")
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}
	for _, rule := range n.rules {
		updated := make(map[string]string, len(out))
		for k, v := range out {
			nk, nv := rule.Apply(k, v)
			updated[nk] = nv
		}
		out = updated
	}
	return out, nil
}

// DefaultRules returns a standard set of normalization rules:
//   - TrimSpace: trims leading/trailing whitespace from both key and value.
//   - UpperKey: converts all keys to upper case.
func DefaultRules() []Rule {
	return []Rule{
		{
			Name: "TrimSpace",
			Apply: func(k, v string) (string, string) {
				return strings.TrimSpace(k), strings.TrimSpace(v)
			},
		},
		{
			Name: "UpperKey",
			Apply: func(k, v string) (string, string) {
				return strings.ToUpper(k), v
			},
		},
	}
}
