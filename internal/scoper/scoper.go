// Package scoper provides scope-based filtering of environment variable maps,
// allowing callers to restrict or project a set of key-value pairs to only
// those keys that match a declared scope prefix or explicit inclusion list.
package scoper

import (
	"errors"
	"strings"
)

// Scoper filters an environment map to a named scope.
type Scoper struct {
	prefix  string
	strip   bool
	explicit map[string]struct{}
}

// New creates a Scoper that matches keys sharing the given prefix.
// If strip is true the prefix is removed from keys in the result.
// prefix must not be empty.
func New(prefix string, strip bool) (*Scoper, error) {
	if strings.TrimSpace(prefix) == "" {
		return nil, errors.New("scoper: prefix must not be blank")
	}
	return &Scoper{
		prefix:   strings.ToUpper(prefix),
		strip:    strip,
		explicit: make(map[string]struct{}),
	}, nil
}

// Include adds an explicit key (case-insensitive) that is always included
// in the scoped output regardless of the prefix rule.
func (s *Scoper) Include(key string) error {
	if strings.TrimSpace(key) == "" {
		return errors.New("scoper: key must not be blank")
	}
	s.explicit[strings.ToUpper(key)] = struct{}{}
	return nil
}

// Scope returns a new map containing only the keys from env that belong to
// the scope. Keys are matched case-insensitively against the prefix.
func (s *Scoper) Scope(env map[string]string) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("scoper: env must not be nil")
	}
	out := make(map[string]string)
	for k, v := range env {
		upper := strings.ToUpper(k)
		_, explicit := s.explicit[upper]
		if explicit {
			out[k] = v
			continue
		}
		if strings.HasPrefix(upper, s.prefix) {
			outKey := k
			if s.strip {
				outKey = k[len(s.prefix):]
				if outKey == "" {
					continue
				}
			}
			out[outKey] = v
		}
	}
	return out, nil
}
