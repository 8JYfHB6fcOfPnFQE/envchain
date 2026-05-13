// Package aliaser provides support for defining and resolving key aliases
// within an environment variable set. An alias maps one or more alternate
// names to a canonical key, enabling consumers to look up values by any
// recognised name.
package aliaser

import (
	"errors"
	"strings"
)

// Aliaser holds the mapping from alias names to canonical keys.
type Aliaser struct {
	// aliases maps lowercase alias → canonical key (original case preserved).
	aliases map[string]string
}

// New returns an empty Aliaser ready for use.
func New() *Aliaser {
	return &Aliaser{aliases: make(map[string]string)}
}

// Register associates one or more alias names with a canonical key.
// Both the canonical key and each alias are matched case-insensitively.
// Returns an error if canonical is blank or any alias is blank.
func (a *Aliaser) Register(canonical string, aliases ...string) error {
	if strings.TrimSpace(canonical) == "" {
		return errors.New("aliaser: canonical key must not be blank")
	}
	for _, alias := range aliases {
		if strings.TrimSpace(alias) == "" {
			return errors.New("aliaser: alias must not be blank")
		}
		a.aliases[strings.ToLower(alias)] = canonical
	}
	// Also register the canonical key as an alias of itself for uniform lookup.
	a.aliases[strings.ToLower(canonical)] = canonical
	return nil
}

// Resolve returns the canonical key for the given name (which may itself be
// the canonical key or any registered alias). The second return value reports
// whether a mapping was found.
func (a *Aliaser) Resolve(name string) (string, bool) {
	canonical, ok := a.aliases[strings.ToLower(name)]
	return canonical, ok
}

// Lookup resolves name to its canonical key and then retrieves the
// corresponding value from env. Returns the value and true on success, or
// an empty string and false when the alias is unknown or the key is absent
// from env.
func (a *Aliaser) Lookup(name string, env map[string]string) (string, bool) {
	canonical, ok := a.Resolve(name)
	if !ok {
		return "", false
	}
	val, found := env[canonical]
	return val, found
}

// Aliases returns all alias names registered for the given canonical key.
// The canonical key itself is not included in the returned slice.
func (a *Aliaser) Aliases(canonical string) []string {
	var result []string
	for alias, canon := range a.aliases {
		if strings.EqualFold(canon, canonical) && !strings.EqualFold(alias, canonical) {
			result = append(result, alias)
		}
	}
	return result
}
