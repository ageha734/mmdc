package renderer

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
	"github.com/ageha734/mmdc/handler"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type ChromedpRenderer struct {
	opts RenderOptions
}

func NewChromedpRenderer(opts RenderOptions) *ChromedpRenderer {
	return &ChromedpRenderer{opts: opts}
}

func (r *ChromedpRenderer) Render(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
	userDataDir, cleanupUserDataDir, err := handler.CreateTempUserDataDir()
	if err != nil {
		return nil, err
	}
	defer cleanupUserDataDir()

	optsChromedp := handler.BuildChromeOptions(userDataDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, optsChromedp...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, r.opts.Quiet)
	defer cancelBrowser()

	htmlGen := NewHTMLGenerator(renderConfig)
	htmlTemplate, err := htmlGen.Generate(diagram.Code, renderConfig.SVGID)
	if err != nil {
		return nil, err
	}

	fileURL, cleanupTempFile, err := handler.CreateTempHTMLFile(htmlTemplate)
	if err != nil {
		return nil, err
	}
	defer cleanupTempFile()

	var consoleLogs []string
	var errorLogs []string
	chromedp.ListenTarget(browserCtx, func(ev interface{}) {
		switch ev := ev.(type) {
		case *runtime.EventConsoleAPICalled:
			if ev.Type == runtime.APITypeError || ev.Type == runtime.APITypeWarning {
				for _, arg := range ev.Args {
					if arg.Value != nil {
						errorLogs = append(errorLogs, fmt.Sprintf("Console %s: %s", ev.Type, string(arg.Value)))
					}
				}
			}
			if !r.opts.Quiet {
				for _, arg := range ev.Args {
					if arg.Value != nil {
						consoleLogs = append(consoleLogs, fmt.Sprintf("Console: %s", string(arg.Value)))
					}
				}
			}
		case *runtime.EventExceptionThrown:
			if ev.ExceptionDetails != nil {
				errorLogs = append(errorLogs, fmt.Sprintf("Exception: %s", ev.ExceptionDetails.Text))
			}
		}
	})

	switch format {
	case domain.OutputFormatSVG:
		return r.renderSVG(browserCtx, fileURL, errorLogs, consoleLogs)
	case domain.OutputFormatPNG:
		return r.renderPNG(browserCtx, fileURL, errorLogs)
	case domain.OutputFormatPDF:
		return r.renderPDF(browserCtx, fileURL, errorLogs)
	default:
		return nil, fmt.Errorf("unsupported output format: %s", format)
	}
}

func (r *ChromedpRenderer) renderSVG(ctx context.Context, fileURL string, errorLogs, consoleLogs []string) ([]byte, error) {
	var svgHTML string
	var rendered bool
	err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(config.ScriptLoadWaitTime),
		chromedp.ActionFunc(func(ctx context.Context) error {
			deadline := time.Now().Add(config.PollingTimeout)
			checkCount := 0
			for time.Now().Before(deadline) {
				checkCount++

				if err := handler.CheckRenderError(ctx, errorLogs); err != nil {
					return err
				}

				var attr string
				var hasSvg bool

				if err := chromedp.AttributeValue("#container", "data-rendered", &attr, nil, chromedp.ByQuery).Do(ctx); err == nil {
					if attr == config.RenderedAttrTrue {
						if err := chromedp.Evaluate(`document.querySelector("#container svg") !== null`, &hasSvg).Do(ctx); err == nil && hasSvg {
							rendered = true
							return nil
						}
					}
				}

				if checkCount%config.PollingCheckInterval == 0 {
					if err := chromedp.Evaluate(`document.querySelector("#container svg") !== null`, &hasSvg).Do(ctx); err == nil && hasSvg {
						rendered = true
						return nil
					}
				}

				time.Sleep(config.PollingInterval)
			}
			timeoutMsg := fmt.Sprintf("timeout waiting for mermaid render (checked %d times)", checkCount)
			if len(errorLogs) > 0 {
				timeoutMsg += "; errors: " + strings.Join(errorLogs, "; ")
			}
			return fmt.Errorf("%s", timeoutMsg)
		}),
		chromedp.WaitVisible("#container svg", chromedp.ByQuery),
		chromedp.Sleep(config.FinalWaitTime),
		chromedp.OuterHTML("#container", &svgHTML, chromedp.ByQuery),
	)
	if err != nil {
		return nil, handler.FormatRenderError("SVG", err, errorLogs, consoleLogs, r.opts.Quiet)
	}
	if !rendered {
		return nil, handler.FormatRenderError("SVG", fmt.Errorf("mermaid render did not complete"), errorLogs, consoleLogs, r.opts.Quiet)
	}

	svgRe := regexp.MustCompile(`(?s)<svg[^>]*>.*?</svg>`)
	svgMatch := svgRe.FindString(svgHTML)
	if svgMatch != "" {
		return []byte(svgMatch), nil
	}
	return []byte(svgHTML), nil
}

func (r *ChromedpRenderer) renderPNG(ctx context.Context, fileURL string, errorLogs []string) ([]byte, error) {
	var result []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(config.ScriptLoadWaitTime),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return handler.CheckRenderError(ctx, errorLogs)
		}),
		chromedp.WaitVisible("#container[data-rendered]", chromedp.ByQuery),
		chromedp.WaitVisible("#container svg", chromedp.ByQuery),
		chromedp.Sleep(config.FinalWaitTime),
		chromedp.CaptureScreenshot(&result),
	)
	if err != nil {
		return nil, handler.FormatRenderError("PNG", err, errorLogs, nil, r.opts.Quiet)
	}
	return result, nil
}

func (r *ChromedpRenderer) renderPDF(ctx context.Context, fileURL string, errorLogs []string) ([]byte, error) {
	var pdfData []byte
	err := chromedp.Run(ctx,
		chromedp.Navigate(fileURL),
		chromedp.WaitReady("body", chromedp.ByQuery),
		chromedp.Sleep(config.ScriptLoadWaitTime),
		chromedp.ActionFunc(func(ctx context.Context) error {
			return handler.CheckRenderError(ctx, errorLogs)
		}),
		chromedp.WaitVisible("#container[data-rendered]", chromedp.ByQuery),
		chromedp.WaitVisible("#container svg", chromedp.ByQuery),
		chromedp.Sleep(config.FinalWaitTime),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdfData, _, err = page.PrintToPDF().
				WithPrintBackground(true).
				WithTransferMode(page.PrintToPDFTransferModeReturnAsBase64).
				Do(ctx)
			if err != nil {
				return fmt.Errorf("failed to generate PDF: %w", err)
			}
			pdfData, err = base64.StdEncoding.DecodeString(string(pdfData))
			if err != nil {
				return fmt.Errorf("failed to decode PDF: %w", err)
			}
			return nil
		}),
	)
	if err != nil {
		return nil, handler.FormatRenderError("PDF", err, errorLogs, nil, r.opts.Quiet)
	}
	return pdfData, nil
}
