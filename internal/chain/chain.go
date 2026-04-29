// Package chain provides functionality for chaining multiple environment
// variable sets together and resolving them in a defined order.
package chain

import (
	"errors"
	"fmt"

	"github.com/envchain/envchain/internal/context"
	"github.com/envchain/envchain/internal/envset"
)

// ErrEmptyChain is returned when attempting to resolve an empty chain.
var ErrEmptyChain = errors.New("chain: no contexts registered in chain")

// Chain holds an ordered sequence of context names to resolve.
type Chain struct {
	names   []string
	ctxReg  *context.Registry
}

// New creates a new Chain using the provided context registry.
func New(reg *context.Registry) *Chain {
	return &Chain{
		names:  []string{},
		ctxReg: reg,
	}
}

// Add appends a context name to the chain.
func (c *Chain) Add(name string) error {
	if name == "" {
		return errors.New("chain: context name must not be empty")
	}
	c.names = append(c.names, name)
	return nil
}

// Resolve walks the chain in order, merging environment sets from each
// context. Later contexts override earlier ones for duplicate keys.
func (c *Chain) Resolve() (*envset.EnvSet, error) {
	if len(c.names) == 0 {
		return nil, ErrEmptyChain
	}

	merged := envset.New()

	for _, name := range c.names {
		ctx, err := c.ctxReg.Get(name)
		if err != nil {
			return nil, fmt.Errorf("chain: failed to get context %q: %w", name, err)
		}

		set, err := ctx.EnvSet()
		if err != nil {
			return nil, fmt.Errorf("chain: failed to load env set for context %q: %w", name, err)
		}

		if err := merged.Merge(set); err != nil {
			return nil, fmt.Errorf("chain: merge failed for context %q: %w", name, err)
		}
	}

	return merged, nil
}

// Names returns the ordered list of context names in the chain.
func (c *Chain) Names() []string {
	out := make([]string, len(c.names))
	copy(out, c.names)
	return out
}
