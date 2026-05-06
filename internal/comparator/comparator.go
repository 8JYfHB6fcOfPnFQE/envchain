package comparator

import (
	"errors"
	"sort"
)

// Result holds the outcome of comparing two environment maps.
type Result struct {
	Matching  []string
	Conflicts map[string][2]string // key -> [left, right]
	LeftOnly  []string
	RightOnly []string
}

// Comparator compares two named environment variable maps.
type Comparator struct {
	leftName  string
	rightName string
}

// New creates a Comparator with the given context names used in output labelling.
func New(leftName, rightName string) (*Comparator, error) {
	if leftName == "" {
		return nil, errors.New("comparator: leftName must not be empty")
	}
	if rightName == "" {
		return nil, errors.New("comparator: rightName must not be empty")
	}
	return &Comparator{leftName: leftName, rightName: rightName}, nil
}

// Compare performs a key-by-key comparison of left and right env maps.
func (c *Comparator) Compare(left, right map[string]string) (*Result, error) {
	if left == nil {
		return nil, errors.New("comparator: left env must not be nil")
	}
	if right == nil {
		return nil, errors.New("comparator: right env must not be nil")
	}

	result := &Result{
		Conflicts: make(map[string][2]string),
	}

	seen := make(map[string]bool)

	for k, lv := range left {
		seen[k] = true
		rv, ok := right[k]
		if !ok {
			result.LeftOnly = append(result.LeftOnly, k)
		} else if lv == rv {
			result.Matching = append(result.Matching, k)
		} else {
			result.Conflicts[k] = [2]string{lv, rv}
		}
	}

	for k := range right {
		if !seen[k] {
			result.RightOnly = append(result.RightOnly, k)
		}
	}

	sort.Strings(result.Matching)
	sort.Strings(result.LeftOnly)
	sort.Strings(result.RightOnly)

	return result, nil
}

// LeftName returns the label for the left environment.
func (c *Comparator) LeftName() string { return c.leftName }

// RightName returns the label for the right environment.
func (c *Comparator) RightName() string { return c.rightName }
