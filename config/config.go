// Package config provides configuration types and utilities for mmdc.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const (
	RenderTimeout        = 120 * time.Second
	PollingTimeout       = 90 * time.Second
	PollingInterval      = 500 * time.Millisecond
	ScriptLoadWaitTime   = 3 * time.Second
	FinalWaitTime        = 1 * time.Second
	PollingCheckInterval = 10
	TempFilePattern      = "mmdc-*.html"
	UserDataDirPattern   = "mmdc-chrome-*"
	RenderedAttrTrue     = "true"
	RenderedAttrError    = "error"
	StdinPath            = "-"
	DefaultOutputFile    = "out.svg"
	FilePermission       = 0o644
	DirPermission        = 0o755
)

type RenderConfig struct {
	Theme           string
	Width           int
	Height          int
	BackgroundColor string
	Scale           int
	ConfigFile      string
	CSSFile         string
	SVGID           string
	PDFFit          bool
}

func NewRenderConfig() *RenderConfig {
	return &RenderConfig{
		Theme:           "default",
		Width:           800,
		Height:          600,
		BackgroundColor: "white",
		Scale:           1,
		SVGID:           "mermaid-svg",
	}
}

func LoadMermaidConfig(renderConfig *RenderConfig) (map[string]interface{}, error) {
	if renderConfig == nil {
		return nil, fmt.Errorf("renderConfig is nil")
	}
	mermaidConfig := map[string]interface{}{
		"theme":       renderConfig.Theme,
		"startOnLoad": false,
	}

	if renderConfig.ConfigFile != "" {
		configData, err := os.ReadFile(renderConfig.ConfigFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		if err := json.Unmarshal(configData, &mermaidConfig); err != nil {
			return nil, fmt.Errorf("failed to parse config file: %w", err)
		}
	}

	return mermaidConfig, nil
}

func DetermineSVGID(svgId string, renderConfig *RenderConfig) string {
	if svgId != "" {
		return svgId
	}
	if renderConfig.SVGID != "" {
		return renderConfig.SVGID
	}
	return "mermaid-svg"
}
