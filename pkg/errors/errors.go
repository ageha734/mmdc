package errors

import (
	"fmt"
	"strings"
)

func WrapReadError(err error, path string) error {
	return fmt.Errorf("failed to read input: %w", err)
}

func WrapRenderError(err error) error {
	return fmt.Errorf("failed to render: %w", err)
}

func WrapWriteError(err error, path string) error {
	return fmt.Errorf("failed to write output file %s: %w", path, err)
}

func WrapMarkdownRenderError(err error, diagramIndex int, codePreview string) error {
	return fmt.Errorf("failed to render diagram %d: %w\nMermaid code preview:\n%s",
		diagramIndex, err, codePreview)
}

func FormatCodePreview(code string, maxLength int, maxLines int) string {
	preview := code
	if len(preview) > maxLength {
		preview = preview[:maxLength] + "..."
	}
	lines := strings.Split(preview, "\n")
	if len(lines) > maxLines {
		preview = strings.Join(lines[:maxLines], "\n") + "\n..."
	}
	return preview
}
