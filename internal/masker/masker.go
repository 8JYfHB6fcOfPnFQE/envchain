// Package masker provides utilities for masking sensitive environment
// variable values before they are written to output or logs.
package masker

import (
	"errors"
	"strings"
)

const defaultMask = "********"

// Masker holds a set of key patterns considered sensitive and replaces
// their values with a fixed mask string.
type Masker struct {
	keys []string
	mask string
}

// New creates a Masker that will obscure values for the given keys.
// Keys are matched case-insensitively. mask may be empty to use the
// default mask string.
func New(keys []string, mask string) (*Masker, error) {
	for _, k := range keys {
		if strings.TrimSpace(k) == "" {
			return nil, errors.New("masker: key must not be blank")
		}
	}
	if mask == "" {
		mask = defaultMask
	}
	normalized := make([]string, len(keys))
	for i, k := range keys {
		normalized[i] = strings.ToLower(k)
	}
	return &Masker{keys: normalized, mask: mask}, nil
}

// Apply returns a copy of env where values whose keys are considered
// sensitive have been replaced with the mask string.
func (m *Masker) Apply(env map[string]string) map[string]string {
	out := make(map[string]string, len(env))
	for k, v := range env {
		if m.isSensitive(k) {
			out[k] = m.mask
		} else {
			out[k] = v
		}
	}
	return out
}

// IsSensitive reports whether the given key should be masked.
func (m *Masker) IsSensitive(key string) bool {
	return m.isSensitive(key)
}

func (m *Masker) isSensitive(key string) bool {
	lower := strings.ToLower(key)
	for _, k := range m.keys {
		if k == lower {
			return true
		}
	}
	return false
}

// Mask returns the mask string used by this Masker.
func (m *Masker) Mask() string {
	return m.mask
}
