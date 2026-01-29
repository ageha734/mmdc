package domain_test

import (
	"testing"

	"github.com/ageha734/mmdc/domain"
)

func TestOutputFormat_String(t *testing.T) {
	tests := []struct {
		name     string
		format   domain.OutputFormat
		expected string
	}{
		{"SVG", domain.OutputFormatSVG, "svg"},
		{"PNG", domain.OutputFormatPNG, "png"},
		{"PDF", domain.OutputFormatPDF, "pdf"},
		{"Empty", domain.OutputFormat(""), ""},
		{"Custom", domain.OutputFormat("custom"), "custom"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.String()
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestOutputFormat_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		format   domain.OutputFormat
		expected bool
	}{
		{"SVG", domain.OutputFormatSVG, true},
		{"PNG", domain.OutputFormatPNG, true},
		{"PDF", domain.OutputFormatPDF, true},
		{"Empty", domain.OutputFormat(""), false},
		{"Invalid", domain.OutputFormat("invalid"), false},
		{"JPG", domain.OutputFormat("jpg"), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.format.IsValid()
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestMermaidDiagram(t *testing.T) {
	diagram := &domain.MermaidDiagram{
		Code:      "graph TD\nA-->B",
		Index:     0,
		FullMatch: "```mermaid\ngraph TD\nA-->B\n```",
	}
	if diagram.Code != "graph TD\nA-->B" {
		t.Errorf("expected Code 'graph TD\\nA-->B', got %q", diagram.Code)
	}
	if diagram.Index != 0 {
		t.Errorf("expected Index 0, got %d", diagram.Index)
	}
	if diagram.FullMatch != "```mermaid\ngraph TD\nA-->B\n```" {
		t.Errorf("unexpected FullMatch: %q", diagram.FullMatch)
	}
}
