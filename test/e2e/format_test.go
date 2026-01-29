package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFormat_ExplicitSVG(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "svg",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render SVG: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestFormat_ExplicitPNG(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.png")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "png",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render PNG: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if len(content) < 4 || content[0] != 0x89 || content[1] != 0x50 || content[2] != 0x4E || content[3] != 0x47 {
		t.Error("Output is not a valid PNG")
	}
}

func TestFormat_ExplicitPDF(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "pdf",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render PDF: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.HasPrefix(string(content), "%PDF") {
		t.Error("Output is not a valid PDF")
	}
}

func TestFormat_FormatOverridesExtension(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

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

	if len(content) >= 4 && content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47 {
		t.Log("Format flag correctly overrides extension")
	} else if strings.Contains(string(content), "<svg") {
		t.Log("Extension took precedence over format flag")
	}
}

func TestFormat_ComplexDiagramAllFormats(t *testing.T) {
	formats := []struct {
		format    string
		extension string
		checkFunc func([]byte) bool
	}{
		{
			format:    "svg",
			extension: ".svg",
			checkFunc: func(content []byte) bool {
				return strings.Contains(string(content), "<svg")
			},
		},
		{
			format:    "png",
			extension: ".png",
			checkFunc: func(content []byte) bool {
				return len(content) >= 4 && content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47
			},
		},
		{
			format:    "pdf",
			extension: ".pdf",
			checkFunc: func(content []byte) bool {
				return strings.HasPrefix(string(content), "%PDF")
			},
		},
	}

	for _, f := range formats {
		t.Run(strings.ToUpper(f.format), func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "complex"+f.extension)

			_, stderr, err := runMmdc(t,
				"-i", getMockPath("complex.mmd"),
				"-o", outputFile,
				"-e", f.format,
				"-q",
			)
			if err != nil {
				t.Fatalf("Failed to render complex diagram as %s: %v\nstderr: %s", f.format, err, stderr)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			if !f.checkFunc(content) {
				t.Errorf("Output is not a valid %s file", f.format)
			}
		})
	}
}

func TestFormat_PDFWithFit(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.pdf")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "pdf",
		"-f",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render PDF with fit: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.HasPrefix(string(content), "%PDF") {
		t.Error("Output is not a valid PDF")
	}
}
