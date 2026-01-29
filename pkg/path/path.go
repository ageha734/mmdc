package path

import (
	"path/filepath"
	"strings"

	"github.com/ageha734/mmdc/domain"
)

const stdinPath = "-"

func IsMarkdownFile(path string) bool {
	if path == stdinPath {
		return false
	}
	return strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")
}

func DetermineOutputPath(inputPath, outputPath string) string {
	if outputPath != "" {
		return outputPath
	}
	if inputPath == stdinPath {
		return stdinPath
	}
	return inputPath + ".svg"
}

func DetermineOutputFormatFromPath(outputPath string) domain.OutputFormat {
	ext := strings.ToLower(filepath.Ext(outputPath))
	switch ext {
	case ".svg":
		return domain.OutputFormatSVG
	case ".png":
		return domain.OutputFormatPNG
	case ".pdf":
		return domain.OutputFormatPDF
	default:
		return domain.OutputFormatSVG
	}
}

func DetermineOutputFormatFromExtension(ext string) domain.OutputFormat {
	ext = strings.TrimPrefix(ext, ".")
	switch ext {
	case "svg", "png", "pdf":
		return domain.OutputFormat(ext)
	default:
		return domain.OutputFormatSVG
	}
}
