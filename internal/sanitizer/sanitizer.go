// Package sanitizer provides utilities for sanitizing environment variable
// keys and values before they are exported or stored. It trims whitespace,
// normalizes key casing, and removes unsafe characters from values.
package sanitizer

import (
	"errors"
	"regexp"
	"strings"
)

var invalidKeyChars = regexp.MustCompile(`[^A-Z0-9_]`)

// Sanitizer cleans environment variable maps.
type Sanitizer struct {
	stripControlChars bool
}

// New creates a new Sanitizer. If stripControlChars is true, control
// characters (e.g. \t, \r, \n) are removed from values.
func New(stripControlChars bool) *Sanitizer {
	return &Sanitizer{stripControlChars: stripControlChars}
}

// SanitizeKey normalizes an environment variable key to uppercase and
// replaces any character that is not A-Z, 0-9, or underscore with an
// underscore. Returns an error if the resulting key is empty.
func (s *Sanitizer) SanitizeKey(key string) (string, error) {
	upper := strings.ToUpper(strings.TrimSpace(key))
	cleaned := invalidKeyChars.ReplaceAllString(upper, "_")
	if cleaned == "" {
		return "", errors.New("sanitizer: key is empty after sanitization")
	}
	return cleaned, nil
}

// SanitizeValue trims leading and trailing whitespace from a value and,
// if configured, strips control characters.
func (s *Sanitizer) SanitizeValue(value string) string {
	v := strings.TrimSpace(value)
	if s.stripControlChars {
		v = strings.Map(func(r rune) rune {
			if r == '\t' || r == '\r' || r == '\n' {
				return -1
			}
			return r
		}, v)
	}
	return v
}

// SanitizeMap applies SanitizeKey and SanitizeValue to every entry in the
// provided map, returning a new sanitized map. If a key cannot be sanitized
// it is skipped and the error is collected. The first encountered error is
// returned alongside any successfully sanitized entries.
func (s *Sanitizer) SanitizeMap(env map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(env))
	var firstErr error
	for k, v := range env {
		cleanKey, err := s.SanitizeKey(k)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			continue
		}
		out[cleanKey] = s.SanitizeValue(v)
	}
	return out, firstErr
}
