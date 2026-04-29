package exporter

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Format represents the output format for exported environment variables.
type Format string

const (
	FormatDotenv Format = "dotenv"
	FormatExport Format = "export"
	FormatJSON   Format = "json"
)

// Exporter writes environment variable sets to an output stream.
type Exporter struct {
	format Format
	writer io.Writer
}

// New creates a new Exporter with the given format and writer.
func New(format Format, w io.Writer) (*Exporter, error) {
	switch format {
	case FormatDotenv, FormatExport, FormatJSON:
		// valid
	default:
		return nil, fmt.Errorf("unsupported format: %q", format)
	}
	return &Exporter{format: format, writer: w}, nil
}

// Write outputs the provided key-value pairs in the configured format.
func (e *Exporter) Write(vars map[string]string) error {
	keys := sortedKeys(vars)
	switch e.format {
	case FormatDotenv:
		return e.writeDotenv(keys, vars)
	case FormatExport:
		return e.writeExport(keys, vars)
	case FormatJSON:
		return e.writeJSON(keys, vars)
	}
	return nil
}

func (e *Exporter) writeDotenv(keys []string, vars map[string]string) error {
	for _, k := range keys {
		_, err := fmt.Fprintf(e.writer, "%s=%s\n", k, quoteValue(vars[k]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeExport(keys []string, vars map[string]string) error {
	for _, k := range keys {
		_, err := fmt.Fprintf(e.writer, "export %s=%s\n", k, quoteValue(vars[k]))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) writeJSON(keys []string, vars map[string]string) error {
	var sb strings.Builder
	sb.WriteString("{\n")
	for i, k := range keys {
		sb.WriteString(fmt.Sprintf("  %q: %q", k, vars[k]))
		if i < len(keys)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("}\n")
	_, err := fmt.Fprint(e.writer, sb.String())
	return err
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func quoteValue(v string) string {
	if strings.ContainsAny(v, " \t\n") {
		return fmt.Sprintf("%q", v)
	}
	return v
}
