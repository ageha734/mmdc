package usecase_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/file"
	rendererpkg "github.com/ageha734/mmdc/infra/renderer"
	mockrenderer "github.com/ageha734/mmdc/test/mock/infra/renderer"
	"github.com/ageha734/mmdc/usecase"
)

var _ rendererpkg.Renderer = (*mockrenderer.MockRenderer)(nil)

type mockLogger struct {
	infoCalled    bool
	successCalled bool
	warnCalled    bool
	errorExitFunc func(msg string)
}

func (m *mockLogger) Info(msg string) {
	m.infoCalled = true
}

func (m *mockLogger) Success(msg string) {
	m.successCalled = true
}

func (m *mockLogger) Warn(msg string) {
	m.warnCalled = true
}

func (m *mockLogger) ErrorExit(msg string) {
	if m.errorExitFunc != nil {
		m.errorExitFunc(msg)
	} else {
		panic("ErrorExit called: " + msg)
	}
}

func TestNewRenderUsecase(t *testing.T) {
	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	if u == nil {
		t.Fatal("NewRenderUsecase returned nil")
	}
}

func TestRenderUsecase_RenderFromInput_SingleFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	outputFile := filepath.Join(tmpDir, "output.svg")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			if diagram.Code != mermaidCode {
				t.Errorf("expected code %q, got %q", mermaidCode, diagram.Code)
			}
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, cfg, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockL.infoCalled {
		t.Error("expected Info to be called")
	}
	if !mockL.successCalled {
		t.Error("expected Success to be called")
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("output file should be created")
	}
}

func TestRenderUsecase_RenderFromInput_MarkdownFile(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	markdownContent := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(markdownContent), 0o644); err != nil {
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

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, cfg, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if renderCount == 0 {
		t.Error("expected render to be called")
	}
}

func TestRenderUsecase_RenderFromInput_EmptyInputPath(t *testing.T) {
	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	if u == nil {
		t.Fatal("NewRenderUsecase returned nil")
	}
}

func TestRenderUsecase_RenderFromInput_NonExistentFile(t *testing.T) {
	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	errorExitCalled := false
	mockL := &mockLogger{
		errorExitFunc: func(msg string) {
			errorExitCalled = true
		},
	}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	defer func() {
		if r := recover(); r != nil {
			if !errorExitCalled {
				t.Error("expected ErrorExit to be called")
			}
		}
	}()

	u.RenderFromInput(ctx, "/nonexistent/file.mmd", "", domain.OutputFormatSVG, cfg, "", false)
}

func TestRenderUsecase_RenderFromInput_QuietMode(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	outputFile := filepath.Join(tmpDir, "output.svg")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockL.infoCalled {
		t.Error("Info should not be called in quiet mode")
	}
	if mockL.successCalled {
		t.Error("Success should not be called in quiet mode")
	}
}

func TestRenderUsecase_RenderFromInput_FormatFromExtension(t *testing.T) {
	tests := []struct {
		name       string
		outputFile string
		wantFormat domain.OutputFormat
	}{
		{
			name:       "PNG format",
			outputFile: "output.png",
			wantFormat: domain.OutputFormatPNG,
		},
		{
			name:       "PDF format",
			outputFile: "output.pdf",
			wantFormat: domain.OutputFormatPDF,
		},
		{
			name:       "SVG format",
			outputFile: "output.svg",
			wantFormat: domain.OutputFormatSVG,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputFile := filepath.Join(tmpDir, "input.mmd")
			outputFile := filepath.Join(tmpDir, tt.outputFile)
			mermaidCode := "graph TD\nA-->B"

			if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
				t.Fatalf("failed to create input file: %v", err)
			}

			var capturedFormat domain.OutputFormat
			mockR := &mockrenderer.MockRenderer{
				RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
					capturedFormat = format
					return []byte("rendered"), nil
				},
			}
			frw := file.NewFileReaderWriter()
			mockL := &mockLogger{}

			u := usecase.NewRenderUsecase(mockR, frw, mockL)
			ctx := context.Background()
			cfg := config.NewRenderConfig()

			err := u.RenderFromInput(ctx, inputFile, outputFile, "", cfg, "", true)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if capturedFormat != tt.wantFormat {
				t.Errorf("expected format %v, got %v", tt.wantFormat, capturedFormat)
			}
		})
	}
}

func TestRenderUsecase_RenderFromInput_RenderError(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	outputFile := filepath.Join(tmpDir, "output.svg")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return nil, fmt.Errorf("render failed")
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, cfg, "", true)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestRenderUsecase_RenderFromInput_InvalidOutputDir(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	outputFile := "/nonexistent/dir/output.svg"
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	errorExitCalled := false
	mockL := &mockLogger{
		errorExitFunc: func(msg string) {
			errorExitCalled = true
		},
	}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	defer func() {
		recover()
		if !errorExitCalled {
			t.Error("expected ErrorExit to be called for invalid output dir")
		}
	}()

	u.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, cfg, "", true)
}

func TestRenderUsecase_RenderFromInput_InvalidOutputExtension(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	outputFile := filepath.Join(tmpDir, "output.invalid")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	errorExitCalled := false
	mockL := &mockLogger{
		errorExitFunc: func(msg string) {
			errorExitCalled = true
		},
	}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	defer func() {
		recover()
		if !errorExitCalled {
			t.Error("expected ErrorExit to be called for invalid output extension")
		}
	}()

	u.RenderFromInput(ctx, inputFile, outputFile, "", cfg, "", true)
}

func TestRenderUsecase_RenderFromInput_MarkdownWithOutputStdin(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	errorExitCalled := false
	mockL := &mockLogger{
		errorExitFunc: func(msg string) {
			errorExitCalled = true
		},
	}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	defer func() {
		recover()
		if !errorExitCalled {
			t.Error("expected ErrorExit to be called for markdown with stdout output")
		}
	}()

	u.RenderFromInput(ctx, inputFile, "-", domain.OutputFormatSVG, cfg, "", true)
}

func TestRenderUsecase_RenderFromInput_MarkdownNoOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
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

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, "", domain.OutputFormatSVG, cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRenderUsecase_RenderFromInput_MarkdownFormatFromExtension(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedFormat domain.OutputFormat
	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			capturedFormat = format
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, "", cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedFormat != domain.OutputFormatSVG {
		t.Errorf("expected SVG format for .md output, got %v", capturedFormat)
	}
}

func TestRenderUsecase_RenderFromInput_MarkdownFormatFromMarkdownExtension(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.markdown")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedFormat domain.OutputFormat
	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			capturedFormat = format
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, "", cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedFormat != domain.OutputFormatSVG {
		t.Errorf("expected SVG format for .markdown output, got %v", capturedFormat)
	}
}

func TestRenderUsecase_RenderFromInput_MarkdownFormatExplicit(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.md")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedFormat domain.OutputFormat
	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			capturedFormat = format
			return []byte("rendered"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatPNG, cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedFormat != domain.OutputFormatPNG {
		t.Errorf("expected PNG format when explicitly set, got %v", capturedFormat)
	}
}

func TestRenderUsecase_RenderFromInput_MarkdownWithPNGOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.md")
	outputFile := filepath.Join(tmpDir, "output.png")
	content := "```mermaid\ngraph TD\nA-->B\n```"

	if err := os.WriteFile(inputFile, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedFormat domain.OutputFormat
	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			capturedFormat = format
			return []byte("rendered"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, outputFile, "", cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedFormat != domain.OutputFormatPNG {
		t.Errorf("expected PNG format for .png output, got %v", capturedFormat)
	}
}

func TestRenderUsecase_RenderFromInput_EmptyInputPath_Warning(t *testing.T) {
	mockR := mockrenderer.NewMockRenderer()
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	_ = u.RenderFromInput(ctx, "", "", domain.OutputFormatSVG, cfg, "", false)

	if !mockL.warnCalled {
		t.Error("expected Warn to be called for empty input path")
	}
}

func TestRenderUsecase_RenderFromInput_StdoutOutput(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedOutput bytes.Buffer
	restore := usecase.SetOsStdoutForTesting(&capturedOutput)
	defer restore()

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, "-", domain.OutputFormatSVG, cfg, "", true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedOutput.String() != "<svg>rendered</svg>" {
		t.Errorf("expected '<svg>rendered</svg>', got %q", capturedOutput.String())
	}
}

func TestRenderUsecase_RenderFromInput_StdoutOutput_NoFormat(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedOutput bytes.Buffer
	restore := usecase.SetOsStdoutForTesting(&capturedOutput)
	defer restore()

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return []byte("<svg>rendered</svg>"), nil
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, "-", "", cfg, "", false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !mockL.warnCalled {
		t.Error("expected Warn to be called for no output format")
	}

	if capturedOutput.String() != "<svg>rendered</svg>" {
		t.Errorf("expected '<svg>rendered</svg>', got %q", capturedOutput.String())
	}
}

func TestRenderUsecase_RenderFromInput_StdoutOutput_RenderError(t *testing.T) {
	tmpDir := t.TempDir()
	inputFile := filepath.Join(tmpDir, "input.mmd")
	mermaidCode := "graph TD\nA-->B"

	if err := os.WriteFile(inputFile, []byte(mermaidCode), 0o644); err != nil {
		t.Fatalf("failed to create input file: %v", err)
	}

	var capturedOutput bytes.Buffer
	restore := usecase.SetOsStdoutForTesting(&capturedOutput)
	defer restore()

	mockR := &mockrenderer.MockRenderer{
		RenderFunc: func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
			return nil, fmt.Errorf("render failed")
		},
	}
	frw := file.NewFileReaderWriter()
	mockL := &mockLogger{}

	u := usecase.NewRenderUsecase(mockR, frw, mockL)
	ctx := context.Background()
	cfg := config.NewRenderConfig()

	err := u.RenderFromInput(ctx, inputFile, "-", domain.OutputFormatSVG, cfg, "", true)
	if err == nil {
		t.Fatal("expected error")
	}
}
