// Package tagger provides functionality for attaching and querying
// named tags on environment variable keys.
//
// Tags are arbitrary key-value labels that can be used to annotate
// environment variables with metadata such as sensitivity level,
// deployment environment, or ownership group.
//
// Example usage:
//
//	tr := tagger.New()
//	_ = tr.Tag("DB_PASSWORD", "sensitivity", "high")
//	_ = tr.Tag("DB_PASSWORD", "env", "production")
//
//	sensitiveKeys := tr.KeysWithTag("sensitivity")
//
// Tags do not affect variable resolution or validation; they serve
// as informational annotations for tooling and reporting.
package tagger
