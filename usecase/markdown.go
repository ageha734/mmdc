package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/file"
	"github.com/ageha734/mmdc/infra/logger"
	"github.com/ageha734/mmdc/infra/renderer"
	"github.com/ageha734/mmdc/pkg/errors"
	"github.com/ageha734/mmdc/pkg/markdown"
)

type MarkdownUsecase struct {
	renderer renderer.Renderer
	fileRW   *file.FileReaderWriter
	logger   logger.Interface
}

func NewMarkdownUsecase(r renderer.Renderer, fileRW *file.FileReaderWriter, log logger.Interface) *MarkdownUsecase {
	return &MarkdownUsecase{
		renderer: r,
		fileRW:   fileRW,
		logger:   log,
	}
}

func (u *MarkdownUsecase) ProcessMarkdownFile(ctx context.Context, inputPath, outputPath string, format domain.OutputFormat, artefactsPath string, renderConfig *config.RenderConfig, quiet bool) error {
	content, err := u.fileRW.ReadString(inputPath)
	if err != nil {
		return errors.WrapReadError(err, inputPath)
	}

	matches := markdown.ExtractAllMermaidFromMarkdown(content)

	if len(matches) == 0 {
		if !quiet {
			u.logger.Info("No mermaid charts found in Markdown input")
		}
		return nil
	}

	if !quiet {
		u.logger.Info(fmt.Sprintf("Found %d mermaid charts in Markdown input", len(matches)))
	}

	if outputPath == "" {
		return nil
	}

	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		absOutputPath = outputPath
	}

	outputDir, err := u.prepareOutputDirectory(absOutputPath, artefactsPath)
	if err != nil {
		return err
	}

	var imageInfos []domain.ImageInfo
	for i, match := range matches {
		outputFile := markdown.GenerateOutputFileName(absOutputPath, i, format)
		if filepath.Dir(outputFile) != outputDir {
			outputFile = filepath.Join(outputDir, filepath.Base(outputFile))
		}

		diagram := &domain.MermaidDiagram{
			Code:      match.Code,
			Index:     i,
			FullMatch: match.FullMatch,
		}

		result, err := u.renderer.Render(ctx, diagram, format, renderConfig)
		if err != nil {
			codePreview := errors.FormatCodePreview(match.Code, 200, 5)
			return errors.WrapMarkdownRenderError(err, i+1, codePreview)
		}

		if err := u.fileRW.Write(outputFile, result, config.FilePermission); err != nil {
			return errors.WrapWriteError(err, outputFile)
		}

		relPath := markdown.CalculateRelativePath(absOutputPath, outputFile)

		if !quiet {
			u.logger.Success(fmt.Sprintf(" ✅ %s", relPath))
		}

		imageInfos = append(imageInfos, domain.ImageInfo{
			URL:   relPath,
			Title: "",
			Alt:   fmt.Sprintf("diagram %d", i+1),
		})
	}

	if strings.HasSuffix(absOutputPath, ".md") || strings.HasSuffix(absOutputPath, ".markdown") {
		newContent := markdown.ReplaceMermaidWithImages(content, imageInfos)

		if err := u.fileRW.Write(absOutputPath, []byte(newContent), config.FilePermission); err != nil {
			return errors.WrapWriteError(err, absOutputPath)
		}
		if !quiet {
			u.logger.Success(fmt.Sprintf(" ✅ %s", absOutputPath))
		}
	}

	return nil
}

func (u *MarkdownUsecase) prepareOutputDirectory(outputPath, artefactsPath string) (string, error) {
	outputDir := filepath.Dir(outputPath)
	if artefactsPath != "" {
		absArtefactsPath, err := filepath.Abs(artefactsPath)
		if err != nil {
			absArtefactsPath = artefactsPath
		}
		if err := u.fileRW.MkdirAll(absArtefactsPath, config.DirPermission); err != nil {
			return "", fmt.Errorf("failed to create artefacts directory: %w", err) //nolint:misspell // intentional
		}
		return absArtefactsPath, nil
	}
	if err := u.fileRW.MkdirAll(outputDir, config.DirPermission); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}
	return outputDir, nil
}
