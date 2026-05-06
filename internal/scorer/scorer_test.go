package scorer_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/scorer"
)

func TestNew_NegativePenaltyReturnsError(t *testing.T) {
	_, err := scorer.New([]string{"KEY"}, nil, -1, 5)
	if err == nil {
		t.Fatal("expected error for negative penaltyEmpty")
	}
}

func TestNew_NegativeBonusReturnsError(t *testing.T) {
	_, err := scorer.New([]string{"KEY"}, nil, 5, -1)
	if err == nil {
		t.Fatal("expected error for negative bonusComplete")
	}
}

func TestEvaluate_NilEnvReturnsError(t *testing.T) {
	s, _ := scorer.New([]string{"KEY"}, nil, 5, 10)
	_, err := s.Evaluate(nil)
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestEvaluate_AllRequiredPresent(t *testing.T) {
	s, err := scorer.New([]string{"HOST", "PORT"}, nil, 5, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{"HOST": "localhost", "PORT": "8080"}
	res, err := s.Evaluate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 2 keys * 10 + 20 bonus = 40
	if res.Score != 40 {
		t.Errorf("expected score 40, got %d", res.Score)
	}
	if res.MaxScore != 40 {
		t.Errorf("expected max score 40, got %d", res.MaxScore)
	}
	if res.Percentage != 100.0 {
		t.Errorf("expected 100%%, got %.2f", res.Percentage)
	}
}

func TestEvaluate_MissingRequiredReducesScore(t *testing.T) {
	s, err := scorer.New([]string{"HOST", "PORT"}, nil, 5, 20)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{"HOST": "localhost"}
	res, err := s.Evaluate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// HOST=10, PORT missing=-5, no bonus; score=5
	if res.Score != 5 {
		t.Errorf("expected score 5, got %d", res.Score)
	}
	if len(res.Reasons) == 0 {
		t.Error("expected at least one reason")
	}
}

func TestEvaluate_SensitiveKeyNoted(t *testing.T) {
	s, err := scorer.New(nil, []string{"SECRET_KEY"}, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env := map[string]string{"SECRET_KEY": "abc123"}
	res, err := s.Evaluate(env)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	found := false
	for _, r := range res.Reasons {
		if r != "" {
			found = true
		}
	}
	if !found {
		t.Error("expected sensitive key reason to be recorded")
	}
}

func TestEvaluate_EmptyRequiredKeys_ZeroMaxScore(t *testing.T) {
	s, err := scorer.New(nil, nil, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	res, err := s.Evaluate(map[string]string{"FOO": "bar"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.MaxScore != 0 {
		t.Errorf("expected max score 0, got %d", res.MaxScore)
	}
	if res.Percentage != 0 {
		t.Errorf("expected percentage 0, got %.2f", res.Percentage)
	}
}
