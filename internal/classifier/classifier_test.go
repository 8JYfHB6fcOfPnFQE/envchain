package classifier_test

import (
	"testing"

	"github.com/yourorg/envchain/internal/classifier"
)

func TestNew_DefaultPatternsWork(t *testing.T) {
	c, err := classifier.New(nil, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil classifier")
	}
}

func TestNew_BlankExtraSecretPatternReturnsError(t *testing.T) {
	_, err := classifier.New([]string{""}, nil)
	if err == nil {
		t.Fatal("expected error for blank secret pattern")
	}
}

func TestNew_BlankExtraInternalPatternReturnsError(t *testing.T) {
	_, err := classifier.New(nil, []string{" "})
	if err == nil {
		t.Fatal("expected error for blank internal pattern")
	}
}

func TestClassify_SecretKey(t *testing.T) {
	c, _ := classifier.New(nil, nil)
	tests := []string{"DB_PASSWORD", "API_TOKEN", "AWS_SECRET_ACCESS_KEY", "PRIVATE_KEY"}
	for _, key := range tests {
		if got := c.Classify(key); got != classifier.LevelSecret {
			t.Errorf("Classify(%q) = %v, want secret", key, got)
		}
	}
}

func TestClassify_InternalKey(t *testing.T) {
	c, _ := classifier.New(nil, nil)
	tests := []string{"DB_HOST", "APP_PORT", "SERVICE_ENDPOINT", "DATABASE_DSN"}
	for _, key := range tests {
		if got := c.Classify(key); got != classifier.LevelInternal {
			t.Errorf("Classify(%q) = %v, want internal", key, got)
		}
	}
}

func TestClassify_PublicKey(t *testing.T) {
	c, _ := classifier.New(nil, nil)
	tests := []string{"APP_ENV", "LOG_LEVEL", "REGION", "VERSION"}
	for _, key := range tests {
		if got := c.Classify(key); got != classifier.LevelPublic {
			t.Errorf("Classify(%q) = %v, want public", key, got)
		}
	}
}

func TestClassify_CaseInsensitive(t *testing.T) {
	c, _ := classifier.New(nil, nil)
	if got := c.Classify("db_Password"); got != classifier.LevelSecret {
		t.Errorf("expected secret, got %v", got)
	}
}

func TestClassify_ExtraPatternOverride(t *testing.T) {
	c, _ := classifier.New([]string{"license"}, nil)
	if got := c.Classify("LICENSE_KEY"); got != classifier.LevelSecret {
		t.Errorf("expected secret for custom pattern, got %v", got)
	}
}

func TestClassifyMap_ReturnsAllKeys(t *testing.T) {
	c, _ := classifier.New(nil, nil)
	env := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_PORT":    "8080",
		"LOG_LEVEL":   "info",
	}
	result := c.ClassifyMap(env)
	if len(result) != len(env) {
		t.Fatalf("expected %d entries, got %d", len(env), len(result))
	}
	if result["DB_PASSWORD"] != classifier.LevelSecret {
		t.Errorf("DB_PASSWORD should be secret")
	}
	if result["APP_PORT"] != classifier.LevelInternal {
		t.Errorf("APP_PORT should be internal")
	}
	if result["LOG_LEVEL"] != classifier.LevelPublic {
		t.Errorf("LOG_LEVEL should be public")
	}
}

func TestLevel_String(t *testing.T) {
	if classifier.LevelSecret.String() != "secret" {
		t.Errorf("unexpected string for LevelSecret")
	}
	if classifier.LevelInternal.String() != "internal" {
		t.Errorf("unexpected string for LevelInternal")
	}
	if classifier.LevelPublic.String() != "public" {
		t.Errorf("unexpected string for LevelPublic")
	}
}
