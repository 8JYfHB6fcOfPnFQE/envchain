package grouper_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/grouper"
)

func TestAdd_ValidEntry(t *testing.T) {
	g := grouper.New()
	if err := g.Add("db", "DB_HOST"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	list := g.List()
	if len(list) != 1 || list[0].Name != "db" {
		t.Fatalf("expected group 'db', got %v", list)
	}
}

func TestAdd_BlankNameReturnsError(t *testing.T) {
	g := grouper.New()
	if err := g.Add("", "DB_HOST"); err == nil {
		t.Fatal("expected error for blank name")
	}
}

func TestAdd_BlankKeyReturnsError(t *testing.T) {
	g := grouper.New()
	if err := g.Add("db", ""); err == nil {
		t.Fatal("expected error for blank key")
	}
}

func TestList_ReturnsSortedGroups(t *testing.T) {
	g := grouper.New()
	_ = g.Add("redis", "REDIS_PORT")
	_ = g.Add("app", "APP_ENV")
	_ = g.Add("db", "DB_HOST")
	list := g.List()
	names := make([]string, len(list))
	for i, gr := range list {
		names[i] = gr.Name
	}
	expected := []string{"app", "db", "redis"}
	for i, n := range expected {
		if names[i] != n {
			t.Errorf("position %d: want %q got %q", i, n, names[i])
		}
	}
}

func TestGroupByPrefix_MatchesPrefix(t *testing.T) {
	g := grouper.New()
	_ = g.Add("db", "placeholder")
	env := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"APP_ENV": "production",
	}
	result := g.GroupByPrefix(env)
	if len(result["db"]) != 2 {
		t.Errorf("expected 2 db keys, got %d", len(result["db"]))
	}
	if len(result["__ungrouped__"]) != 1 {
		t.Errorf("expected 1 ungrouped key, got %d", len(result["__ungrouped__"]))
	}
}

func TestGroupByPrefix_NilEnvReturnsEmptyMap(t *testing.T) {
	g := grouper.New()
	_ = g.Add("db", "DB_HOST")
	result := g.GroupByPrefix(nil)
	if len(result) != 0 {
		t.Errorf("expected empty result for nil env, got %v", result)
	}
}

func TestGroupByPrefix_CaseInsensitivePrefix(t *testing.T) {
	g := grouper.New()
	_ = g.Add("App", "placeholder")
	env := map[string]string{"APP_DEBUG": "true", "OTHER": "val"}
	result := g.GroupByPrefix(env)
	if len(result["App"]) != 1 || result["App"][0] != "APP_DEBUG" {
		t.Errorf("expected APP_DEBUG in App group, got %v", result)
	}
}
