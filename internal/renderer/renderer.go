package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/yourorg/envchain/internal/exporter"
)

// Renderer combines a resolved variable map with an Exporter to produce
// formatted output, optionally filtering or masking sensitive keys.
type Renderer struct {
	vars    map[string]string
	maskSet map[string]bool
}

// New creates a Renderer from a resolved variable map.
func New(vars map[string]string) *Renderer {
	copy := make(map[string]string, len(vars))
	for k, v := range vars {
		copy[k] = v
	}
	return &Renderer{vars: copy, maskSet: make(map[string]bool)}
}

// Mask marks the given keys so their values are replaced with "***" on output.
func (r *Renderer) Mask(keys ...string) {
	for _, k := range keys {
		r.maskSet[strings.ToUpper(k)] = true
	}
}

// Render writes the variables to w using the specified format.
func (r *Renderer) Render(format exporter.Format, w io.Writer) error {
	e, err := exporter.New(format, w)
	if err != nil {
		return fmt.Errorf("renderer: %w", err)
	}
	output := make(map[string]string, len(r.vars))
	for k, v := range r.vars {
		if r.maskSet[strings.ToUpper(k)] {
			output[k] = "***"
		} else {
			output[k] = v
		}
	}
	return e.Write(output)
}

// Keys returns all variable keys held by the Renderer.
func (r *Renderer) Keys() []string {
	keys := make([]string, 0, len(r.vars))
	for k := range r.vars {
		keys = append(keys, k)
	}
	return keys
}
