package envset_test

import (
	"testing"

	"github.com/yourusername/envchain/internal/envset"
)

func TestNew(t *testing.T) {
	es := envset.New("test", []string{"FOO"}, []string{"BAR"})
	if es.Name != "test" {
		t.Errorf("expected name %q, got %q", "test", es.Name)
	}
	if len(es.Required) != 1 || es.Required[0] != "FOO" {
		t.Errorf("unexpected required vars: %v", es.Required)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	es := envset.New("prod", []string{"DATABASE_URL", "SECRET_KEY"}, nil)
	// Do not load from env; Values is empty
	err := es.Validate()
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
	msg := err.Error()
	if msg == "" {
		t.Fatal("expected non-empty error message")
	}
}

func TestValidate_AllPresent(t *testing.T) {
	es := envset.New("staging", []string{"APP_PORT"}, []string{"LOG_LEVEL"})
	es.Values["APP_PORT"] = "8080"
	if err := es.Validate(); err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestLoad_ReadsFromEnv(t *testing.T) {
	t.Setenv("ENVCHAIN_TEST_VAR", "hello")
	es := envset.New("load-test", []string{"ENVCHAIN_TEST_VAR"}, nil)
	es.Load()
	val, ok := es.Get("ENVCHAIN_TEST_VAR")
	if !ok {
		t.Fatal("expected variable to be loaded")
	}
	if val != "hello" {
		t.Errorf("expected %q, got %q", "hello", val)
	}
}

func TestGet_NotFound(t *testing.T) {
	es := envset.New("empty", nil, nil)
	_, ok := es.Get("NONEXISTENT")
	if ok {
		t.Error("expected variable to not be found")
	}
}
