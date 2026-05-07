// Package tracer provides event tracing for environment variable resolution
// steps across deployment contexts.
//
// A Tracer records timestamped events that describe what happened to each
// environment variable during chain resolution — such as when a key was set,
// overridden, or skipped. Events can be filtered by context name or summarised
// into a human-readable report.
//
// Basic usage:
//
//	tr := tracer.New()
//	_ = tr.Record("production", "DB_URL", "resolved", "from base context")
//	_ = tr.Record("production", "API_KEY", "overridden", "by overlay context")
//
//	events := tr.FilterByContext("production")
//	fmt.Println(tr.Summary())
//
// Tracer is safe for concurrent use.
package tracer
