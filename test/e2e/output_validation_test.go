package e2e

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type SVGElement struct {
	XMLName xml.Name
	ID      string `xml:"id,attr"`
	Width   string `xml:"width,attr"`
	Height  string `xml:"height,attr"`
	ViewBox string `xml:"viewBox,attr"`
	Inner   []byte `xml:",innerxml"`
}

func TestOutputValidation_SVGStructure(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	var svg SVGElement
	if err := xml.Unmarshal(content, &svg); err != nil {
		t.Fatalf("Output is not valid XML: %v", err)
	}

	if svg.XMLName.Local != "svg" {
		t.Errorf("Expected root element 'svg', got '%s'", svg.XMLName.Local)
	}
}

func TestOutputValidation_SVGContainsExpectedElements(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	contentStr := string(content)

	expectedElements := []string{
		"<svg",
		"</svg>",
	}

	for _, elem := range expectedElements {
		if !strings.Contains(contentStr, elem) {
			t.Errorf("SVG output should contain %q", elem)
		}
	}
}

func TestOutputValidation_SVGWithCustomId(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")
	customId := "my-test-diagram"

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-I", customId,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), customId) {
		t.Error("SVG output should contain custom ID")
	}
}

func TestOutputValidation_PNGMagicBytes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.png")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "png",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	expectedMagic := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
	if len(content) < 8 {
		t.Fatal("PNG file is too small")
	}

	for i, b := range expectedMagic {
		if content[i] != b {
			t.Errorf("PNG magic byte %d: expected 0x%02X, got 0x%02X", i, b, content[i])
		}
	}
}

func TestOutputValidation_PNGMinimumSize(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.png")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "png",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to stat output: %v", err)
	}

	if stat.Size() < 100 {
		t.Errorf("PNG file seems too small: %d bytes", stat.Size())
	}
}

func TestOutputValidation_PDFMagicBytes(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "pdf",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.HasPrefix(string(content), "%PDF-") {
		t.Error("PDF file does not start with %PDF-")
	}
}

func TestOutputValidation_PDFMinimumSize(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "pdf",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	stat, err := os.Stat(outputFile)
	if err != nil {
		t.Fatalf("Failed to stat output: %v", err)
	}

	if stat.Size() < 100 {
		t.Errorf("PDF file seems too small: %d bytes", stat.Size())
	}
}

func TestOutputValidation_PDFContainsEOFMarker(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "pdf",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "%%EOF") {
		t.Errorf("%s", "PDF file does not contain %%EOF marker")
	}
}

func TestOutputValidation_ComplexDiagramSVG(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "complex.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("complex.mmd"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "<g") {
		t.Error("Complex SVG should contain <g> (group) elements")
	}
}

func TestOutputValidation_MarkdownOutput(t *testing.T) {
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
		t.Fatalf("Failed to render: %v\nstderr: %s", err, stderr)
	}

	mdContent, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read markdown output: %v", err)
	}

	mdStr := string(mdContent)

	if strings.Contains(mdStr, "```mermaid") {
		t.Error("Markdown output should not contain mermaid code blocks")
	}

	if !strings.Contains(mdStr, "![") {
		t.Error("Markdown output should contain image references")
	}

	files, err := os.ReadDir(artefactsDir)
	if err != nil {
		t.Fatalf("Failed to read artefacts directory: %v", err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".svg") {
			svgPath := filepath.Join(artefactsDir, f.Name())
			svgContent, err := os.ReadFile(svgPath)
			if err != nil {
				t.Fatalf("Failed to read SVG file %s: %v", f.Name(), err)
			}

			if !strings.Contains(string(svgContent), "<svg") {
				t.Errorf("SVG file %s does not contain SVG tag", f.Name())
			}
		}
	}
}

func TestOutputValidation_FileSizeConsistency(t *testing.T) {
	tmpDir := t.TempDir()

	sizes := make([]int64, 3)
	for i := 0; i < 3; i++ {
		outputFile := filepath.Join(tmpDir, "output"+string(rune('A'+i))+".svg")

		_, stderr, err := runMmdc(t,
			"-i", getMockPath("simple.mmd"),
			"-o", outputFile,
			"-q",
		)
		if err != nil {
			t.Fatalf("Failed to render iteration %d: %v\nstderr: %s", i, err, stderr)
		}

		stat, err := os.Stat(outputFile)
		if err != nil {
			t.Fatalf("Failed to stat output: %v", err)
		}
		sizes[i] = stat.Size()
	}

	for i, size := range sizes {
		if size == 0 {
			t.Errorf("Output %d has zero size", i)
		}
	}

	avgSize := (sizes[0] + sizes[1] + sizes[2]) / 3
	for i, size := range sizes {
		diff := size - avgSize
		if diff < 0 {
			diff = -diff
		}
		if float64(diff) > float64(avgSize)*0.1 {
			t.Errorf("Output %d size %d varies significantly from average %d", i, size, avgSize)
		}
	}
}
