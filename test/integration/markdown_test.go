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
	"github.com/ageha734/mmdc/infra/file"
	"github.com/ageha734/mmdc/infra/renderer"
	mockLogger "github.com/ageha734/mmdc/test/mock/infra/logger"
	"github.com/ageha734/mmdc/usecase"
)

func TestMarkdownUsecase_ProcessMarkdownWithMermaid(t *testing.T) {
	tmpDir := t.TempDir()

	mdContent := `# Test Document

This is a test document with Mermaid diagrams.

` + "```mermaid" + `
graph TD
    A[Start] --> B[Process]
    B --> C[End]
` + "```" + `

Some text between diagrams.

` + "```mermaid" + `
sequenceDiagram
    Alice->>Bob: Hello
    Bob->>Alice: Hi
` + "```" + `

End of document.
`

	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	mdUsecase := usecase.NewMarkdownUsecase(r, fileRW, log)
	err := mdUsecase.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", renderConfig, true)
	if err != nil {
		t.Fatalf("Failed to process markdown file: %v", err)
	}

	outputContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	content := string(outputContent)
	if strings.Contains(content, "```mermaid") {
		t.Error("Output markdown should not contain mermaid code blocks")
	}
	if !strings.Contains(content, "![") {
		t.Error("Output markdown should contain image references")
	}
}

func TestMarkdownUsecase_ProcessMarkdownWithArtefacts(t *testing.T) {
	tmpDir := t.TempDir()
	artefactsDir := filepath.Join(tmpDir, "images")

	mdContent := `# Test

` + "```mermaid" + `
graph TD
    A --> B
` + "```" + `
`

	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	mdUsecase := usecase.NewMarkdownUsecase(r, fileRW, log)
	err := mdUsecase.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, artefactsDir, renderConfig, true)
	if err != nil {
		t.Fatalf("Failed to process markdown file: %v", err)
	}

	if _, err := os.Stat(artefactsDir); os.IsNotExist(err) {
		t.Error("Artefacts directory should be created")
	}

	files, err := os.ReadDir(artefactsDir)
	if err != nil {
		t.Fatalf("Failed to read artefacts directory: %v", err)
	}

	svgFound := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".svg") {
			svgFound = true
			break
		}
	}
	if !svgFound {
		t.Error("SVG file should exist in artefacts directory")
	}
}

func TestMarkdownUsecase_ProcessMarkdownNoMermaid(t *testing.T) {
	tmpDir := t.TempDir()

	mdContent := `# Test Document

This is a plain markdown document without mermaid diagrams.

- Item 1
- Item 2
- Item 3
`

	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	mdUsecase := usecase.NewMarkdownUsecase(r, fileRW, log)
	err := mdUsecase.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", renderConfig, true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if _, err := os.Stat(outputFile); !os.IsNotExist(err) {
		t.Error("Output file should not be created when there are no mermaid diagrams")
	}
}

func TestMarkdownUsecase_ProcessMarkdownToPNG(t *testing.T) {
	tmpDir := t.TempDir()

	mdContent := `# Test

` + "```mermaid" + `
graph TD
    A --> B
` + "```" + `
`

	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")
	artefactsDir := filepath.Join(tmpDir, "images")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	mdUsecase := usecase.NewMarkdownUsecase(r, fileRW, log)
	err := mdUsecase.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatPNG, artefactsDir, renderConfig, true)
	if err != nil {
		t.Fatalf("Failed to process markdown file: %v", err)
	}

	files, err := os.ReadDir(artefactsDir)
	if err != nil {
		t.Fatalf("Failed to read artefacts directory: %v", err)
	}

	pngFound := false
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".png") {
			pngFound = true

			pngPath := filepath.Join(artefactsDir, f.Name())
			pngData, err := os.ReadFile(pngPath)
			if err != nil {
				t.Fatalf("Failed to read PNG file: %v", err)
			}
			if len(pngData) < 4 || pngData[0] != 0x89 || pngData[1] != 0x50 || pngData[2] != 0x4E || pngData[3] != 0x47 {
				t.Error("PNG file does not have valid PNG magic bytes")
			}
			break
		}
	}
	if !pngFound {
		t.Error("PNG file should exist in artefacts directory")
	}
}

func TestMarkdownUsecase_ProcessMarkdownMultipleDiagrams(t *testing.T) {
	tmpDir := t.TempDir()
	artefactsDir := filepath.Join(tmpDir, "images")

	mdContent := `# Multiple Diagrams

` + "```mermaid" + `
graph TD
    A --> B
` + "```" + `

` + "```mermaid" + `
sequenceDiagram
    A->>B: Hello
` + "```" + `

` + "```mermaid" + `
pie
    "A" : 50
    "B" : 50
` + "```" + `
`

	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	mdUsecase := usecase.NewMarkdownUsecase(r, fileRW, log)
	err := mdUsecase.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, artefactsDir, renderConfig, true)
	if err != nil {
		t.Fatalf("Failed to process markdown file: %v", err)
	}

	files, err := os.ReadDir(artefactsDir)
	if err != nil {
		t.Fatalf("Failed to read artefacts directory: %v", err)
	}

	svgCount := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".svg") {
			svgCount++
		}
	}
	if svgCount != 3 {
		t.Errorf("Expected 3 SVG files, got %d", svgCount)
	}
}

func TestMarkdownUsecase_ProcessMarkdownWithInvalidDiagram(t *testing.T) {
	tmpDir := t.TempDir()

	mdContent := `# Test

` + "```mermaid" + `
invalid mermaid syntax %%%
` + "```" + `
`

	inputFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(inputFile, []byte(mdContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	mdUsecase := usecase.NewMarkdownUsecase(r, fileRW, log)
	err := mdUsecase.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", renderConfig, true)
	if err == nil {
		t.Error("Expected error for invalid mermaid diagram, but got nil")
	}
}
