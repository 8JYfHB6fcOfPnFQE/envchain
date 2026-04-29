// Package merger provides functionality for merging multiple environment
// variable sets into a single resolved map, with support for precedence ordering.
package merger

import "fmt"

// Merger combines multiple environment variable maps into one,
// applying precedence rules where later sources override earlier ones.
type Merger struct {
	sources []map[string]string
	names   []string
}

// New creates a new Merger instance.
func New() *Merger {
	return &Merger{
		sources: make([]map[string]string, 0),
		names:   make([]string, 0),
	}
}

// Add appends a named environment variable map as a source.
// Sources added later take precedence over earlier ones.
func (m *Merger) Add(name string, env map[string]string) error {
	if name == "" {
		return fmt.Errorf("merger: source name must not be empty")
	}
	if env == nil {
		return fmt.Errorf("merger: env map for source %q must not be nil", name)
	}
	m.names = append(m.names, name)
	m.sources = append(m.sources, env)
	return nil
}

// Merge combines all added sources into a single map.
// Keys from later sources overwrite keys from earlier sources.
// Returns an error if no sources have been added.
func (m *Merger) Merge() (map[string]string, error) {
	if len(m.sources) == 0 {
		return nil, fmt.Errorf("merger: no sources added")
	}

	result := make(map[string]string)
	for _, src := range m.sources {
		for k, v := range src {
			result[k] = v
		}
	}
	return result, nil
}

// Sources returns the names of all registered sources in order.
func (m *Merger) Sources() []string {
	out := make([]string, len(m.names))
	copy(out, m.names)
	return out
}
