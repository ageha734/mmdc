package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/file"
	"github.com/ageha734/mmdc/infra/logger"
	"github.com/ageha734/mmdc/infra/renderer"
	"github.com/ageha734/mmdc/pkg/errors"
	"github.com/ageha734/mmdc/pkg/path"
	"github.com/ageha734/mmdc/pkg/validation"
)

var osStdout io.Writer = os.Stdout

func SetOsStdoutForTesting(w io.Writer) func() {
	original := osStdout
	osStdout = w
	return func() {
		osStdout = original
	}
}

type RenderUsecase struct {
	renderer renderer.Renderer
	fileRW   *file.FileReaderWriter
	logger   logger.Interface
}

func NewRenderUsecase(r renderer.Renderer, fileRW *file.FileReaderWriter, log logger.Interface) *RenderUsecase {
	return &RenderUsecase{
		renderer: r,
		fileRW:   fileRW,
		logger:   log,
	}
}

func (u *RenderUsecase) RenderFromInput(ctx context.Context, inputPath string, outputPath string, format domain.OutputFormat, renderConfig *config.RenderConfig, artefactsPath string, quiet bool) error {
	inputPath = u.normalizeInputPath(inputPath)

	inputData, err := u.fileRW.ReadString(inputPath)
	if err != nil {
		return errors.WrapReadError(err, inputPath)
	}

	if path.IsMarkdownFile(inputPath) {
		return u.renderMarkdown(ctx, inputPath, outputPath, format, renderConfig, artefactsPath, quiet)
	}

	return u.renderSingle(ctx, inputData, inputPath, outputPath, format, renderConfig, quiet)
}

func (u *RenderUsecase) normalizeInputPath(inputPath string) string {
	if inputPath == "" {
		u.logger.Warn("No input file specified, reading from stdin. " +
			"If you want to specify an input file, please use `-i <input>`. " +
			"You can use `-i -` to read from stdin and to suppress this warning.")
		return config.StdinPath
	}
	if inputPath != config.StdinPath {
		if err := validation.ValidateInputFile(inputPath); err != nil {
			u.logger.ErrorExit(err.Error())
		}
	}
	return inputPath
}

func (u *RenderUsecase) renderSingle(ctx context.Context, mermaidCode, inputPath, outputPath string, format domain.OutputFormat, renderConfig *config.RenderConfig, quiet bool) error {
	outputPath = path.DetermineOutputPath(inputPath, outputPath)

	if outputPath == config.StdinPath {
		return u.renderToStdout(ctx, mermaidCode, format, renderConfig, quiet)
	}

	if err := u.validateOutputPath(outputPath); err != nil {
		u.logger.ErrorExit(err.Error())
	}

	format = u.determineFormat(outputPath, format)
	if err := validation.ValidateOutputFormat(format); err != nil {
		u.logger.ErrorExit(err.Error())
	}

	if !quiet {
		u.logger.Info("Generating single mermaid chart")
	}

	result, err := u.renderMermaid(ctx, mermaidCode, format, renderConfig, quiet)
	if err != nil {
		return errors.WrapRenderError(err)
	}

	if err := u.fileRW.Write(outputPath, result, config.FilePermission); err != nil {
		return errors.WrapWriteError(err, outputPath)
	}

	if !quiet {
		u.logger.Success(fmt.Sprintf(" ✅ %s", outputPath))
	}

	return nil
}

func (u *RenderUsecase) renderToStdout(ctx context.Context, mermaidCode string, format domain.OutputFormat, renderConfig *config.RenderConfig, quiet bool) error {
	if format == "" {
		format = domain.OutputFormatSVG
		u.logger.Warn("No output format specified, using svg. " +
			"If you want to specify an output format and suppress this warning, " +
			"please use `-e <format>`.")
	}
	result, err := u.renderMermaid(ctx, mermaidCode, format, renderConfig, quiet)
	if err != nil {
		return errors.WrapRenderError(err)
	}
	if _, err := osStdout.Write(result); err != nil {
		return fmt.Errorf("failed to write to stdout: %w", err)
	}
	return nil
}

func (u *RenderUsecase) validateOutputPath(outputPath string) error {
	if err := validation.ValidateOutputFileExtension(outputPath); err != nil {
		return err
	}
	outputDir := filepath.Dir(outputPath)
	return validation.ValidateOutputDir(outputDir)
}

func (u *RenderUsecase) determineFormat(outputPath string, format domain.OutputFormat) domain.OutputFormat {
	if format != "" {
		return format
	}
	return path.DetermineOutputFormatFromPath(outputPath)
}

func (u *RenderUsecase) renderMermaid(ctx context.Context, mermaidCode string, format domain.OutputFormat, renderConfig *config.RenderConfig, _ bool) ([]byte, error) {
	diagram := &domain.MermaidDiagram{Code: mermaidCode}
	return u.renderer.Render(ctx, diagram, format, renderConfig)
}

func (u *RenderUsecase) renderMarkdown(ctx context.Context, inputPath, outputPath string, format domain.OutputFormat, renderConfig *config.RenderConfig, artefactsPath string, quiet bool) error {
	outputPath = u.determineMarkdownOutputPath(inputPath, outputPath)
	if outputPath == "" {
		return nil
	}
	format = u.determineMarkdownFormat(outputPath, format)

	markdownUsecase := NewMarkdownUsecase(u.renderer, u.fileRW, u.logger)
	return markdownUsecase.ProcessMarkdownFile(ctx, inputPath, outputPath, format, artefactsPath, renderConfig, quiet)
}

func (u *RenderUsecase) determineMarkdownOutputPath(inputPath, outputPath string) string {
	switch outputPath {
	case "":
		absInputPath, err := filepath.Abs(inputPath)
		if err != nil {
			return strings.TrimSuffix(inputPath, filepath.Ext(inputPath)) + ".svg"
		}
		return strings.TrimSuffix(absInputPath, filepath.Ext(absInputPath)) + ".svg"
	case config.StdinPath:
		u.logger.ErrorExit(validation.ValidateMarkdownOutput(outputPath).Error())
		return ""
	default:
		return outputPath
	}
}

func (u *RenderUsecase) determineMarkdownFormat(outputPath string, format domain.OutputFormat) domain.OutputFormat {
	if format != "" {
		return format
	}
	ext := strings.ToLower(filepath.Ext(outputPath))
	if ext == ".md" || ext == ".markdown" {
		return domain.OutputFormatSVG
	}
	return path.DetermineOutputFormatFromExtension(strings.TrimPrefix(ext, "."))
}
