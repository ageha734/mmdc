package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestAdvanced_SVGId(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-I", "my-custom-svg-id",
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with custom SVG ID: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "my-custom-svg-id") {
		t.Error("Output should contain custom SVG ID")
	}
}

func TestAdvanced_BackgroundColor(t *testing.T) {
	colors := []string{"white", "black", "#f0f0f0", "#000000", "transparent"}

	for _, color := range colors {
		t.Run(color, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output.svg")

			_, stderr, err := runMmdc(t,
				"-i", getMockPath("simple.mmd"),
				"-o", outputFile,
				"-b", color,
				"-q",
			)
			if err != nil {
				t.Fatalf("Failed to render with background color %s: %v\nstderr: %s", color, err, stderr)
			}

			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				t.Fatalf("Output file was not created for color %s", color)
			}
		})
	}
}

func TestAdvanced_QuietMode(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	stdout, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render in quiet mode: %v\nstderr: %s", err, stderr)
	}

	if strings.Contains(stdout, "Generating") || strings.Contains(stdout, "mermaid") {
		t.Error("Quiet mode should suppress log output")
	}
}

func TestAdvanced_VerboseMode(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	stdout, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
	)
	if err != nil {
		t.Fatalf("Failed to render in verbose mode: %v\nstderr: %s", err, stderr)
	}

	t.Logf("Verbose mode stdout: %s", stdout)
	t.Logf("Verbose mode stderr: %s", stderr)
}

func TestAdvanced_PDFFit(t *testing.T) {
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

func TestAdvanced_CombinedOptions(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("complex.mmd"),
		"-o", outputFile,
		"-t", "dark",
		"-w", "1200",
		"-H", "800",
		"-s", "2",
		"-b", "#1a1a1a",
		"-I", "complex-diagram",
		"-c", getMockPath("mermaid-config.json"),
		"-C", getMockPath("custom.css"),
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with combined options: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "<svg") {
		t.Error("Output does not contain SVG tag")
	}
	if !strings.Contains(contentStr, "complex-diagram") {
		t.Error("Output should contain custom SVG ID")
	}
}

func TestAdvanced_LongSVGId(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	longID := "this-is-a-very-long-svg-id-that-tests-the-limits-of-the-id-parameter"

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-I", longID,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with long SVG ID: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), longID) {
		t.Error("Output should contain the long SVG ID")
	}
}

func TestAdvanced_SpecialCharactersInSVGId(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "output.svg")

	validID := "diagram-123-test"

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("simple.mmd"),
		"-o", outputFile,
		"-I", validID,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render with special SVG ID: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), validID) {
		t.Error("Output should contain the SVG ID with special characters")
	}
}

func TestAdvanced_AllDiagramTypes(t *testing.T) {
	diagramTypes := []struct {
		name string
		code string
	}{
		{
			name: "flowchart",
			code: "graph TD\n    A --> B",
		},
		{
			name: "sequence",
			code: "sequenceDiagram\n    A->>B: Hello",
		},
		{
			name: "class",
			code: "classDiagram\n    class Animal",
		},
		{
			name: "state",
			code: "stateDiagram-v2\n    [*] --> State1",
		},
		{
			name: "er",
			code: "erDiagram\n    CUSTOMER ||--o{ ORDER : places",
		},
		{
			name: "gantt",
			code: "gantt\n    title A Gantt\n    dateFormat YYYY-MM-DD\n    section A\n    Task1 :2024-01-01, 1d",
		},
		{
			name: "pie",
			code: "pie\n    title Pets\n    \"Dogs\" : 50\n    \"Cats\" : 50",
		},
		{
			name: "gitgraph",
			code: "gitGraph\n    commit\n    commit",
		},
	}

	for _, dt := range diagramTypes {
		t.Run(dt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			inputFile := filepath.Join(tmpDir, "input.mmd")
			outputFile := filepath.Join(tmpDir, "output.svg")

			if err := os.WriteFile(inputFile, []byte(dt.code), 0o644); err != nil {
				t.Fatalf("Failed to write input file: %v", err)
			}

			_, stderr, err := runMmdc(t,
				"-i", inputFile,
				"-o", outputFile,
				"-q",
			)
			if err != nil {
				t.Fatalf("Failed to render %s diagram: %v\nstderr: %s", dt.name, err, stderr)
			}

			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output: %v", err)
			}

			if !strings.Contains(string(content), "<svg") {
				t.Errorf("%s diagram output does not contain SVG tag", dt.name)
			}
		})
	}
}
