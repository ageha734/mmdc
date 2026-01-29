// Package main provides the CLI entrypoint for mmdc.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/infra/file"
	"github.com/ageha734/mmdc/infra/logger"
	"github.com/ageha734/mmdc/infra/renderer"
	"github.com/ageha734/mmdc/usecase"
	"github.com/spf13/cobra"
)

var version = "1.0.0"

type renderOptions struct {
	input                 string
	output                string
	outputFormat          string
	theme                 string
	width                 int
	height                int
	backgroundColor       string
	scale                 int
	quiet                 bool
	configFile            string
	cssFile               string
	svgID                 string
	pdfFit                bool
	artefacts             string //nolint:misspell // intentional British spelling
	puppeteerConfig       string
	iconPacks             []string
	iconPacksNamesAndUrls []string
}

func main() {
	rootCmd := &cobra.Command{
		Use:   "mmdc",
		Short: "Mermaid CLI - Convert Mermaid diagrams to SVG/PNG/PDF",
		Long: `mmdc is a command-line tool for converting Mermaid diagrams to various image formats.
It uses chromedp to render diagrams in a headless browser.`,
		Version: version,
	}

	renderCmd := newRenderCommand()
	renderCmd.Use = ""
	rootCmd.AddCommand(renderCmd)

	rootCmd.RunE = renderCmd.RunE
	rootCmd.Flags().AddFlagSet(renderCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func newRenderCommand() *cobra.Command {
	opts := &renderOptions{}

	cmd := &cobra.Command{
		Use:   "render",
		Short: "Render Mermaid diagram",
		Long:  "Render a Mermaid diagram from input file to output format (SVG/PNG/PDF)",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runRender(opts)
		},
	}

	cmd.Flags().StringVarP(&opts.input, "input", "i", "", "Input mermaid file or markdown file (required)")
	cmd.Flags().StringVarP(&opts.output, "output", "o", "", "Output file path (default: input + .svg)")
	cmd.Flags().StringVarP(&opts.outputFormat, "outputFormat", "e", "", "Output format: svg, png, or pdf")
	cmd.Flags().StringVarP(&opts.theme, "theme", "t", "default", "Theme: default, forest, dark, or neutral")
	cmd.Flags().IntVarP(&opts.width, "width", "w", 800, "Page width")
	cmd.Flags().IntVarP(&opts.height, "height", "H", 600, "Page height")
	cmd.Flags().StringVarP(&opts.backgroundColor, "backgroundColor", "b", "white", "Background color")
	cmd.Flags().IntVarP(&opts.scale, "scale", "s", 1, "Scale factor")
	cmd.Flags().BoolVarP(&opts.quiet, "quiet", "q", false, "Suppress log output")
	cmd.Flags().StringVarP(&opts.configFile, "configFile", "c", "", "Mermaid config JSON file")
	cmd.Flags().StringVarP(&opts.cssFile, "cssFile", "C", "", "CSS file")
	cmd.Flags().StringVarP(&opts.svgID, "svgID", "I", "", "SVG element ID")
	cmd.Flags().BoolVarP(&opts.pdfFit, "pdfFit", "f", false, "Scale PDF to fit chart")
	cmd.Flags().StringVarP(&opts.artefacts, "artefacts", "a", "", "Output artefacts path (only for Markdown input)") //nolint:misspell // intentional
	cmd.Flags().StringVarP(&opts.puppeteerConfig, "puppeteerConfigFile", "p", "", "Puppeteer config JSON file")
	cmd.Flags().StringArrayVar(&opts.iconPacks, "iconPacks", []string{}, "Icon packs to use")
	cmd.Flags().StringArrayVar(&opts.iconPacksNamesAndUrls, "iconPacksNamesAndUrls", []string{}, "Icon packs with custom URLs (format: prefix#url)")

	return cmd
}

func runRender(opts *renderOptions) error {
	renderConfig := config.NewRenderConfig()
	renderConfig.Theme = opts.theme
	renderConfig.Width = opts.width
	renderConfig.Height = opts.height
	renderConfig.BackgroundColor = opts.backgroundColor
	renderConfig.Scale = opts.scale
	renderConfig.ConfigFile = opts.configFile
	renderConfig.CSSFile = opts.cssFile
	renderConfig.SVGID = opts.svgID
	renderConfig.PDFFit = opts.pdfFit

	var format domain.OutputFormat
	if opts.outputFormat != "" {
		format = domain.OutputFormat(opts.outputFormat)
	}

	log := logger.New(opts.quiet)
	fileRW := file.NewFileReaderWriter()
	rendererImpl := renderer.NewChromedpRenderer(renderer.RenderOptions{Quiet: opts.quiet})
	usecase := usecase.NewRenderUsecase(rendererImpl, fileRW, log)

	ctx := context.Background()
	return usecase.RenderFromInput(ctx, opts.input, opts.output, format, renderConfig, opts.artefacts, opts.quiet)
}
