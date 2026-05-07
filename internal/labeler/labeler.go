package labeler

import (
	"errors"
	"fmt"
	"strings"
)

// Label represents a named annotation attached to an environment key.
type Label struct {
	Name  string
	Value string
}

// Labeler attaches and queries string labels on environment keys.
type Labeler struct {
	labels map[string][]Label
}

// New creates a new Labeler instance.
func New() *Labeler {
	return &Labeler{
		labels: make(map[string][]Label),
	}
}

// Attach adds a label with the given name and value to a key.
// Returns an error if key, name, or value is blank.
func (l *Labeler) Attach(key, name, value string) error {
	key = strings.TrimSpace(key)
	name = strings.TrimSpace(name)
	value = strings.TrimSpace(value)

	if key == "" {
		return errors.New("labeler: key must not be blank")
	}
	if name == "" {
		return errors.New("labeler: label name must not be blank")
	}
	if value == "" {
		return errors.New("labeler: label value must not be blank")
	}

	normKey := strings.ToUpper(key)
	l.labels[normKey] = append(l.labels[normKey], Label{Name: name, Value: value})
	return nil
}

// Get returns all labels attached to the given key.
// Returns nil if the key has no labels.
func (l *Labeler) Get(key string) []Label {
	normKey := strings.ToUpper(strings.TrimSpace(key))
	result, ok := l.labels[normKey]
	if !ok {
		return nil
	}
	copy := make([]Label, len(result))
	for i, lbl := range result {
		copy[i] = lbl
	}
	return copy
}

// FindByLabel returns all keys that have a label matching the given name and value.
func (l *Labeler) FindByLabel(name, value string) []string {
	var keys []string
	for k, lbls := range l.labels {
		for _, lbl := range lbls {
			if strings.EqualFold(lbl.Name, name) && strings.EqualFold(lbl.Value, value) {
				keys = append(keys, k)
				break
			}
		}
	}
	return keys
}

// Summary returns a human-readable summary of all labels for a key.
func (l *Labeler) Summary(key string) string {
	lbls := l.Get(key)
	if len(lbls) == 0 {
		return fmt.Sprintf("%s: (no labels)", strings.ToUpper(key))
	}
	parts := make([]string, len(lbls))
	for i, lbl := range lbls {
		parts[i] = fmt.Sprintf("%s=%s", lbl.Name, lbl.Value)
	}
	return fmt.Sprintf("%s: [%s]", strings.ToUpper(key), strings.Join(parts, ", "))
}
