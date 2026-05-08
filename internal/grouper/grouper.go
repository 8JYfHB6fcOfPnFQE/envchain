package grouper

import (
	"errors"
	"sort"
	"strings"
)

// Group holds a named collection of environment variable keys.
type Group struct {
	Name string
	Keys []string
}

// Grouper organises environment variable keys into named groups by prefix or
// explicit assignment, enabling bulk operations on related variables.
type Grouper struct {
	groups map[string][]string
}

// New returns an initialised Grouper.
func New() *Grouper {
	return &Grouper{groups: make(map[string][]string)}
}

// Add registers key under the named group. Both name and key must be non-blank.
func (g *Grouper) Add(name, key string) error {
	name = strings.TrimSpace(name)
	key = strings.TrimSpace(key)
	if name == "" {
		return errors.New("grouper: group name must not be blank")
	}
	if key == "" {
		return errors.New("grouper: key must not be blank")
	}
	g.groups[name] = append(g.groups[name], key)
	return nil
}

// GroupByPrefix scans env and places every key whose prefix matches a
// registered group name (case-insensitive, separated by "_") into that group.
// Keys that match no group are placed into the "__ungrouped__" group.
func (g *Grouper) GroupByPrefix(env map[string]string) map[string][]string {
	result := make(map[string][]string)
	for k := range env {
		matched := false
		upper := strings.ToUpper(k)
		for name := range g.groups {
			prefix := strings.ToUpper(name) + "_"
			if strings.HasPrefix(upper, prefix) {
				result[name] = append(result[name], k)
				matched = true
				break
			}
		}
		if !matched {
			result["__ungrouped__"] = append(result["__ungrouped__"], k)
		}
	}
	for name := range result {
		sort.Strings(result[name])
	}
	return result
}

// List returns all registered groups in sorted order.
func (g *Grouper) List() []Group {
	names := make([]string, 0, len(g.groups))
	for n := range g.groups {
		names = append(names, n)
	}
	sort.Strings(names)
	out := make([]Group, 0, len(names))
	for _, n := range names {
		keys := make([]string, len(g.groups[n]))
		copy(keys, g.groups[n])
		sort.Strings(keys)
		out = append(out, Group{Name: n, Keys: keys})
	}
	return out
}
