package markdown_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/pkg/markdown"
)

func TestExtractAllMermaidFromMarkdown(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedCount int
	}{
		{
			name:          "コードブロック1つ",
			content:       "```mermaid\ngraph TD\nA-->B\n```",
			expectedCount: 1,
		},
		{
			name:          "フェンス記法1つ",
			content:       ":::mermaid\ngraph TD\nA-->B\n:::",
			expectedCount: 1,
		},
		{
			name:          "コードブロックとフェンス記法",
			content:       "```mermaid\ngraph TD\nA-->B\n```\n:::mermaid\nsequenceDiagram\nA->>B: Hello\n:::",
			expectedCount: 2,
		},
		{
			name:          "Mermaidなし",
			content:       "普通のMarkdownテキスト",
			expectedCount: 0,
		},
		{
			name:          "空文字列",
			content:       "",
			expectedCount: 0,
		},
		{
			name:          "複数のコードブロック",
			content:       "```mermaid\ngraph TD\nA-->B\n```\nテキスト\n```mermaid\nsequenceDiagram\nA->>B: Hello\n```",
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := markdown.ExtractAllMermaidFromMarkdown(tt.content)
			if len(matches) != tt.expectedCount {
				t.Errorf("expected %d matches, got %d", tt.expectedCount, len(matches))
			}
			for i, match := range matches {
				if match.Code == "" {
					t.Errorf("match[%d].Code should not be empty", i)
				}
				if match.FullMatch == "" {
					t.Errorf("match[%d].FullMatch should not be empty", i)
				}
			}
		})
	}
}

func TestGenerateOutputFileName(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		index      int
		format     domain.OutputFormat
		checkFunc  func(*testing.T, string)
	}{
		{
			name:       "Markdownファイル",
			outputPath: "test.md",
			index:      0,
			format:     domain.OutputFormatSVG,
			checkFunc: func(t *testing.T, result string) {
				baseName := filepath.Base(result)
				if baseName != "test-1.svg" {
					t.Errorf("expected basename 'test-1.svg', got %q", baseName)
				}
			},
		},
		{
			name:       "SVGファイル",
			outputPath: "test.svg",
			index:      1,
			format:     domain.OutputFormatPNG,
			checkFunc: func(t *testing.T, result string) {
				baseName := filepath.Base(result)
				if baseName != "test-2.svg" {
					t.Errorf("expected basename 'test-2.svg', got %q", baseName)
				}
			},
		},
		{
			name:       "複数の図",
			outputPath: "test.md",
			index:      2,
			format:     domain.OutputFormatSVG,
			checkFunc: func(t *testing.T, result string) {
				baseName := filepath.Base(result)
				if baseName != "test-3.svg" {
					t.Errorf("expected basename 'test-3.svg', got %q", baseName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := markdown.GenerateOutputFileName(tt.outputPath, tt.index, tt.format)
			tt.checkFunc(t, result)
		})
	}
}

func TestCalculateRelativePath(t *testing.T) {
	tests := []struct {
		name       string
		outputPath string
		outputFile string
		checkFunc  func(*testing.T, string)
	}{
		{
			name:       "同じディレクトリ",
			outputPath: "dir/output.md",
			outputFile: "dir/diagram-1.svg",
			checkFunc: func(t *testing.T, result string) {
				if result != "./diagram-1.svg" {
					t.Errorf("expected './diagram-1.svg', got %q", result)
				}
			},
		},
		{
			name:       "親ディレクトリ",
			outputPath: "dir/sub/output.md",
			outputFile: "dir/diagram-1.svg",
			checkFunc: func(t *testing.T, result string) {
				if result == "" {
					t.Errorf("expected non-empty relative path")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			outputPath := filepath.Join(tmpDir, tt.outputPath)
			outputFile := filepath.Join(tmpDir, tt.outputFile)

			os.MkdirAll(filepath.Dir(outputPath), 0o755)
			os.MkdirAll(filepath.Dir(outputFile), 0o755)

			result := markdown.CalculateRelativePath(outputPath, outputFile)
			tt.checkFunc(t, result)
		})
	}
}

func TestReplaceMermaidWithImages(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		imageInfos []domain.ImageInfo
		checkFunc  func(*testing.T, string)
	}{
		{
			name:    "1つのコードブロック",
			content: "```mermaid\ngraph TD\nA-->B\n```",
			imageInfos: []domain.ImageInfo{
				{URL: "./diagram-1.svg", Alt: "diagram 1"},
			},
			checkFunc: func(t *testing.T, result string) {
				if !contains(result, "![diagram 1](./diagram-1.svg)") {
					t.Errorf("expected image link, got %q", result)
				}
				if contains(result, "```mermaid") {
					t.Errorf("mermaid code block should be replaced")
				}
			},
		},
		{
			name:    "複数のコードブロック",
			content: "```mermaid\ngraph TD\nA-->B\n```\nテキスト\n```mermaid\nsequenceDiagram\nA->>B: Hello\n```",
			imageInfos: []domain.ImageInfo{
				{URL: "./diagram-1.svg", Alt: "diagram 1"},
				{URL: "./diagram-2.svg", Alt: "diagram 2"},
			},
			checkFunc: func(t *testing.T, result string) {
				if !contains(result, "![diagram 1](./diagram-1.svg)") {
					t.Errorf("expected first image link")
				}
				if !contains(result, "![diagram 2](./diagram-2.svg)") {
					t.Errorf("expected second image link")
				}
			},
		},
		{
			name:    "フェンス記法",
			content: ":::mermaid\ngraph TD\nA-->B\n:::",
			imageInfos: []domain.ImageInfo{
				{URL: "./diagram-1.svg", Alt: "diagram 1"},
			},
			checkFunc: func(t *testing.T, result string) {
				if !contains(result, "![diagram 1](./diagram-1.svg)") {
					t.Errorf("expected image link for fence syntax")
				}
			},
		},
		{
			name:       "Mermaidなし",
			content:    "普通のテキスト",
			imageInfos: []domain.ImageInfo{},
			checkFunc: func(t *testing.T, result string) {
				if result != "普通のテキスト" {
					t.Errorf("content should not change: %q", result)
				}
			},
		},
		{
			name:    "imageInfosが足りない場合",
			content: "```mermaid\ngraph TD\nA-->B\n```\nテキスト\n```mermaid\nsequenceDiagram\nA->>B: Hello\n```",
			imageInfos: []domain.ImageInfo{
				{URL: "./diagram-1.svg", Alt: "diagram 1"},
			},
			checkFunc: func(t *testing.T, result string) {
				if !contains(result, "![diagram 1](./diagram-1.svg)") {
					t.Errorf("expected first image link")
				}
				if !contains(result, "```mermaid") {
					t.Errorf("second mermaid block should remain unchanged")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := markdown.ReplaceMermaidWithImages(tt.content, tt.imageInfos)
			tt.checkFunc(t, result)
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
