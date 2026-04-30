// Package pipeline provides a high-level orchestration layer for envchain.
//
// A Pipeline combines a Resolver and an Exporter to implement the full
// environment variable processing flow in a single Run call:
//
//  1. The Resolver walks the configured chain of context names, merging
//     each context's variables in order and validating the final result.
//
//  2. The Exporter writes the resolved key-value map to the configured
//     output in the chosen format (dotenv, export, or JSON).
//
// Example usage:
//
//	p, err := pipeline.New(pipeline.Config{
//	    Resolver: res,
//	    Exporter: exp,
//	})
//	if err != nil { ... }
//	if err := p.Run([]string{"base", "production"}); err != nil { ... }
package pipeline
