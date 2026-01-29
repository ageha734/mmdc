package errors_test

import (
	stderrors "errors"
	"strings"
	"testing"

	"github.com/ageha734/mmdc/pkg/errors"
)

func TestWrapReadError(t *testing.T) {
	originalErr := stderrors.New("file not found")
	wrappedErr := errors.WrapReadError(originalErr, "/path/to/file")

	if wrappedErr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(wrappedErr.Error(), "failed to read input") {
		t.Errorf("error message should contain 'failed to read input': %v", wrappedErr)
	}
	if !stderrors.Is(wrappedErr, originalErr) {
		t.Error("wrapped error should wrap original error")
	}
}

func TestWrapRenderError(t *testing.T) {
	originalErr := stderrors.New("render failed")
	wrappedErr := errors.WrapRenderError(originalErr)

	if wrappedErr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(wrappedErr.Error(), "failed to render") {
		t.Errorf("error message should contain 'failed to render': %v", wrappedErr)
	}
	if !stderrors.Is(wrappedErr, originalErr) {
		t.Error("wrapped error should wrap original error")
	}
}

func TestWrapWriteError(t *testing.T) {
	originalErr := stderrors.New("write failed")
	wrappedErr := errors.WrapWriteError(originalErr, "/path/to/output")

	if wrappedErr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(wrappedErr.Error(), "failed to write output file") {
		t.Errorf("error message should contain 'failed to write output file': %v", wrappedErr)
	}
	if !strings.Contains(wrappedErr.Error(), "/path/to/output") {
		t.Errorf("error message should contain path: %v", wrappedErr)
	}
	if !stderrors.Is(wrappedErr, originalErr) {
		t.Error("wrapped error should wrap original error")
	}
}

func TestWrapMarkdownRenderError(t *testing.T) {
	originalErr := stderrors.New("render failed")
	codePreview := "graph TD\nA-->B"
	wrappedErr := errors.WrapMarkdownRenderError(originalErr, 1, codePreview)

	if wrappedErr == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(wrappedErr.Error(), "failed to render diagram 1") {
		t.Errorf("error message should contain 'failed to render diagram 1': %v", wrappedErr)
	}
	if !strings.Contains(wrappedErr.Error(), codePreview) {
		t.Errorf("error message should contain code preview: %v", wrappedErr)
	}
	if !stderrors.Is(wrappedErr, originalErr) {
		t.Error("wrapped error should wrap original error")
	}
}

func TestFormatCodePreview(t *testing.T) {
	tests := []struct {
		name      string
		code      string
		maxLength int
		maxLines  int
		checkFunc func(*testing.T, string)
	}{
		{
			name:      "短いコード",
			code:      "graph TD\nA-->B",
			maxLength: 200,
			maxLines:  5,
			checkFunc: func(t *testing.T, result string) {
				if !strings.Contains(result, "graph TD") {
					t.Errorf("result should contain code: %q", result)
				}
			},
		},
		{
			name:      "長いコード（maxLength制限）",
			code:      strings.Repeat("a", 300),
			maxLength: 100,
			maxLines:  5,
			checkFunc: func(t *testing.T, result string) {
				if len(result) > 100+10 {
					t.Errorf("result should be truncated to maxLength: len=%d", len(result))
				}
				if !strings.HasSuffix(result, "...") {
					t.Errorf("result should end with '...': %q", result)
				}
			},
		},
		{
			name:      "複数行コード（maxLines制限）",
			code:      strings.Repeat("line\n", 10),
			maxLength: 200,
			maxLines:  3,
			checkFunc: func(t *testing.T, result string) {
				lines := strings.Split(result, "\n")
				if len(lines) > 4 {
					t.Errorf("result should be truncated to maxLines: lines=%d", len(lines))
				}
				if !strings.HasSuffix(result, "...") {
					t.Errorf("result should end with '...': %q", result)
				}
			},
		},
		{
			name:      "空のコード",
			code:      "",
			maxLength: 200,
			maxLines:  5,
			checkFunc: func(t *testing.T, result string) {
				if result != "" {
					t.Errorf("result should be empty for empty code: %q", result)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errors.FormatCodePreview(tt.code, tt.maxLength, tt.maxLines)
			tt.checkFunc(t, result)
		})
	}
}
