// Package classifier categorises environment variable keys by sensitivity level
// based on configurable keyword patterns.
package classifier

import (
	"errors"
	"strings"
)

// Level represents the sensitivity classification of an environment variable.
type Level int

const (
	LevelPublic Level = iota
	LevelInternal
	LevelSecret
)

// String returns a human-readable label for the Level.
func (l Level) String() string {
	switch l {
	case LevelSecret:
		return "secret"
	case LevelInternal:
		return "internal"
	default:
		return "public"
	}
}

// Classifier assigns sensitivity levels to environment variable keys.
type Classifier struct {
	secretPatterns   []string
	internalPatterns []string
}

// New creates a Classifier with default keyword patterns.
// Additional patterns may be supplied; they are merged with the defaults.
func New(extraSecret, extraInternal []string) (*Classifier, error) {
	defaultSecret := []string{"secret", "password", "passwd", "token", "apikey", "api_key", "private_key", "credentials"}
	defaultInternal := []string{"internal", "host", "port", "endpoint", "url", "dsn", "addr"}

	for _, p := range extraSecret {
		if strings.TrimSpace(p) == "" {
			return nil, errors.New("classifier: secret pattern must not be blank")
		}
	}
	for _, p := range extraInternal {
		if strings.TrimSpace(p) == "" {
			return nil, errors.New("classifier: internal pattern must not be blank")
		}
	}

	return &Classifier{
		secretPatterns:   append(defaultSecret, extraSecret...),
		internalPatterns: append(defaultInternal, extraInternal...),
	}, nil
}

// Classify returns the sensitivity Level for the given key.
// Matching is case-insensitive and checks whether any pattern is a substring
// of the normalised key.
func (c *Classifier) Classify(key string) Level {
	norm := strings.ToLower(key)
	for _, p := range c.secretPatterns {
		if strings.Contains(norm, strings.ToLower(p)) {
			return LevelSecret
		}
	}
	for _, p := range c.internalPatterns {
		if strings.Contains(norm, strings.ToLower(p)) {
			return LevelInternal
		}
	}
	return LevelPublic
}

// ClassifyMap returns a map of key → Level for every key in the provided set.
func (c *Classifier) ClassifyMap(env map[string]string) map[string]Level {
	out := make(map[string]Level, len(env))
	for k := range env {
		out[k] = c.Classify(k)
	}
	return out
}
