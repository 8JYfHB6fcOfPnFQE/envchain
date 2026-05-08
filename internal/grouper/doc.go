// Package grouper provides utilities for organising environment variable keys
// into named groups.
//
// A Grouper maintains a registry of named groups. Keys can be added to a group
// explicitly via Add, or discovered automatically from a live environment map
// using GroupByPrefix, which matches keys whose names start with
// "<GROUP_NAME>_" (case-insensitive).
//
// Keys that do not match any registered group prefix are collected under the
// reserved "__ungrouped__" group name, making it straightforward to identify
// variables that have not yet been classified.
//
// Example usage:
//
//	g := grouper.New()
//	g.Add("db", "DB_HOST")
//	g.Add("db", "DB_PORT")
//	result := g.GroupByPrefix(os.Environ())
package grouper
