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

func TestRenderUsecase_RenderFromInputToFile(t *testing.T) {
	tmpDir := t.TempDir()

	inputContent := `graph TD
    A[Start] --> B[Process]
    B --> C[End]`

	inputFile := filepath.Join(tmpDir, "input.mmd")
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	uc := usecase.NewRenderUsecase(r, fileRW, log)
	err := uc.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, renderConfig, "", true)
	if err != nil {
		t.Fatalf("Failed to render from input: %v", err)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output file does not contain SVG content")
	}
}

func TestRenderUsecase_RenderToPNG(t *testing.T) {
	tmpDir := t.TempDir()

	inputContent := `graph LR
    A --> B --> C`

	inputFile := filepath.Join(tmpDir, "input.mmd")
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.png")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	uc := usecase.NewRenderUsecase(r, fileRW, log)
	err := uc.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatPNG, renderConfig, "", true)
	if err != nil {
		t.Fatalf("Failed to render to PNG: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if len(content) < 4 || content[0] != 0x89 || content[1] != 0x50 || content[2] != 0x4E || content[3] != 0x47 {
		t.Error("Output file is not a valid PNG")
	}
}

func TestRenderUsecase_RenderToPDF(t *testing.T) {
	tmpDir := t.TempDir()

	inputContent := `sequenceDiagram
    Alice->>Bob: Hello
    Bob->>Alice: Hi`

	inputFile := filepath.Join(tmpDir, "input.mmd")
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.pdf")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	uc := usecase.NewRenderUsecase(r, fileRW, log)
	err := uc.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatPDF, renderConfig, "", true)
	if err != nil {
		t.Fatalf("Failed to render to PDF: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if !strings.HasPrefix(string(content), "%PDF") {
		t.Error("Output file is not a valid PDF")
	}
}

func TestRenderUsecase_RenderMarkdownFile(t *testing.T) {
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

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	uc := usecase.NewRenderUsecase(r, fileRW, log)
	err := uc.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, renderConfig, "", true)
	if err != nil {
		t.Fatalf("Failed to render markdown file: %v", err)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	if strings.Contains(string(content), "```mermaid") {
		t.Error("Output should not contain mermaid code blocks")
	}
	if !strings.Contains(string(content), "![") {
		t.Error("Output should contain image references")
	}
}

func TestRenderUsecase_OutputFormatFromExtension(t *testing.T) {
	testCases := []struct {
		name      string
		extension string
		format    domain.OutputFormat
	}{
		{"SVG", ".svg", domain.OutputFormatSVG},
		{"PNG", ".png", domain.OutputFormatPNG},
		{"PDF", ".pdf", domain.OutputFormatPDF},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			inputContent := `graph TD
    A --> B`

			inputFile := filepath.Join(tmpDir, "input.mmd")
			if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
				t.Fatalf("Failed to write input file: %v", err)
			}

			outputFile := filepath.Join(tmpDir, "output"+tc.extension)

			ctx := context.Background()
			renderConfig := config.NewRenderConfig()
			opts := renderer.RenderOptions{Quiet: true}
			r := renderer.NewChromedpRenderer(opts)
			fileRW := file.NewFileReaderWriter()
			log := mockLogger.NewMockLogger()

			uc := usecase.NewRenderUsecase(r, fileRW, log)
			err := uc.RenderFromInput(ctx, inputFile, outputFile, "", renderConfig, "", true)
			if err != nil {
				t.Fatalf("Failed to render: %v", err)
			}

			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				t.Fatal("Output file was not created")
			}
		})
	}
}

func TestFileReaderWriter_ReadWrite(t *testing.T) {
	tmpDir := t.TempDir()

	testContent := "test content for file I/O"
	testFile := filepath.Join(tmpDir, "test.txt")

	fileRW := file.NewFileReaderWriter()

	if err := fileRW.Write(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	content, err := fileRW.Read(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if string(content) != testContent {
		t.Errorf("Content mismatch: expected %q, got %q", testContent, string(content))
	}
}

func TestFileReaderWriter_ReadString(t *testing.T) {
	tmpDir := t.TempDir()

	testContent := "test string content"
	testFile := filepath.Join(tmpDir, "test.txt")

	if err := os.WriteFile(testFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	fileRW := file.NewFileReaderWriter()

	content, err := fileRW.ReadString(testFile)
	if err != nil {
		t.Fatalf("Failed to read string: %v", err)
	}

	if content != testContent {
		t.Errorf("Content mismatch: expected %q, got %q", testContent, content)
	}
}

func TestFileReaderWriter_MkdirAll(t *testing.T) {
	tmpDir := t.TempDir()

	nestedDir := filepath.Join(tmpDir, "a", "b", "c")

	fileRW := file.NewFileReaderWriter()

	if err := fileRW.MkdirAll(nestedDir, 0755); err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	if _, err := os.Stat(nestedDir); os.IsNotExist(err) {
		t.Error("Nested directory was not created")
	}
}

func TestRenderUsecase_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	configContent := `{
		"theme": "dark",
		"flowchart": {
			"curve": "basis"
		}
	}`
	configFile := filepath.Join(tmpDir, "mermaid-config.json")
	if err := os.WriteFile(configFile, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	inputContent := `graph TD
    A --> B`
	inputFile := filepath.Join(tmpDir, "input.mmd")
	if err := os.WriteFile(inputFile, []byte(inputContent), 0644); err != nil {
		t.Fatalf("Failed to write input file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	renderConfig.ConfigFile = configFile
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()
	log := mockLogger.NewMockLogger()

	uc := usecase.NewRenderUsecase(r, fileRW, log)
	err := uc.RenderFromInput(ctx, inputFile, outputFile, domain.OutputFormatSVG, renderConfig, "", true)
	if err != nil {
		t.Fatalf("Failed to render with config file: %v", err)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestRenderUsecase_ReadNonExistentFile(t *testing.T) {
	ctx := context.Background()
	renderConfig := config.NewRenderConfig()
	opts := renderer.RenderOptions{Quiet: true}
	r := renderer.NewChromedpRenderer(opts)
	fileRW := file.NewFileReaderWriter()

	var errorExitCalled bool
	log := &mockLogger.MockLogger{
		ErrorExitFunc: func(msg string) {
			errorExitCalled = true
		},
	}

	uc := usecase.NewRenderUsecase(r, fileRW, log)
	err := uc.RenderFromInput(ctx, "/non/existent/file.mmd", "output.svg", domain.OutputFormatSVG, renderConfig, "", true)

	if err == nil && !errorExitCalled {
		t.Error("Expected error for non-existent file")
	}
}
