package pipeline_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/yourorg/envchain/internal/chain"
	"github.com/yourorg/envchain/internal/context"
	"github.com/yourorg/envchain/internal/exporter"
	"github.com/yourorg/envchain/internal/merger"
	"github.com/yourorg/envchain/internal/pipeline"
	"github.com/yourorg/envchain/internal/resolver"
	"github.com/yourorg/envchain/internal/validator"
)

// TestIntegration_ChainedContextsOverride verifies that later contexts in the
// chain override earlier ones and that the final output is exported correctly.
func TestIntegration_ChainedContextsOverride(t *testing.T) {
	reg := context.New()

	base, _ := context.NewContext("base", map[string]string{
		"APP_ENV": "development",
		"LOG_LEVEL": "debug",
	})
	prod, _ := context.NewContext("prod", map[string]string{
		"APP_ENV": "production",
		"PORT": "443",
	})

	_ = reg.Register("base", base)
	_ = reg.Register("prod", prod)

	ch := chain.New(reg)
	_ = ch.Add("base")
	_ = ch.Add("prod")

	m := merger.New()
	v := validator.New(nil)
	res, err := resolver.New(resolver.Config{Chain: ch, Merger: m, Validator: v})
	if err != nil {
		t.Fatalf("resolver.New: %v", err)
	}

	buf := &bytes.Buffer{}
	exp, err := exporter.New("dotenv", buf)
	if err != nil {
		t.Fatalf("exporter.New: %v", err)
	}

	p, err := pipeline.New(pipeline.Config{Resolver: res, Exporter: exp})
	if err != nil {
		t.Fatalf("pipeline.New: %v", err)
	}

	if err := p.Run([]string{"base", "prod"}); err != nil {
		t.Fatalf("Run: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "APP_ENV=production") {
		t.Errorf("expected APP_ENV=production in output, got:\n%s", out)
	}
	if !strings.Contains(out, "LOG_LEVEL=debug") {
		t.Errorf("expected LOG_LEVEL=debug in output, got:\n%s", out)
	}
	if !strings.Contains(out, "PORT=443") {
		t.Errorf("expected PORT=443 in output, got:\n%s", out)
	}
}
