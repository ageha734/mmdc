package validation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/pkg/validation"
)

func TestValidateInputFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
		setup   func() string
	}{
		{
			name:    "空文字列",
			path:    "",
			wantErr: false,
		},
		{
			name:    "stdinパス",
			path:    "-",
			wantErr: false,
		},
		{
			name:    "存在するファイル",
			wantErr: false,
			setup: func() string {
				tmpFile := filepath.Join(t.TempDir(), "test.txt")
				os.WriteFile(tmpFile, []byte("test"), 0o644)
				return tmpFile
			},
		},
		{
			name:    "存在しないファイル",
			path:    "/nonexistent/file.txt",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.path
			if tt.setup != nil {
				path = tt.setup()
			}
			err := validation.ValidateInputFile(path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInputFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputFileExtension(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"SVG", "test.svg", false},
		{"PNG", "test.png", false},
		{"PDF", "test.pdf", false},
		{"Markdown", "test.md", false},
		{"Markdown alt", "test.markdown", false},
		{"大文字SVG", "test.SVG", false},
		{"無効な拡張子", "test.txt", true},
		{"拡張子なし", "test", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateOutputFileExtension(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputFileExtension() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputFormat(t *testing.T) {
	tests := []struct {
		name    string
		format  domain.OutputFormat
		wantErr bool
	}{
		{"SVG", domain.OutputFormatSVG, false},
		{"PNG", domain.OutputFormatPNG, false},
		{"PDF", domain.OutputFormatPDF, false},
		{"無効", domain.OutputFormat("invalid"), true},
		{"空", domain.OutputFormat(""), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateOutputFormat(tt.format)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateOutputDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
		setup   func() string
	}{
		{
			name:    "カレントディレクトリ",
			dir:     ".",
			wantErr: false,
		},
		{
			name:    "空文字列",
			dir:     "",
			wantErr: false,
		},
		{
			name:    "存在するディレクトリ",
			wantErr: false,
			setup:   func() string { return t.TempDir() },
		},
		{
			name:    "存在しないディレクトリ",
			dir:     "/nonexistent/dir",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.dir
			if tt.setup != nil {
				dir = tt.setup()
			}
			err := validation.ValidateOutputDir(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMarkdownOutput(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"stdin", "-", true},
		{"通常のパス", "output.md", false},
		{"空文字列", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validation.ValidateMarkdownOutput(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMarkdownOutput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
