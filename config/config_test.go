package config_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/ageha734/mmdc/config"
)

func TestNewRenderConfig(t *testing.T) {
	cfg := config.NewRenderConfig()
	if cfg == nil {
		t.Fatal("NewRenderConfig returned nil")
	}
	if cfg.Theme != "default" {
		t.Errorf("expected Theme 'default', got %q", cfg.Theme)
	}
	if cfg.Width != 800 {
		t.Errorf("expected Width 800, got %d", cfg.Width)
	}
	if cfg.Height != 600 {
		t.Errorf("expected Height 600, got %d", cfg.Height)
	}
	if cfg.BackgroundColor != "white" {
		t.Errorf("expected BackgroundColor 'white', got %q", cfg.BackgroundColor)
	}
	if cfg.Scale != 1 {
		t.Errorf("expected Scale 1, got %d", cfg.Scale)
	}
	if cfg.SVGID != "mermaid-svg" {
		t.Errorf("expected SVGID 'mermaid-svg', got %q", cfg.SVGID)
	}
}

func TestLoadMermaidConfig_NilConfig(t *testing.T) {
	_, err := config.LoadMermaidConfig(nil)
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestLoadMermaidConfig_DefaultConfig(t *testing.T) {
	cfg := config.NewRenderConfig()
	mermaidConfig, err := config.LoadMermaidConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mermaidConfig["theme"] != "default" {
		t.Errorf("expected theme 'default', got %v", mermaidConfig["theme"])
	}
	if mermaidConfig["startOnLoad"] != false {
		t.Errorf("expected startOnLoad false, got %v", mermaidConfig["startOnLoad"])
	}
}

func TestLoadMermaidConfig_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "mermaid-config.json")
	configData := map[string]interface{}{
		"theme": "dark",
		"flowchart": map[string]interface{}{
			"useMaxWidth": true,
		},
	}
	jsonData, _ := json.Marshal(configData)
	if err := os.WriteFile(configFile, jsonData, 0o644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	cfg := config.NewRenderConfig()
	cfg.ConfigFile = configFile
	mermaidConfig, err := config.LoadMermaidConfig(cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mermaidConfig["theme"] != "dark" {
		t.Errorf("expected theme 'dark', got %v", mermaidConfig["theme"])
	}
	if mermaidConfig["flowchart"] == nil {
		t.Error("expected flowchart config to be loaded")
	}
}

func TestLoadMermaidConfig_InvalidConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "invalid-config.json")
	if err := os.WriteFile(configFile, []byte("invalid json"), 0o644); err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}

	cfg := config.NewRenderConfig()
	cfg.ConfigFile = configFile
	_, err := config.LoadMermaidConfig(cfg)
	if err == nil {
		t.Fatal("expected error for invalid config file")
	}
}

func TestLoadMermaidConfig_NonExistentConfigFile(t *testing.T) {
	cfg := config.NewRenderConfig()
	cfg.ConfigFile = "/nonexistent/file.json"
	_, err := config.LoadMermaidConfig(cfg)
	if err == nil {
		t.Fatal("expected error for non-existent config file")
	}
}

func TestDetermineSVGID_WithSVGID(t *testing.T) {
	cfg := config.NewRenderConfig()
	result := config.DetermineSVGID("custom-id", cfg)
	if result != "custom-id" {
		t.Errorf("expected 'custom-id', got %q", result)
	}
}

func TestDetermineSVGID_WithConfigSVGID(t *testing.T) {
	cfg := config.NewRenderConfig()
	cfg.SVGID = "config-id"
	result := config.DetermineSVGID("", cfg)
	if result != "config-id" {
		t.Errorf("expected 'config-id', got %q", result)
	}
}

func TestDetermineSVGID_Default(t *testing.T) {
	cfg := config.NewRenderConfig()
	cfg.SVGID = ""
	result := config.DetermineSVGID("", cfg)
	if result != "mermaid-svg" {
		t.Errorf("expected 'mermaid-svg', got %q", result)
	}
}

func TestConstants(t *testing.T) {
	if config.StdinPath != "-" {
		t.Errorf("expected StdinPath '-', got %q", config.StdinPath)
	}
	if config.DefaultOutputFile != "out.svg" {
		t.Errorf("expected DefaultOutputFile 'out.svg', got %q", config.DefaultOutputFile)
	}
	if config.RenderedAttrTrue != "true" {
		t.Errorf("expected RenderedAttrTrue 'true', got %q", config.RenderedAttrTrue)
	}
	if config.RenderedAttrError != "error" {
		t.Errorf("expected RenderedAttrError 'error', got %q", config.RenderedAttrError)
	}
	if config.FilePermission != 0o644 {
		t.Errorf("expected FilePermission 0644, got %o", config.FilePermission)
	}
	if config.DirPermission != 0o755 {
		t.Errorf("expected DirPermission 0755, got %o", config.DirPermission)
	}
}
