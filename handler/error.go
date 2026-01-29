package handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/ageha734/mmdc/config"
	"github.com/chromedp/chromedp"
)

func CheckRenderError(ctx context.Context, errorLogs []string) error {
	if len(errorLogs) > 0 {
		return fmt.Errorf("mermaid syntax error detected: %s", strings.Join(errorLogs, "; "))
	}

	var attr string
	if err := chromedp.AttributeValue("#container", "data-rendered", &attr, nil, chromedp.ByQuery).Do(ctx); err == nil {
		if attr == config.RenderedAttrError {
			var errorMsg string
			_ = chromedp.Evaluate(`document.querySelector("#container p")?.textContent || "Unknown error"`, &errorMsg).Do(ctx)
			if len(errorLogs) > 0 {
				errorMsg += "; " + strings.Join(errorLogs, "; ")
			}
			return fmt.Errorf("mermaid render error: %s", errorMsg)
		}
	}
	return nil
}

func FormatRenderError(format string, err error, errorLogs []string, consoleLogs []string, quiet bool) error {
	errorMsg := fmt.Sprintf("failed to render %s: %v", format, err)
	if len(errorLogs) > 0 {
		errorMsg += fmt.Sprintf("\nErrors: %s", strings.Join(errorLogs, "\n"))
	}
	if len(consoleLogs) > 0 && !quiet {
		errorMsg += fmt.Sprintf("\nConsole logs:\n%s", strings.Join(consoleLogs, "\n"))
	}
	return fmt.Errorf("%s", errorMsg)
}
