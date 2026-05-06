package tagger

import (
	"errors"
	"sort"
	"strings"
)

// Tag represents a label attached to an environment variable key.
type Tag struct {
	Name  string
	Value string
}

// Tagger assigns and retrieves tags for environment variable keys.
type Tagger struct {
	tags map[string][]Tag // key -> tags
}

// New creates a new Tagger instance.
func New() *Tagger {
	return &Tagger{
		tags: make(map[string][]Tag),
	}
}

// Tag attaches a named tag with an optional value to the given key.
// Returns an error if key or tag name is blank.
func (t *Tagger) Tag(key, name, value string) error {
	key = strings.TrimSpace(key)
	name = strings.TrimSpace(name)
	if key == "" {
		return errors.New("tagger: key must not be blank")
	}
	if name == "" {
		return errors.New("tagger: tag name must not be blank")
	}
	t.tags[key] = append(t.tags[key], Tag{Name: name, Value: value})
	return nil
}

// GetTags returns all tags associated with the given key.
// Returns nil if the key has no tags.
func (t *Tagger) GetTags(key string) []Tag {
	tags, ok := t.tags[strings.TrimSpace(key)]
	if !ok {
		return nil
	}
	result := make([]Tag, len(tags))
	copy(result, tags)
	return result
}

// KeysWithTag returns all keys that have a tag matching the given name.
func (t *Tagger) KeysWithTag(name string) []string {
	name = strings.TrimSpace(name)
	var keys []string
	for key, tags := range t.tags {
		for _, tag := range tags {
			if tag.Name == name {
				keys = append(keys, key)
				break
			}
		}
	}
	sort.Strings(keys)
	return keys
}

// Remove removes all tags from the given key.
func (t *Tagger) Remove(key string) {
	delete(t.tags, strings.TrimSpace(key))
}
