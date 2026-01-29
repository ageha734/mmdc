// Package validation provides input validation functions.
package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ageha734/mmdc/domain"
)

func ValidateInputFile(inputPath string) error {
	if inputPath == "" || inputPath == "-" {
		return nil
	}
	if !fileExists(inputPath) {
		return fmt.Errorf(`input file "%s" doesn't exist`, inputPath)
	}
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ValidateOutputFileExtension(outputPath string) error {
	ext := strings.ToLower(filepath.Ext(outputPath))
	validExts := []string{".svg", ".png", ".pdf", ".md", ".markdown"}
	for _, validExt := range validExts {
		if ext == validExt {
			return nil
		}
	}
	return fmt.Errorf(`output file must end with ".md"/".markdown", ".svg", ".png" or ".pdf"`)
}

func ValidateOutputFormat(format domain.OutputFormat) error {
	if !format.IsValid() {
		return fmt.Errorf(`output format must be one of "svg", "png" or "pdf"`)
	}
	return nil
}

func ValidateOutputDir(outputDir string) error {
	if outputDir != "." && outputDir != "" {
		if !fileExists(outputDir) {
			return fmt.Errorf(`output directory "%s/" doesn't exist`, outputDir)
		}
	}
	return nil
}

func ValidateMarkdownOutput(outputPath string) error {
	if outputPath == "-" {
		return fmt.Errorf("cannot use stdout with markdown input")
	}
	return nil
}
