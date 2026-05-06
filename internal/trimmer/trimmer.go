// Package trimmer provides utilities for normalising environment variable
// maps by stripping keys or values that do not meet basic hygiene criteria.
package trimmer

import (
	"errors"
	"strings"
)

// Trimmer removes empty keys, blank values, and optionally keys matching a
// set of exact names or prefixes from an environment variable map.
type Trimmer struct {
	excludeKeys    map[string]struct{}
	excludePrefixes []string
	dropBlankValues bool
}

// New constructs a Trimmer. excludeKeys is a set of exact key names (compared
// case-insensitively) to drop. excludePrefixes is a list of key prefixes to
// drop. dropBlankValues controls whether keys whose value is empty or
// whitespace-only are removed.
func New(excludeKeys []string, excludePrefixes []string, dropBlankValues bool) (*Trimmer, error) {
	eks := make(map[string]struct{}, len(excludeKeys))
	for _, k := range excludeKeys {
		if k == "" {
			return nil, errors.New("trimmer: exclude key must not be empty")
		}
		eks[strings.ToUpper(k)] = struct{}{}
	}
	for _, p := range excludePrefixes {
		if p == "" {
			return nil, errors.New("trimmer: exclude prefix must not be empty")
		}
	}
	return &Trimmer{
		excludeKeys:     eks,
		excludePrefixes: excludePrefixes,
		dropBlankValues: dropBlankValues,
	}, nil
}

// Trim returns a new map with unwanted entries removed. The original map is
// never mutated.
func (t *Trimmer) Trim(env map[string]string) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("trimmer: env map must not be nil")
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		if k == "" {
			continue
		}
		if t.dropBlankValues && strings.TrimSpace(v) == "" {
			continue
		}
		upper := strings.ToUpper(k)
		if _, excluded := t.excludeKeys[upper]; excluded {
			continue
		}
		if t.hasExcludedPrefix(upper) {
			continue
		}
		out[k] = v
	}
	return out, nil
}

func (t *Trimmer) hasExcludedPrefix(upperKey string) bool {
	for _, p := range t.excludePrefixes {
		if strings.HasPrefix(upperKey, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}
