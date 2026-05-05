package profiler

import (
	"errors"
	"fmt"
	"strings"
)

// MatchResult describes the outcome of matching an env map against a profile.
type MatchResult struct {
	Profile  string
	Missing  []string
	Extra    []string
	Matched  []string
}

// Match checks which keys from the profile are present in env, and which keys
// in env are not declared in the profile. Returns an error if the profile is nil.
func Match(profile *Profile, env map[string]string) (*MatchResult, error) {
	if profile == nil {
		return nil, errors.New("profiler: profile must not be nil")
	}
	if env == nil {
		env = map[string]string{}
	}

	profKeys := profile.Keys()
	profSet := make(map[string]struct{}, len(profKeys))
	for _, k := range profKeys {
		profSet[strings.ToUpper(k)] = struct{}{}
	}

	envSet := make(map[string]struct{}, len(env))
	for k := range env {
		envSet[strings.ToUpper(k)] = struct{}{}
	}

	var missing, extra, matched []string

	for _, k := range profKeys {
		uk := strings.ToUpper(k)
		if _, ok := envSet[uk]; ok {
			matched = append(matched, k)
		} else {
			missing = append(missing, k)
		}
	}

	for k := range env {
		uk := strings.ToUpper(k)
		if _, ok := profSet[uk]; !ok {
			extra = append(extra, k)
		}
	}

	return &MatchResult{
		Profile: profile.Name(),
		Missing: missing,
		Extra:   extra,
		Matched: matched,
	}, nil
}

// Summary returns a human-readable summary of the match result.
func (r *MatchResult) Summary() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "profile=%s matched=%d missing=%d extra=%d",
		r.Profile, len(r.Matched), len(r.Missing), len(r.Extra))
	return sb.String()
}

// IsComplete returns true when no keys are missing.
func (r *MatchResult) IsComplete() bool {
	return len(r.Missing) == 0
}
