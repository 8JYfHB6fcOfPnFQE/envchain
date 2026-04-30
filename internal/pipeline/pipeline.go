// Package pipeline orchestrates the full envchain resolution flow,
// combining loading, merging, validation, and exporting into a single
// composable execution unit.
package pipeline

import (
	"fmt"

	"github.com/yourorg/envchain/internal/exporter"
	"github.com/yourorg/envchain/internal/resolver"
)

// Pipeline coordinates the end-to-end processing of an environment chain.
type Pipeline struct {
	resolver *resolver.Resolver
	exporter *exporter.Exporter
}

// Config holds the options used to construct a Pipeline.
type Config struct {
	Resolver *resolver.Resolver
	Exporter *exporter.Exporter
}

// New creates a Pipeline from the provided Config.
// Returns an error if any required component is nil.
func New(cfg Config) (*Pipeline, error) {
	if cfg.Resolver == nil {
		return nil, fmt.Errorf("pipeline: resolver must not be nil")
	}
	if cfg.Exporter == nil {
		return nil, fmt.Errorf("pipeline: exporter must not be nil")
	}
	return &Pipeline{
		resolver: cfg.Resolver,
		exporter: cfg.Exporter,
	}, nil
}

// Run executes the pipeline for the given ordered list of context names.
// It resolves and merges the contexts, then writes the result via the exporter.
func (p *Pipeline) Run(contexts []string) error {
	if len(contexts) == 0 {
		return fmt.Errorf("pipeline: at least one context name is required")
	}

	resolved, err := p.resolver.Resolve(contexts)
	if err != nil {
		return fmt.Errorf("pipeline: resolve failed: %w", err)
	}

	if err := p.exporter.Write(resolved); err != nil {
		return fmt.Errorf("pipeline: export failed: %w", err)
	}

	return nil
}
