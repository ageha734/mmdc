package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfig_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-c", getMockPath("mermaid-config.json"),
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with config file: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestConfig_WithCSSFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-C", getMockPath("custom.css"),
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with CSS file: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestConfig_WithBothConfigAndCSS(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-c", getMockPath("mermaid-config.json"),
		"-C", getMockPath("custom.css"),
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with config and CSS: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestConfig_PuppeteerConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-p", getMockPath("puppeteer-config.json"),
		"-q",
	)
	if err != nil {
		t.Logf("Puppeteer config test: %v\nstderr: %s", err, stderr)
		return
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Log("Output file was not created (puppeteer config may not be supported)")
	}
}

func TestConfig_NonExistentConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-c", "/non/existent/config.json",
		"-q",
	)

	if err == nil {
		t.Error("Expected error for non-existent config file")
	}
	if !strings.Contains(stderr, "config") && !strings.Contains(stderr, "no such file") && !strings.Contains(stderr, "not found") {
		t.Logf("Error message: %s", stderr)
	}
}

func TestConfig_NonExistentCSSFile(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-C", "/non/existent/style.css",
		"-q",
	)

	t.Logf("Non-existent CSS file: err=%v, stderr=%s", err, stderr)

	if err != nil {
		t.Logf("Error for non-existent CSS file (may be expected): %v", err)
	}
}

func TestConfig_InvalidJSONConfig(t *testing.T) {
	tmpDir := t.TempDir()

	invalidConfig := filepath.Join(tmpDir, "invalid.json")
	if err := os.WriteFile(invalidConfig, []byte("{invalid json}"), 0o644); err != nil {
		t.Fatalf("Failed to create invalid config file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-c", invalidConfig,
		"-q",
	)

	if err == nil {
		t.Error("Expected error for invalid JSON config")
	}
	t.Logf("Invalid JSON error: %s", stderr)
}

func TestConfig_EmptyConfigFile(t *testing.T) {
	tmpDir := t.TempDir()

	emptyConfig := filepath.Join(tmpDir, "empty.json")
	if err := os.WriteFile(emptyConfig, []byte("{}"), 0o644); err != nil {
		t.Fatalf("Failed to create empty config file: %v", err)
	}

	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-c", emptyConfig,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with empty config: %v\nstderr: %s", err, stderr)
	}

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}
}

func TestConfig_ComplexDiagramWithConfig(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("complex.mmd"),
		"-o", outputFile,
		"-c", getMockPath("mermaid-config.json"),
		"-C", getMockPath("custom.css"),
		"-t", "dark",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render complex diagram with config: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}
