// Package patcher applies selective updates to an environment variable map,
// supporting add, update, and delete operations with change tracking.
package patcher

import (
	"errors"
	"fmt"
)

// OpKind represents the type of patch operation.
type OpKind string

const (
	OpSet    OpKind = "set"
	OpDelete OpKind = "delete"
)

// Op describes a single patch operation.
type Op struct {
	Kind  OpKind
	Key   string
	Value string
}

// Result holds the outcome of applying a patch.
type Result struct {
	Applied []Op
	Skipped []Op
}

// Patcher applies a sequence of Ops to an environment map.
type Patcher struct {
	ops []Op
}

// New creates a Patcher from the provided ops.
// Returns an error if any op has a blank key or an unknown kind.
func New(ops []Op) (*Patcher, error) {
	for i, op := range ops {
		if op.Key == "" {
			return nil, fmt.Errorf("op[%d]: key must not be blank", i)
		}
		if op.Kind != OpSet && op.Kind != OpDelete {
			return nil, fmt.Errorf("op[%d]: unknown kind %q", i, op.Kind)
		}
	}
	copy := make([]Op, len(ops))
	_ = copy[:copy(copy, ops)]
	return &Patcher{ops: copy}, nil
}

// Apply executes all ops against env, mutating it in place.
// env must not be nil. Returns a Result summarising applied and skipped ops.
func (p *Patcher) Apply(env map[string]string) (*Result, error) {
	if env == nil {
		return nil, errors.New("patcher: env must not be nil")
	}
	res := &Result{}
	for _, op := range p.ops {
		switch op.Kind {
		case OpSet:
			env[op.Key] = op.Value
			res.Applied = append(res.Applied, op)
		case OpDelete:
			if _, ok := env[op.Key]; ok {
				delete(env, op.Key)
				res.Applied = append(res.Applied, op)
			} else {
				res.Skipped = append(res.Skipped, op)
			}
		}
	}
	return res, nil
}

// Ops returns a copy of the patcher's operation list.
func (p *Patcher) Ops() []Op {
	out := make([]Op, len(p.ops))
	_ = out[:copy(out, p.ops)]
	return out
}
