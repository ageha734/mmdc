package usecase_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/file"
	mockrenderer "github.com/ageha734/mmdc/test/mock/infra/renderer"
	"github.com/ageha734/mmdc/usecase"
)

func TestNewMarkdownUsecase(t *testing.T) {
	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	if u == nil {
		t.Fatal("NewMarkdownUsecase returned nil")
	}
}

func TestMarkdownUsecase_ProcessMarkdownFile_NoMermaid(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	content := "普通のMarkdownテキスト"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockL.infoCalled {
		t.Error("expected Info to be called")
	}
}

func TestMarkdownUsecase_ProcessMarkdownFile_WithMermaid(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	renderCount := 0
	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			renderCount++
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	t.Skip("Skipping test that requires actual browser (markdown.go creates ChromedpRenderer internally)")

	err := u.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", cfg, false)
	_ = err
	_ = renderCount
}

func TestMarkdownUsecase_ProcessMarkdownFile_WithArtefactsPath(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	artefactsPath := filepath.Join(tmpDir, "artefacts")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, artefactsPath, cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if _, err := os.Stat(artefactsPath); os.IsNotExist(err) {
		t.Error("artefacts directory should be created")
	}
}

func TestMarkdownUsecase_ProcessMarkdownFile_NonMarkdownOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.txt")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", cfg, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMarkdownUsecase_ProcessMarkdownFile_ReadError(t *testing.T) {
	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.ProcessMarkdownFile(ctx, "/nonexistent/input.md", "output.md", domain.OutputFormatSVG, "", cfg, false)
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestMarkdownUsecase_ProcessMarkdownFile_QuietMode(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	content := "普通のMarkdownテキスト"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", cfg, true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockL.infoCalled {
		t.Error("Info should not be called in quiet mode")
	}
}

func TestMarkdownUsecase_ProcessMarkdownFile_RenderError(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return nil, fmt.Errorf("render failed")
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewMarkdownUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.ProcessMarkdownFile(ctx, inputFile, outputFile, domain.OutputFormatSVG, "", cfg, true)
	if err == nil {
		t.Fatal("expected error")
	}
}
