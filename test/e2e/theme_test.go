package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTheme_Default(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-t", "default",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with default theme: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestTheme_Dark(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-t", "dark",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with dark theme: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestTheme_Forest(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-t", "forest",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with forest theme: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestTheme_Neutral(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-t", "neutral",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with neutral theme: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestTheme_AllThemes(t *testing.T) {
	themes := []string{"default", "dark", "forest", "neutral"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.svg")

			_, stderr, err := runMmdc(t,
				"-i", getMockPath("simple.mmd"),
				"-o", outputFile,
				"-t", theme,
				"-q",
			)
			if err != nil {
				t.Fatalf("Failed to render with %s theme: %v\nstderr: %s", theme, err, stderr)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			if !strings.Contains(string(content), "<svg") {
				t.Errorf("Output with %s theme does not contain SVG tag", theme)
			}
		})
	}
}

func TestTheme_ComplexDiagram(t *testing.T) {
	themes := []string{"default", "dark", "forest", "neutral"}

	for _, theme := range themes {
		t.Run(theme, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "complex.svg")

			_, stderr, err := runMmdc(t,
				"-i", getMockPath("complex.mmd"),
				"-o", outputFile,
				"-t", theme,
				"-q",
			)
			if err != nil {
				t.Fatalf("Failed to render complex diagram with %s theme: %v\nstderr: %s", theme, err, stderr)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			if !strings.Contains(string(content), "<svg") {
				t.Errorf("Output with %s theme does not contain SVG tag", theme)
			}
		})
	}
}
