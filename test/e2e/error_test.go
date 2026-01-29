package e2e

import (
	"os"
	"path/filepath"
	"testing"
)

func TestError_NoInputFile(t *testing.T) {
	_, _, err := runMmdc(t)

	t.Logf("No input file result: %v", err)
}

func TestError_NonExistentInputFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", "/non/existent/file.mmd",
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Error("Expected error for non-existent input file")
	}
	t.Logf("Non-existent file error: %s", stderr)
}

func TestError_InvalidMermaidSyntax(t *testing.T) {
	tmpDir := t.TempDir()

	invalidFile := filepath.Join(tmpDir, "invalid.mmd")
	if err := os.WriteFile(invalidFile, []byte("invalid mermaid syntax %%%"), 0o644); err != nil {
		t.Fatalf("Failed to create invalid file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", invalidFile,
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Error("Expected error for invalid mermaid syntax")
	}
	t.Logf("Invalid syntax error: %s", stderr)
}

func TestError_InvalidOutputDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "nonexistent", "subdir", "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Error("Expected error for invalid output directory")
	}
	t.Logf("Invalid output directory error: %s", stderr)
}

func TestError_InvalidOutputFormat(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.xyz")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-e", "invalid_format",
		"-q",
	)

	if err == nil {
		t.Error("Expected error for invalid output format")
	}
	t.Logf("Invalid format error: %s", stderr)
}

func TestError_EmptyInputFile(t *testing.T) {
	tmpDir := t.TempDir()

	emptyFile := filepath.Join(tmpDir, "empty.mmd")
	if err := os.WriteFile(emptyFile, []byte(""), 0o644); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", emptyFile,
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Log("Empty file rendered without error (may be valid)")
	} else {
		t.Logf("Empty file error: %s", stderr)
	}
}

func TestError_WhitespaceOnlyInput(t *testing.T) {
	tmpDir := t.TempDir()

	whitespaceFile := filepath.Join(tmpDir, "whitespace.mmd")
	if err := os.WriteFile(whitespaceFile, []byte("   \n\t\n  "), 0o644); err != nil {
		t.Fatalf("Failed to create whitespace file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", whitespaceFile,
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Log("Whitespace-only file rendered without error (may be valid)")
	} else {
		t.Logf("Whitespace file error: %s", stderr)
	}
}

func TestError_ReadOnlyOutputDirectory(t *testing.T) {
	if os.Getuid() == 0 {
		t.Skip("Skipping test as root user")
	}

	tmpDir := t.TempDir()

	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0o555); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	defer os.Chmod(readOnlyDir, 0o755)

	outputFile := filepath.Join(readOnlyDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Error("Expected error for read-only output directory")
	}
	t.Logf("Read-only directory error: %s", stderr)
}

func TestError_InvalidNegativeWidth(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-w", "-100",
		"-q",
	)

	if err == nil {
		t.Log("Negative width rendered without error (implementation dependent)")
	} else {
		t.Logf("Negative width error: %s", stderr)
	}
}

func TestError_InvalidNegativeHeight(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-H", "-100",
		"-q",
	)

	if err == nil {
		t.Log("Negative height rendered without error (implementation dependent)")
	} else {
		t.Logf("Negative height error: %s", stderr)
	}
}

func TestError_InvalidScale(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-s", "-1",
		"-q",
	)

	if err == nil {
		t.Log("Negative scale rendered without error (implementation dependent)")
	} else {
		t.Logf("Negative scale error: %s", stderr)
	}
}

func TestError_MarkdownWithInvalidDiagram(t *testing.T) {
	tmpDir := t.TempDir()

	invalidMd := filepath.Join(tmpDir, "invalid.md")
	content := "# Test\n\n```mermaid\ninvalid mermaid %%%\n```\n"
	if err := os.WriteFile(invalidMd, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create invalid markdown: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.md")

	_, stderr, err := runMmdc(t,
		"-i", invalidMd,
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Error("Expected error for markdown with invalid mermaid")
	}
	t.Logf("Invalid markdown diagram error: %s", stderr)
}

func TestError_BinaryInputFile(t *testing.T) {
	tmpDir := t.TempDir()

	binaryFile := filepath.Join(tmpDir, "binary.mmd")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03, 0xFF, 0xFE, 0xFD}
	if err := os.WriteFile(binaryFile, binaryContent, 0o644); err != nil {
		t.Fatalf("Failed to create binary file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", binaryFile,
		"-o", outputFile,
		"-q",
	)

	if err == nil {
		t.Log("Binary file rendered without error (unexpected)")
	} else {
		t.Logf("Binary file error: %s", stderr)
	}
}

func TestError_VeryLargeInputFile(t *testing.T) {
	tmpDir := t.TempDir()

	largeFile := filepath.Join(tmpDir, "large.mmd")

	content := "graph TD\n"
	for i := 0; i < 100; i++ {
		content += "    A" + string(rune(i)) + " --> B" + string(rune(i)) + "\n"
	}
	if err := os.WriteFile(largeFile, []byte(content), 0o644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", largeFile,
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Logf("Large file error (might be expected): %s", stderr)
	}
}
