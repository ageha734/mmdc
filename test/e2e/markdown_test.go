package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMarkdown_ProcessWithDiagrams(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.md")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("with_diagrams.md"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to process markdown: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if strings.Contains(string(content), "```mermaid") {
		t.Error("Output should not contain mermaid code blocks")
	}

	if !strings.Contains(string(content), "![") {
		t.Error("Output should contain image references")
	}
}

func TestMarkdown_WithArtefactsDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.md")
	artefactsDir := filepath.Join(tmpDir, "images")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("with_diagrams.md"),
		"-o", outputFile,
		"-a", artefactsDir,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to process markdown: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(artefactsDir); os.IsNotExist(err) {
		t.Fatal("Artefacts directory was not created")
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

	if svgCount < 1 {
		t.Error("Expected SVG files in artefacts directory")
	}
}

func TestMarkdown_WithPNGFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.md")
	artefactsDir := filepath.Join(tmpDir, "images")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("with_diagrams.md"),
		"-o", outputFile,
		"-a", artefactsDir,
		"-e", "png",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to process markdown: %v\nstderr: %s", err, stderr)
	}

	files, err := os.ReadDir(artefactsDir)
	if err != nil {
		t.Fatalf("Failed to read artefacts directory: %v", err)
	}

	pngCount := 0
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".png") {
			pngCount++
		}
	}

	if pngCount < 1 {
		t.Error("Expected PNG files in artefacts directory")
	}
}

func TestMarkdown_PreservesNonMermaidContent(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.md")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("with_diagrams.md"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to process markdown: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "# Test Document with Mermaid Diagrams") {
		t.Error("Title should be preserved")
	}
	if !strings.Contains(contentStr, "## Flow Chart") {
		t.Error("Section header should be preserved")
	}
	if !strings.Contains(contentStr, "End of document.") {
		t.Error("End text should be preserved")
	}
}

func TestMarkdown_MultipleDiagrams(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.md")
	artefactsDir := filepath.Join(tmpDir, "images")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("with_diagrams.md"),
		"-o", outputFile,
		"-a", artefactsDir,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to process markdown: %v\nstderr: %s", err, stderr)
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

func TestMarkdown_ImageReferencesAreRelative(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.md")
	artefactsDir := filepath.Join(tmpDir, "images")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("with_diagrams.md"),
		"-o", outputFile,
		"-a", artefactsDir,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to process markdown: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	contentStr := string(content)

	if strings.Contains(contentStr, tmpDir) {
		t.Error("Image references should be relative, not absolute")
	}

	if !strings.Contains(contentStr, "images/") && !strings.Contains(contentStr, "./") {
		t.Error("Image references should be relative paths")
	}
}
