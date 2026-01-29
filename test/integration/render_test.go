//go:build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/renderer"
)

func TestChromedpRenderer_RenderSVG(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `graph TD
    A[Start] --> B[Process]
    B --> C[End]`,
	}

	result, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
	if err != nil {
		t.Fatalf("Failed to render SVG: %v", err)
	}

	svgContent := string(result)
	if !strings.Contains(svgContent, "<svg") {
		t.Error("Result does not contain SVG tag")
	}
	if !strings.Contains(svgContent, "</svg>") {
		t.Error("Result does not contain closing SVG tag")
	}
}

func TestChromedpRenderer_RenderPNG(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `graph LR
    A --> B --> C`,
	}

	result, err := r.Render(ctx, diagram, domain.OutputFormatPNG, renderConfig)
	if err != nil {
		t.Fatalf("Failed to render PNG: %v", err)
	}

	if len(result) < 100 {
		t.Error("PNG result is too small")
	}

	if len(result) >= 4 && (result[0] != 0x89 || result[1] != 0x50 || result[2] != 0x4E || result[3] != 0x47) {
		t.Error("Result does not have PNG magic bytes")
	}
}

func TestChromedpRenderer_RenderPDF(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `sequenceDiagram
    Alice->>Bob: Hello
    Bob->>Alice: Hi`,
	}

	result, err := r.Render(ctx, diagram, domain.OutputFormatPDF, renderConfig)
	if err != nil {
		t.Fatalf("Failed to render PDF: %v", err)
	}

	if len(result) < 100 {
		t.Error("PDF result is too small")
	}

	if !strings.HasPrefix(string(result), "%PDF") {
		t.Error("Result does not have PDF magic bytes")
	}
}

func TestChromedpRenderer_RenderWithTheme(t *testing.T) {
	themes := []string{"default", "dark", "forest", "neutral"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			ctx := context.Background()
			renderConfig := config.NewRenderConfig()
			renderConfig.Theme = theme
			opts := renderer.RenderOptions{Quiet: true}
			r := renderer.NewChromedpRenderer(opts)

			diagram := &domain.MermaidDiagram{
				Code: `graph TD
    A --> B`,
			}

			result, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
			if err != nil {
				t.Fatalf("Failed to render SVG with theme %s: %v", theme, err)
			}

			if !strings.Contains(string(result), "<svg") {
				t.Errorf("Theme %s: Result does not contain SVG tag", theme)
			}
		})
	}
}

func TestChromedpRenderer_RenderWithCustomCSS(t *testing.T) {
	tmpDir := t.TempDir()
	cssFile := filepath.Join(tmpDir, "custom.css")

	cssContent := `.node rect { fill: #ff0000 !important; }`
	if err := os.WriteFile(cssFile, []byte(cssContent), 0644); err != nil {
		t.Fatalf("Failed to write CSS file: %v", err)
	}

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	renderConfig.CSSFile = cssFile
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `graph TD
    A[Test Node] --> B[Another Node]`,
	}

	result, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
	if err != nil {
		t.Fatalf("Failed to render SVG with custom CSS: %v", err)
	}

	if !strings.Contains(string(result), "<svg") {
		t.Error("Result does not contain SVG tag")
	}
}

func TestChromedpRenderer_RenderWithSVGID(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	renderConfig.SVGID = "custom-svg-id"
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `graph TD
    A --> B`,
	}

	result, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
	if err != nil {
		t.Fatalf("Failed to render SVG with custom ID: %v", err)
	}

	svgContent := string(result)
	if !strings.Contains(svgContent, "custom-svg-id") {
		t.Error("Result does not contain custom SVG ID")
	}
}

func TestChromedpRenderer_RenderWithBackgroundColor(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	renderConfig.BackgroundColor = "#f0f0f0"
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `graph TD
    A --> B`,
	}

	result, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
	if err != nil {
		t.Fatalf("Failed to render SVG with background color: %v", err)
	}

	if !strings.Contains(string(result), "<svg") {
		t.Error("Result does not contain SVG tag")
	}
}

func TestChromedpRenderer_RenderComplexDiagrams(t *testing.T) {
	testCases := []struct {
		name string
		code string
	}{
		{
			name: "Flowchart",
			code: `flowchart TB
    subgraph IDE1[IDE]
        direction TB
        subgraph SDK[SDK]
            A[Browser] --> B[Server]
        end
    end`,
		},
		{
			name: "Sequence",
			code: `sequenceDiagram
    participant Alice
    participant Bob
    Alice->>John: Hello John, how are you?
    loop Healthcheck
        John->>John: Fight against hypochondria
    end
    Note right of John: Rational thoughts <br/>prevail!
    John-->>Alice: Great!
    John->>Bob: How about you?
    Bob-->>John: Jolly good!`,
		},
		{
			name: "ClassDiagram",
			code: `classDiagram
    Animal <|-- Duck
    Animal <|-- Fish
    Animal <|-- Zebra
    Animal : +int age
    Animal : +String gender
    Animal: +isMammal()
    Animal: +mate()
    class Duck{
        +String beakColor
        +swim()
        +quack()
    }`,
		},
		{
			name: "StateDiagram",
			code: `stateDiagram-v2
    [*] --> Still
    Still --> [*]
    Still --> Moving
    Moving --> Still
    Moving --> Crash
    Crash --> [*]`,
		},
		{
			name: "ERDiagram",
			code: `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE-ITEM : contains
    CUSTOMER }|..|{ DELIVERY-ADDRESS : uses`,
		},
		{
			name: "Gantt",
			code: `gantt
    title A Gantt Diagram
    dateFormat  YYYY-MM-DD
    section Section
    A task           :a1, 2024-01-01, 30d
    Another task     :after a1, 20d`,
		},
		{
			name: "PieChart",
			code: `pie title Pets adopted by volunteers
    "Dogs" : 386
    "Cats" : 85
    "Rats" : 15`,
		},
		{
			name: "GitGraph",
			code: `gitGraph
    commit
    branch develop
    checkout develop
    commit
    commit
    checkout main
    merge develop
    commit`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			renderConfig := config.NewRenderConfig()
			opts := renderer.RenderOptions{Quiet: true}
			r := renderer.NewChromedpRenderer(opts)

			diagram := &domain.MermaidDiagram{Code: tc.code}

			result, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
			if err != nil {
				t.Fatalf("Failed to render %s diagram: %v", tc.name, err)
			}

			if !strings.Contains(string(result), "<svg") {
				t.Errorf("%s: Result does not contain SVG tag", tc.name)
			}
		})
	}
}

func TestChromedpRenderer_RenderInvalidDiagram(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `invalid mermaid syntax %%%`,
	}

	_, err := r.Render(ctx, diagram, domain.OutputFormatSVG, renderConfig)
	if err == nil {
		t.Error("Expected error for invalid diagram, but got nil")
	}
}

func TestChromedpRenderer_RenderUnsupportedFormat(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)

	diagram := &domain.MermaidDiagram{
		Code: `graph TD
    A --> B`,
	}

	_, err := r.Render(ctx, diagram, domain.OutputFormat("gif"), renderConfig)
	if err == nil {
		t.Error("Expected error for unsupported format, but got nil")
	}
}
