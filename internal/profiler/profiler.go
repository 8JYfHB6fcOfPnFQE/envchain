package profiler

import (
	"errors"
	"sort"
	"strings"
)

// Profile represents a named set of environment variable keys expected in a
// deployment context, along with optional metadata tags.
type Profile struct {
	name string
	keys []string
	tags map[string]string
}

// Profiler manages a collection of named profiles.
type Profiler struct {
	profiles map[string]*Profile
}

// New returns a new Profiler instance.
func New() *Profiler {
	return &Profiler{profiles: make(map[string]*Profile)}
}

// Define registers a new profile with the given name, keys, and tags.
// Returns an error if the name is empty, keys list is empty, or the profile
// already exists.
func (p *Profiler) Define(name string, keys []string, tags map[string]string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("profiler: name must not be empty")
	}
	if len(keys) == 0 {
		return errors.New("profiler: keys must not be empty")
	}
	if _, exists := p.profiles[name]; exists {
		return errors.New("profiler: profile already defined: " + name)
	}
	copiedKeys := make([]string, len(keys))
	copy(copiedKeys, keys)
	copiedTags := make(map[string]string, len(tags))
	for k, v := range tags {
		copiedTags[k] = v
	}
	p.profiles[name] = &Profile{name: name, keys: copiedKeys, tags: copiedTags}
	return nil
}

// Get returns the Profile registered under name, or an error if not found.
func (p *Profiler) Get(name string) (*Profile, error) {
	pr, ok := p.profiles[name]
	if !ok {
		return nil, errors.New("profiler: profile not found: " + name)
	}
	return pr, nil
}

// List returns all registered profile names in sorted order.
func (p *Profiler) List() []string {
	names := make([]string, 0, len(p.profiles))
	for n := range p.profiles {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

// Keys returns the keys associated with the profile.
func (pr *Profile) Keys() []string {
	out := make([]string, len(pr.keys))
	copy(out, pr.keys)
	return out
}

// Tags returns the tags associated with the profile.
func (pr *Profile) Tags() map[string]string {
	out := make(map[string]string, len(pr.tags))
	for k, v := range pr.tags {
		out[k] = v
	}
	return out
}

// Name returns the profile name.
func (pr *Profile) Name() string { return pr.name }
