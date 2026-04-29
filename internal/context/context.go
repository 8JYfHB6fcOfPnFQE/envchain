// Package context manages named deployment contexts for envchain.
// A context associates a name (e.g. "production", "staging") with
// a specific set of environment variable definitions and validation rules.
package context

import (
	"fmt"
	"sync"
)

// Context represents a named deployment context containing
// one or more named environment variable sets.
type Context struct {
	Name    string
	EnvSets map[string]EnvSetRef
}

// EnvSetRef holds a reference to an env set by name and its source file path.
type EnvSetRef struct {
	SetName  string
	FilePath string
}

// Registry holds all registered deployment contexts.
type Registry struct {
	mu       sync.RWMutex
	contexts map[string]*Context
}

// New creates a new empty Registry.
func New() *Registry {
	return &Registry{
		contexts: make(map[string]*Context),
	}
}

// Register adds or replaces a context in the registry.
func (r *Registry) Register(ctx *Context) error {
	if ctx == nil {
		return fmt.Errorf("context must not be nil")
	}
	if ctx.Name == "" {
		return fmt.Errorf("context name must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.contexts[ctx.Name] = ctx
	return nil
}

// Get retrieves a context by name. Returns an error if not found.
func (r *Registry) Get(name string) (*Context, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	ctx, ok := r.contexts[name]
	if !ok {
		return nil, fmt.Errorf("context %q not found", name)
	}
	return ctx, nil
}

// List returns all registered context names in no guaranteed order.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.contexts))
	for name := range r.contexts {
		names = append(names, name)
	}
	return names
}

// Remove deletes a context from the registry by name.
func (r *Registry) Remove(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.contexts, name)
}
