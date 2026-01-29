package path_test

import (
	"testing"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/pkg/path"
)

func TestIsMarkdownFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{"Markdown .md", "test.md", true},
		{"Markdown .markdown", "test.markdown", true},
		{"大文字MD", "test.MD", false},
		{"stdin", config.StdinPath, false},
		{"テキストファイル", "test.txt", false},
		{"SVG", "test.svg", false},
		{"空文字列", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := path.IsMarkdownFile(tt.path)
			if result != tt.expected {
				t.Errorf("IsMarkdownFile(%q) = %v, expected %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestDetermineOutputPath(t *testing.T) {
	tests := []struct {
		name       string
		inputPath  string
		outputPath string
		expected   string
	}{
		{"出力パス指定あり", "input.mmd", "output.svg", "output.svg"},
		{"出力パス未指定", "input.mmd", "", "input.mmd.svg"},
		{"stdin入力", config.StdinPath, "", config.StdinPath},
		{"stdin入力・出力パス指定", config.StdinPath, "output.svg", "output.svg"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := path.DetermineOutputPath(tt.inputPath, tt.outputPath)
			if result != tt.expected {
				t.Errorf("DetermineOutputPath(%q, %q) = %q, expected %q",
					tt.inputPath, tt.outputPath, result, tt.expected)
			}
		})
	}
}

func TestDetermineOutputFormatFromPath(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		expected   domain.OutputFormat
	}{
		{"SVG", "test.svg", domain.OutputFormatSVG},
		{"PNG", "test.png", domain.OutputFormatPNG},
		{"PDF", "test.pdf", domain.OutputFormatPDF},
		{"大文字SVG", "test.SVG", domain.OutputFormatSVG},
		{"拡張子なし", "test", domain.OutputFormatSVG},
		{"無効な拡張子", "test.txt", domain.OutputFormatSVG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := path.DetermineOutputFormatFromPath(tt.outputPath)
			if result != tt.expected {
				t.Errorf("DetermineOutputFormatFromPath(%q) = %v, expected %v",
					tt.outputPath, result, tt.expected)
			}
		})
	}
}

func TestDetermineOutputFormatFromExtension(t *testing.T) {
	tests := []struct {
		name     string
		ext      string
		expected domain.OutputFormat
	}{
		{"SVG", "svg", domain.OutputFormatSVG},
		{"PNG", "png", domain.OutputFormatPNG},
		{"PDF", "pdf", domain.OutputFormatPDF},
		{"ドット付きSVG", ".svg", domain.OutputFormatSVG},
		{"ドット付きPNG", ".png", domain.OutputFormatPNG},
		{"無効な拡張子", "txt", domain.OutputFormatSVG},
		{"空文字列", "", domain.OutputFormatSVG},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := path.DetermineOutputFormatFromExtension(tt.ext)
			if result != tt.expected {
				t.Errorf("DetermineOutputFormatFromExtension(%q) = %v, expected %v",
					tt.ext, result, tt.expected)
			}
		})
	}
}
