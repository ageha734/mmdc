package renderer_test

import (
	"context"
	"testing"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/renderer"
)

func TestNewChromedpRenderer(t *testing.T) {
	opts := renderer.RenderOptions{Quiet: false}
	r := renderer.NewChromedpRenderer(opts)
	if r == nil {
		t.Fatal("NewChromedpRenderer returned nil")
	}

	opts2 := renderer.RenderOptions{Quiet: true}
	r2 := renderer.NewChromedpRenderer(opts2)
	if r2 == nil {
		t.Fatal("NewChromedpRenderer returned nil")
	}
}

func TestChromedpRenderer_Render_UnsupportedFormat(t *testing.T) {
	opts := renderer.RenderOptions{Quiet: false}
	r := renderer.NewChromedpRenderer(opts)
	ctx := context.Background()
	cfg := config.NewRenderConfig()
	diagram := &domain.MermaidDiagram{Code: "graph TD\nA-->B"}

	_, err := r.Render(ctx, diagram, domain.OutputFormat("unsupported"), cfg)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestChromedpRenderer_Render_ContextCancelled(t *testing.T) {
	opts := renderer.RenderOptions{Quiet: false}
	r := renderer.NewChromedpRenderer(opts)
	cfg := config.NewRenderConfig()
	diagram := &domain.MermaidDiagram{Code: "graph TD\nA-->B"}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := r.Render(ctx, diagram, domain.OutputFormatSVG, cfg)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
}

func TestChromedpRenderer_Render_InvalidConfigFile(t *testing.T) {
	opts := renderer.RenderOptions{Quiet: false}
	r := renderer.NewChromedpRenderer(opts)
	ctx := context.Background()
	cfg := config.NewRenderConfig()
	cfg.ConfigFile = "/nonexistent/config.json"
	diagram := &domain.MermaidDiagram{Code: "graph TD\nA-->B"}

	_, err := r.Render(ctx, diagram, domain.OutputFormatSVG, cfg)
	if err == nil {
		t.Error("expected error for invalid config file")
	}
}

func TestChromedpRenderer_Render_AllFormats_ContextCancelled(t *testing.T) {
	opts := renderer.RenderOptions{Quiet: false}
	r := renderer.NewChromedpRenderer(opts)
	cfg := config.NewRenderConfig()
	diagram := &domain.MermaidDiagram{Code: "graph TD\nA-->B"}

	formats := []domain.OutputFormat{
		domain.OutputFormatSVG,
		domain.OutputFormatPNG,
		domain.OutputFormatPDF,
	}

	for _, format := range formats {
		t.Run(format.String(), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			_, err := r.Render(ctx, diagram, format, cfg)
			if err == nil {
				t.Error("expected error when context is cancelled")
			}
		})
	}
}

func TestChromedpRenderer_Render_QuietMode(t *testing.T) {
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	cfg := config.NewRenderConfig()
	diagram := &domain.MermaidDiagram{Code: "graph TD\nA-->B"}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := r.Render(ctx, diagram, domain.OutputFormatSVG, cfg)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
}
