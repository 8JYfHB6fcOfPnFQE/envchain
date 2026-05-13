// Package flattener provides utilities for collapsing nested environment
// variable maps into a single flat key-value map, applying optional prefix
// namespacing and collision strategies.
package flattener

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// CollisionStrategy determines how key collisions are resolved when merging
// multiple sources into a flat map.
type CollisionStrategy int

const (
	// CollisionKeepFirst retains the first value seen for a colliding key.
	CollisionKeepFirst CollisionStrategy = iota
	// CollisionKeepLast overwrites with the most recently seen value.
	CollisionKeepLast
	// CollisionError returns an error when a collision is detected.
	CollisionError
)

// Flattener merges one or more named env maps into a single flat map.
type Flattener struct {
	strategy CollisionStrategy
	sources  []source
}

type source struct {
	prefix string
	env    map[string]string
}

// New creates a Flattener with the given collision strategy.
func New(strategy CollisionStrategy) (*Flattener, error) {
	if strategy < CollisionKeepFirst || strategy > CollisionError {
		return nil, errors.New("flattener: unknown collision strategy")
	}
	return &Flattener{strategy: strategy}, nil
}

// Add registers a named env map with an optional prefix.
// If prefix is non-empty, all keys from env will be namespaced as PREFIX_KEY.
func (f *Flattener) Add(prefix string, env map[string]string) error {
	if env == nil {
		return errors.New("flattener: env map must not be nil")
	}
	f.sources = append(f.sources, source{prefix: strings.ToUpper(prefix), env: env})
	return nil
}

// Flatten merges all registered sources into a single flat map.
// Keys are uppercased; prefix (if any) is prepended with an underscore separator.
func (f *Flattener) Flatten() (map[string]string, error) {
	if len(f.sources) == 0 {
		return nil, errors.New("flattener: no sources registered")
	}

	result := make(map[string]string)

	for _, s := range f.sources {
		keys := make([]string, 0, len(s.env))
		for k := range s.env {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			flatKey := strings.ToUpper(k)
			if s.prefix != "" {
				flatKey = s.prefix + "_" + flatKey
			}
			v := s.env[k]

			if existing, ok := result[flatKey]; ok {
				switch f.strategy {
				case CollisionKeepFirst:
					// retain existing — do nothing
				case CollisionKeepLast:
					result[flatKey] = v
				case CollisionError:
					return nil, fmt.Errorf("flattener: key collision on %q (existing=%q, new=%q)", flatKey, existing, v)
				}
			} else {
				result[flatKey] = v
			}
		}
	}

	return result, nil
}
