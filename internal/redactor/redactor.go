// Package redactor provides utilities for masking sensitive environment
// variable values before they are logged, exported, or displayed to users.
package redactor

import (
	"errors"
	"strings"
)

const defaultMask = "********"

// Redactor masks the values of sensitive keys in an environment map.
type Redactor struct {
	sensitiveKeys map[string]struct{}
	mask          string
}

// New creates a new Redactor with the given set of sensitive key names.
// Key matching is case-insensitive. An empty sensitiveKeys slice is valid.
func New(sensitiveKeys []string, mask string) (*Redactor, error) {
	if mask == "" {
		mask = defaultMask
	}
	keySet := make(map[string]struct{}, len(sensitiveKeys))
	for _, k := range sensitiveKeys {
		if k == "" {
			return nil, errors.New("redactor: sensitive key name must not be empty")
		}
		keySet[strings.ToUpper(k)] = struct{}{}
	}
	return &Redactor{
		sensitiveKeys: keySet,
		mask:          mask,
	}, nil
}

// Redact returns a copy of the provided map with sensitive values replaced
// by the configured mask string. The original map is not modified.
func (r *Redactor) Redact(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if _, sensitive := r.sensitiveKeys[strings.ToUpper(k)]; sensitive {
			out[k] = r.mask
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether the given key is considered sensitive.
// The check is case-insensitive.
func (r *Redactor) IsSensitive(key string) bool {
	_, ok := r.sensitiveKeys[strings.ToUpper(key)]
	return ok
}

// Mask returns the mask string used by this Redactor.
func (r *Redactor) Mask() string {
	return r.mask
}
