// Package profiler provides tools for defining and matching named profiles of
// expected environment variable keys against a given environment map.
//
// A Profile declares the set of keys that a deployment context is expected to
// supply. The Profiler registry allows multiple profiles to be defined and
// retrieved by name.
//
// The Match function compares a profile's declared keys against a live
// environment map, reporting which keys are present (matched), absent
// (missing), and undeclared (extra). Key comparisons are case-insensitive.
//
// Typical usage:
//
//	p := profiler.New()
//	p.Define("production", []string{"DB_URL", "API_KEY"}, map[string]string{"tier": "high"})
//	pr, _ := p.Get("production")
//	result, _ := profiler.Match(pr, os.Environ()...)
//	if !result.IsComplete() {
//		log.Printf("incomplete profile: %s", result.Summary())
//	}
package profiler
