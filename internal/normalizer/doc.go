// Package normalizer provides a composable pipeline for normalizing environment
// variable maps before they are validated, exported, or compared.
//
// A Normalizer is constructed with an ordered slice of Rules. Each Rule carries
// a human-readable Name and an Apply function that receives a key-value pair and
// returns a (possibly transformed) key-value pair. Rules are applied in
// registration order, and the output of one rule feeds into the next.
//
// DefaultRules returns a standard set suitable for most deployments:
//
//	- TrimSpace  – strips leading/trailing whitespace from keys and values.
//	- UpperKey   – converts all keys to upper case for consistency.
//
// Example usage:
//
//	n, err := normalizer.New(normalizer.DefaultRules())
//	if err != nil { ... }
//	clean, err := n.Normalize(rawEnv)
package normalizer
