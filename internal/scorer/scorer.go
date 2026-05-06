// Package scorer evaluates the health and completeness of a resolved
// environment variable set by assigning a numeric score based on coverage,
// validation results, and classification outcomes.
package scorer

import (
	"errors"
	"fmt"
)

// Result holds the outcome of a scoring evaluation.
type Result struct {
	Score      int
	MaxScore   int
	Percentage float64
	Reasons    []string
}

// Scorer evaluates an environment map and produces a Result.
type Scorer struct {
	requiredKeys  []string
	sensitiveKeys []string
	penaltyEmpty  int
	bonusComplete int
}

// New creates a Scorer with the given required and sensitive key lists.
// penaltyEmpty is subtracted per missing required key.
// bonusComplete is added when all required keys are present.
func New(requiredKeys, sensitiveKeys []string, penaltyEmpty, bonusComplete int) (*Scorer, error) {
	if penaltyEmpty < 0 {
		return nil, errors.New("scorer: penaltyEmpty must be non-negative")
	}
	if bonusComplete < 0 {
		return nil, errors.New("scorer: bonusComplete must be non-negative")
	}
	return &Scorer{
		requiredKeys:  requiredKeys,
		sensitiveKeys: sensitiveKeys,
		penaltyEmpty:  penaltyEmpty,
		bonusComplete: bonusComplete,
	}, nil
}

// Evaluate scores the provided environment map.
func (s *Scorer) Evaluate(env map[string]string) (*Result, error) {
	if env == nil {
		return nil, errors.New("scorer: env map must not be nil")
	}

	maxScore := len(s.requiredKeys)*10 + s.bonusComplete
	score := 0
	var reasons []string

	missingCount := 0
	for _, key := range s.requiredKeys {
		val, ok := env[key]
		if !ok || val == "" {
			missingCount++
			penalty := s.penaltyEmpty
			score -= penalty
			reasons = append(reasons, fmt.Sprintf("missing required key %q (-%d)", key, penalty))
		} else {
			score += 10
		}
	}

	if missingCount == 0 && len(s.requiredKeys) > 0 {
		score += s.bonusComplete
		reasons = append(reasons, fmt.Sprintf("all required keys present (+%d bonus)", s.bonusComplete))
	}

	for _, key := range s.sensitiveKeys {
		if val, ok := env[key]; ok && val != "" {
			reasons = append(reasons, fmt.Sprintf("sensitive key %q is set", key))
		}
	}

	var pct float64
	if maxScore > 0 {
		pct = float64(score) / float64(maxScore) * 100
	}

	return &Result{
		Score:      score,
		MaxScore:   maxScore,
		Percentage: pct,
		Reasons:    reasons,
	}, nil
}
