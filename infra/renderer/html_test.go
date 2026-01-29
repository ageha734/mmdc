package renderer_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/infra/renderer"
)

func TestNewHTMLGenerator(t *testing.T) {
	cfg := config.NewRenderConfig()
	gen := renderer.NewHTMLGenerator(cfg)
	if gen == nil {
		t.Fatal("NewHTMLGenerator returned nil")
	}
}

func TestHTMLGenerator_Generate(t *testing.T) {
	cfg := config.NewRenderConfig()
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA-->B", "test-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if html == "" {
		t.Fatal("expected HTML to be generated")
	}
	if !strings.Contains(html, "graph TD") {
		t.Error("HTML should contain mermaid code")
	}
	if !strings.Contains(html, "test-id") {
		t.Error("HTML should contain SVG ID")
	}
	if !strings.Contains(html, "<html>") {
		t.Error("HTML should be valid HTML")
	}
}

func TestHTMLGenerator_Generate_WithCSSFile(t *testing.T) {
	tmpDir := t.TempDir()
	cssFile := filepath.Join(tmpDir, "style.css")
	cssContent := "body { background: red; }"
	if err := os.WriteFile(cssFile, []byte(cssContent), 0o644); err != nil {
		t.Fatalf("failed to create CSS file: %v", err)
	}

	cfg := config.NewRenderConfig()
	cfg.CSSFile = cssFile
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA-->B", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, cssContent) {
		t.Error("HTML should contain CSS content")
	}
}

func TestHTMLGenerator_Generate_WithNonExistentCSSFile(t *testing.T) {
	cfg := config.NewRenderConfig()
	cfg.CSSFile = "/nonexistent/style.css"
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA-->B", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if html == "" {
		t.Fatal("expected HTML to be generated even with non-existent CSS file")
	}
}

func TestHTMLGenerator_Generate_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "mermaid-config.json")
	configData := `{"theme": "dark"}`
	if err := os.WriteFile(configFile, []byte(configData), 0o644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	cfg := config.NewRenderConfig()
	cfg.ConfigFile = configFile
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA-->B", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "dark") {
		t.Error("HTML should contain config theme")
	}
}

func TestHTMLGenerator_Generate_WithBackgroundColor(t *testing.T) {
	cfg := config.NewRenderConfig()
	cfg.BackgroundColor = "black"
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA-->B", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(html, "black") {
		t.Error("HTML should contain background color")
	}
}

func TestHTMLGenerator_Generate_EmptySVGID(t *testing.T) {
	cfg := config.NewRenderConfig()
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA-->B", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if html == "" {
		t.Fatal("expected HTML to be generated")
	}
}

func TestHTMLGenerator_Generate_WithSpecialCharacters(t *testing.T) {
	cfg := config.NewRenderConfig()
	gen := renderer.NewHTMLGenerator(cfg)

	html, err := gen.Generate("graph TD\nA[\"Test <>&'\\\"\"]\nA-->B", "svg-id")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if html == "" {
		t.Fatal("expected HTML to be generated with special characters escaped")
	}
}

func TestHTMLGenerator_Generate_WithNonExistentConfigFile(t *testing.T) {
	cfg := config.NewRenderConfig()
	cfg.ConfigFile = "/nonexistent/config.json"
	gen := renderer.NewHTMLGenerator(cfg)

	_, err := gen.Generate("graph TD\nA-->B", "")
	if err == nil {
		t.Fatal("expected error for non-existent config file")
	}
}

func TestHTMLGenerator_Generate_WithInvalidConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid-config.json")
	if err := os.WriteFile(configFile, []byte("not valid json"), 0o644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	cfg := config.NewRenderConfig()
	cfg.ConfigFile = configFile
	gen := renderer.NewHTMLGenerator(cfg)

	_, err := gen.Generate("graph TD\nA-->B", "")
	if err == nil {
		t.Fatal("expected error for invalid config file")
	}
}
