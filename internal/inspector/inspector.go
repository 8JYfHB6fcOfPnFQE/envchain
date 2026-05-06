// Package inspector provides utilities for inspecting and summarising
// the contents of an environment variable map, including key counts,
// classification breakdowns, and value presence statistics.
package inspector

import (
	"errors"
	"sort"
)

// Summary holds aggregate statistics about an environment variable map.
type Summary struct {
	TotalKeys    int
	PresentKeys  int
	MissingKeys  int
	KeyNames     []string
}

// Inspector analyses an environment variable map and produces summaries.
type Inspector struct {
	required []string
}

// New creates an Inspector that treats the given keys as required when
// computing missing/present counts. Keys must not be empty.
func New(required []string) (*Inspector, error) {
	for _, k := range required {
		if k == "" {
			return nil, errors.New("inspector: required key must not be empty")
		}
	}
	copy := make([]string, len(required))
	copy = append(copy[:0], required...)
	return &Inspector{required: copy}, nil
}

// Inspect analyses env and returns a Summary. env must not be nil.
func (i *Inspector) Inspect(env map[string]string) (*Summary, error) {
	if env == nil {
		return nil, errors.New("inspector: env map must not be nil")
	}

	keys := make([]string, 0, len(env))
	for k := range env {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	present := 0
	missing := 0
	for _, req := range i.required {
		if v, ok := env[req]; ok && v != "" {
			present++
		} else {
			missing++
		}
	}

	return &Summary{
		TotalKeys:   len(env),
		PresentKeys: present,
		MissingKeys: missing,
		KeyNames:    keys,
	}, nil
}

// MissingRequired returns the subset of required keys absent or empty in env.
func (i *Inspector) MissingRequired(env map[string]string) ([]string, error) {
	if env == nil {
		return nil, errors.New("inspector: env map must not be nil")
	}
	var missing []string
	for _, req := range i.required {
		if v, ok := env[req]; !ok || v == "" {
			missing = append(missing, req)
		}
	}
	sort.Strings(missing)
	return missing, nil
}
