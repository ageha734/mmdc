package e2e

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binaryPath string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "mmdc-e2e-*")
	if err != nil {
		panic("Failed to create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	binaryPath = filepath.Join(tmpDir, "mmdc")
	if os.Getenv("GOOS") == "windows" {
		binaryPath += ".exe"
	}

	cmd := exec.Command("go", "build", "-o", binaryPath, "../../cmd")
	cmd.Dir = getTestDir()
	if output, err := cmd.CombinedOutput(); err != nil {
		panic("Failed to build binary: " + err.Error() + "\n" + string(output))
	}

	os.Exit(m.Run())
}

func getTestDir() string {
	dir, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}
	return dir
}

func getMockPath(filename string) string {
	return filepath.Join(getTestDir(), "mock", filename)
}

func runMmdc(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := exec.Command(binaryPath, args...)
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}

func TestBasicRender_SVG(t *testing.T) {
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

	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Fatal("Output file was not created")
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestBasicRender_PNG(t *testing.T) {
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

	if len(content) < 4 || content[0] != 0x89 || content[1] != 0x50 || content[2] != 0x4E || content[3] != 0x47 {
		t.Error("Output is not a valid PNG file")
	}
}

func TestBasicRender_PDF(t *testing.T) {
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

	if !strings.HasPrefix(string(content), "%PDF") {
		t.Error("Output is not a valid PDF file")
	}
}

func TestBasicRender_AutoDetectFormat(t *testing.T) {
	testCases := []struct {
		name      string
		extension string
		checkFunc func([]byte) bool
	}{
		{
			name:      "SVG",
			extension: ".svg",
			checkFunc: func(content []byte) bool {
				return strings.Contains(string(content), "<svg")
			},
		},
		{
			name:      "PNG",
			extension: ".png",
			checkFunc: func(content []byte) bool {
				return len(content) >= 4 && content[0] == 0x89 && content[1] == 0x50 && content[2] == 0x4E && content[3] == 0x47
			},
		},
		{
			name:      "PDF",
			extension: ".pdf",
			checkFunc: func(content []byte) bool {
				return strings.HasPrefix(string(content), "%PDF")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputFile := filepath.Join(tmpDir, "output"+tc.extension)

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

			if !tc.checkFunc(content) {
				t.Errorf("Output format validation failed for %s", tc.name)
			}
		})
	}
}

func TestBasicRender_ComplexDiagram(t *testing.T) {
	tmpDir := t.TempDir()
	outputFile := filepath.Join(tmpDir, "complex.svg")

	_, stderr, err := runMmdc(t,
		"-i", getMockPath("complex.mmd"),
		"-o", outputFile,
		"-q",
	)
	if err != nil {
		t.Fatalf("Failed to render complex diagram: %v\nstderr: %s", err, stderr)
	}

	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	if !strings.Contains(string(content), "<svg") {
		t.Error("Output does not contain SVG tag")
	}
}

func TestBasicRender_Version(t *testing.T) {
	stdout, _, err := runMmdc(t, "--version")
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if !strings.Contains(stdout, "1.0.0") && !strings.Contains(stdout, "mmdc") {
		t.Error("Version output does not contain expected version string")
	}
}

func TestBasicRender_Help(t *testing.T) {
	stdout, _, err := runMmdc(t, "--help")
	if err != nil {
		t.Fatalf("Failed to get help: %v", err)
	}

	requiredStrings := []string{
		"--input",
		"--output",
		"--outputFormat",
		"--theme",
		"--quiet",
	}

	for _, s := range requiredStrings {
		if !strings.Contains(stdout, s) {
			t.Errorf("Help output does not contain %q", s)
		}
	}
}
