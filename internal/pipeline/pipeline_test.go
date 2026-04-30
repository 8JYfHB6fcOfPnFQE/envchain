package pipeline_test

import (
	"bytes"
	"testing"

	"github.com/yourorg/envchain/internal/chain"
	"github.com/yourorg/envchain/internal/context"
	"github.com/yourorg/envchain/internal/exporter"
	"github.com/yourorg/envchain/internal/merger"
	"github.com/yourorg/envchain/internal/pipeline"
	"github.com/yourorg/envchain/internal/resolver"
	"github.com/yourorg/envchain/internal/validator"
)

func buildPipeline(t *testing.T, buf *bytes.Buffer) *pipeline.Pipeline {
	t.Helper()

	reg := context.New()
	ctx, _ := context.NewContext("base", map[string]string{"APP_ENV": "test", "PORT": "8080"})
	_ = reg.Register("base", ctx)

	ch := chain.New(reg)
	_ = ch.Add("base")

	m := merger.New()
	v := validator.New(nil)
	res, _ := resolver.New(resolver.Config{Chain: ch, Merger: m, Validator: v})

	exp, _ := exporter.New("dotenv", buf)

	p, _ := pipeline.New(pipeline.Config{Resolver: res, Exporter: exp})
	return p
}

func TestNew_NilResolverReturnsError(t *testing.T) {
	buf := &bytes.Buffer{}
	exp, _ := exporter.New("dotenv", buf)
	_, err := pipeline.New(pipeline.Config{Resolver: nil, Exporter: exp})
	if err == nil {
		t.Fatal("expected error for nil resolver")
	}
}

func TestNew_NilExporterReturnsError(t *testing.T) {
	reg := context.New()
	ch := chain.New(reg)
	m := merger.New()
	v := validator.New(nil)
	res, _ := resolver.New(resolver.Config{Chain: ch, Merger: m, Validator: v})

	_, err := pipeline.New(pipeline.Config{Resolver: res, Exporter: nil})
	if err == nil {
		t.Fatal("expected error for nil exporter")
	}
}

func TestRun_EmptyContextsReturnsError(t *testing.T) {
	buf := &bytes.Buffer{}
	p := buildPipeline(t, buf)
	if err := p.Run(nil); err == nil {
		t.Fatal("expected error for empty contexts")
	}
}

func TestRun_WritesOutput(t *testing.T) {
	buf := &bytes.Buffer{}
	p := buildPipeline(t, buf)
	if err := p.Run([]string{"base"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.Len() == 0 {
		t.Fatal("expected non-empty output from pipeline")
	}
}

func TestRun_UnknownContextReturnsError(t *testing.T) {
	buf := &bytes.Buffer{}
	p := buildPipeline(t, buf)
	if err := p.Run([]string{"nonexistent"}); err == nil {
		t.Fatal("expected error for unknown context")
	}
}
