// Package sorter provides utilities for ordering environment variable maps
// by key name, value length, or custom comparison functions.
package sorter

import (
	"errors"
	"sort"
)

// Order defines the sort direction.
type Order int

const (
	Ascending  Order = iota
	Descending
)

// Sorter sorts a map of environment variables into an ordered slice of pairs.
type Sorter struct {
	order    Order
	byValue  bool
}

// Pair holds a single key-value entry.
type Pair struct {
	Key   string
	Value string
}

// New creates a Sorter. order controls ascending/descending direction;
// byValue, when true, sorts by value length instead of key name.
func New(order Order, byValue bool) (*Sorter, error) {
	if order != Ascending && order != Descending {
		return nil, errors.New("sorter: invalid order value")
	}
	return &Sorter{order: order, byValue: byValue}, nil
}

// Sort returns the key-value pairs from env in the configured order.
// Returns an error if env is nil.
func (s *Sorter) Sort(env map[string]string) ([]Pair, error) {
	if env == nil {
		return nil, errors.New("sorter: env must not be nil")
	}

	pairs := make([]Pair, 0, len(env))
	for k, v := range env {
		pairs = append(pairs, Pair{Key: k, Value: v})
	}

	sort.Slice(pairs, func(i, j int) bool {
		var less bool
		if s.byValue {
			less = len(pairs[i].Value) < len(pairs[j].Value)
		} else {
			less = pairs[i].Key < pairs[j].Key
		}
		if s.order == Descending {
			return !less
		}
		return less
	})

	return pairs, nil
}

// Keys returns only the sorted key names from env.
func (s *Sorter) Keys(env map[string]string) ([]string, error) {
	pairs, err := s.Sort(env)
	if err != nil {
		return nil, err
	}
	keys := make([]string, len(pairs))
	for i, p := range pairs {
		keys[i] = p.Key
	}
	return keys, nil
}
