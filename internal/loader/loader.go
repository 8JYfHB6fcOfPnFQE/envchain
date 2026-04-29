package loader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// FileLoader reads environment variable definitions from a .env file.
type FileLoader struct {
	path string
}

// New creates a new FileLoader for the given file path.
func New(path string) *FileLoader {
	return &FileLoader{path: path}
}

// Load reads key=value pairs from the file and returns them as a map.
// Lines starting with '#' and empty lines are ignored.
func (fl *FileLoader) Load() (map[string]string, error) {
	f, err := os.Open(fl.path)
	if err != nil {
		return nil, fmt.Errorf("loader: open %q: %w", fl.path, err)
	}
	defer f.Close()

	result := make(map[string]string)
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("loader: %q line %d: invalid format, expected KEY=VALUE", fl.path, lineNum)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Strip optional surrounding quotes from value
		if len(value) >= 2 {
			if (value[0] == '"' && value[len(value)-1] == '"') ||
				(value[0] == '\'' && value[len(value)-1] == '\'') {
				value = value[1 : len(value)-1]
			}
		}

		result[key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("loader: scan %q: %w", fl.path, err)
	}

	return result, nil
}
