// Package labeler provides a mechanism for attaching and querying
// named string labels on environment variable keys.
//
// Labels are arbitrary key-value annotations (e.g. env=production,
// tier=backend, sensitivity=high) that can be used by other pipeline
// components to filter, route, or audit environment variables without
// modifying their values.
//
// Usage:
//
//	l := labeler.New()
//	l.Attach("DB_PASSWORD", "sensitivity", "high")
//	l.Attach("DB_PASSWORD", "env", "production")
//
//	lbls := l.Get("DB_PASSWORD")
//	keys := l.FindByLabel("sensitivity", "high")
//
// Key matching is case-insensitive; keys are normalised to uppercase
// internally.
package labeler
