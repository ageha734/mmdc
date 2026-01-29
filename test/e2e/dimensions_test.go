package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDimensions_Width(t *testing.T) {
	widths := []int{400, 800, 1200, 1920}

	for _, width := range widths {
		t.Run(string(rune(width)), func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.svg")

			_, stderr, err := runMmdc(t,
				"-i", getMockPath("simple.mmd"),
				"-o", outputFile,
				"-w", string(rune(width)),
				"-q",
			)
			if err != nil && !strings.Contains(stderr, "invalid") {
				t.Logf("Render with width %d: %v", width, err)
			}
		})
	}
}

func TestDimensions_DefaultWidth(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with default width: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestDimensions_CustomDimensions(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-w", "1200",
		"-H", "800",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with custom dimensions: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestDimensions_Scale(t *testing.T) {
	scales := []string{"1", "2", "3"}

	for _, scale := range scales {
		t.Run("Scale_"+scale, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.svg")

			_, stderr, err := runMmdc(t,
				"-i", getMockPath("simple.mmd"),
				"-o", outputFile,
				"-s", scale,
				"-q",
			)
			if err != nil {
				t.Fatalf("Failed to render with scale %s: %v\nstderr: %s", scale, err, stderr)
			}

			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				t.Fatalf("Output file was not created for scale %s", scale)
			}
		})
	}
}

func TestDimensions_PNGWithScale(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile1 := filepath.Join(tmpDir, "scale1.png")
	outputFile2 := filepath.Join(tmpDir, "scale2.png")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile1,
		"-e", "png",
		"-s", "1",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with scale 1: %v\nstderr: %s", err, stderr)
	}

	_, stderr, err = runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile2,
		"-e", "png",
		"-s", "2",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with scale 2: %v\nstderr: %s", err, stderr)
	}

	stat1, err := os.Stat(outputFile1)
	if err != nil {
		t.Fatal("Scale 1 output file not created")
	}

	stat2, err := os.Stat(outputFile2)
	if err != nil {
		t.Fatal("Scale 2 output file not created")
	}

	t.Logf("Scale 1 file size: %d, Scale 2 file size: %d", stat1.Size(), stat2.Size())
}

func TestDimensions_SmallDimensions(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-w", "200",
		"-H", "150",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with small dimensions: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestDimensions_LargeDimensions(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-w", "3840",
		"-H", "2160",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with large dimensions: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}
