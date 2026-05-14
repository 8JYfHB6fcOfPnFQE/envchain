// Package interpolator provides variable interpolation within environment
// variable values, expanding references such as ${VAR} or $VAR using a
// provided environment map as the source of substitutions.
package interpolator

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Interpolator expands variable references within environment values.
type Interpolator struct {
	fallbackToOS bool
	maxDepth     int
}

// New creates an Interpolator. When fallbackToOS is true, unresolved
// references are looked up in the process environment before being left
// unexpanded. maxDepth limits recursive expansion to prevent cycles.
func New(fallbackToOS bool, maxDepth int) (*Interpolator, error) {
	if maxDepth <= 0 {
		return nil, errors.New("interpolator: maxDepth must be greater than zero")
	}
	return &Interpolator{fallbackToOS: fallbackToOS, maxDepth: maxDepth}, nil
}

// Expand returns a new map where every value has its ${VAR} and $VAR
// references replaced with values from env, optionally falling back to
// os.Getenv. Returns an error if env is nil.
func (i *Interpolator) Expand(env map[string]string) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("interpolator: env must not be nil")
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		expanded, err := i.expandValue(v, env, 0)
		if err != nil {
			return nil, fmt.Errorf("interpolator: key %q: %w", k, err)
		}
		out[k] = expanded
	}
	return out, nil
}

// ExpandValue expands a single value string using env as the variable source.
func (i *Interpolator) ExpandValue(value string, env map[string]string) (string, error) {
	if env == nil {
		return "", errors.New("interpolator: env must not be nil")
	}
	return i.expandValue(value, env, 0)
}

func (i *Interpolator) expandValue(value string, env map[string]string, depth int) (string, error) {
	if depth > i.maxDepth {
		return "", errors.New("max expansion depth exceeded; possible cycle")
	}
	result := os.Expand(value, func(key string) string {
		if v, ok := env[key]; ok {
			expanded, _ := i.expandValue(v, env, depth+1)
			return expanded
		}
		if i.fallbackToOS {
			return os.Getenv(key)
		}
		return "$" + key
	})
	_ = strings.TrimSpace // keep import used via os.Expand
	return result, nil
}
