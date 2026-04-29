// Package resolver provides functionality for resolving environment variable
// values across chained deployment contexts, applying merging and validation.
package resolver

import (
	"errors"
	"fmt"

	"github.com/user/envchain/internal/chain"
	"github.com/user/envchain/internal/merger"
	"github.com/user/envchain/internal/validator"
)

// Resolver combines a chain, merger, and validator to produce a final
// resolved map of environment variables for a given context chain.
type Resolver struct {
	chain     *chain.Chain
	merger    *merger.Merger
	validator *validator.Validator
}

// New creates a new Resolver with the provided chain, merger, and validator.
// Returns an error if any argument is nil.
func New(c *chain.Chain, m *merger.Merger, v *validator.Validator) (*Resolver, error) {
	if c == nil {
		return nil, errors.New("resolver: chain must not be nil")
	}
	if m == nil {
		return nil, errors.New("resolver: merger must not be nil")
	}
	if v == nil {
		return nil, errors.New("resolver: validator must not be nil")
	}
	return &Resolver{chain: c, merger: m, validator: v}, nil
}

// Resolve walks the chain, merges all resolved context maps in order,
// validates the final merged result, and returns it.
func (r *Resolver) Resolve() (map[string]string, error) {
	names, err := r.chain.Resolve()
	if err != nil {
		return nil, fmt.Errorf("resolver: chain resolution failed: %w", err)
	}

	for _, name := range names {
		env, err := r.chain.Registry().Get(name)
		if err != nil {
			return nil, fmt.Errorf("resolver: context %q not found: %w", name, err)
		}
		if err := r.merger.Add(name, env); err != nil {
			return nil, fmt.Errorf("resolver: failed to add context %q to merger: %w", name, err)
		}
	}

	merged, err := r.merger.Merge()
	if err != nil {
		return nil, fmt.Errorf("resolver: merge failed: %w", err)
	}

	if err := r.validator.Validate(merged); err != nil {
		return nil, fmt.Errorf("resolver: validation failed: %w", err)
	}

	return merged, nil
}
