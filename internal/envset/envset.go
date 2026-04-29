package envset

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// EnvSet represents a named set of environment variables with optional validation rules.
type EnvSet struct {
	Name     string
	Required []string
	Optional []string
	Values   map[string]string
}

// ErrMissingVar is returned when a required environment variable is not set.
type ErrMissingVar struct {
	VarName string
	SetName string
}

func (e *ErrMissingVar) Error() string {
	return fmt.Sprintf("envset %q: required variable %q is not set", e.SetName, e.VarName)
}

// New creates a new EnvSet with the given name.
func New(name string, required, optional []string) *EnvSet {
	return &EnvSet{
		Name:     name,
		Required: required,
		Optional: optional,
		Values:   make(map[string]string),
	}
}

// Load reads the environment variables defined in the set from the current process environment.
func (e *EnvSet) Load() {
	for _, key := range e.Required {
		if val, ok := os.LookupEnv(key); ok {
			e.Values[key] = val
		}
	}
	for _, key := range e.Optional {
		if val, ok := os.LookupEnv(key); ok {
			e.Values[key] = val
		}
	}
}

// Validate checks that all required variables are present in Values.
func (e *EnvSet) Validate() error {
	var errs []string
	for _, key := range e.Required {
		if _, ok := e.Values[key]; !ok {
			errs = append(errs, (&ErrMissingVar{VarName: key, SetName: e.Name}).Error())
		}
	}
	if len(errs) > 0 {
		return errors.New(strings.Join(errs, "; "))
	}
	return nil
}

// Get returns the value of a variable in the set, and whether it was found.
func (e *EnvSet) Get(key string) (string, bool) {
	val, ok := e.Values[key]
	return val, ok
}
