// Package handler provides browser and file handling utilities.
package handler

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/ageha734/mmdc/config"
	"github.com/chromedp/chromedp"
)

func BuildChromeOptions(userDataDir string) []chromedp.ExecAllocatorOption {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.NoSandbox,
		chromedp.UserDataDir(userDataDir),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("disable-setuid-sandbox", true),
		chromedp.Flag("disable-web-security", true),
		chromedp.Flag("disable-default-apps", true),
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("disable-background-timer-throttling", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-breakpad", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-component-update", true),
		chromedp.Flag("disable-features", "TranslateUI,BlinkGenPropertyTrees"),
		chromedp.Flag("disable-hang-monitor", true),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-popup-blocking", true),
		chromedp.Flag("disable-prompt-on-repost", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-translate", true),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("no-crash-upload", true),
		chromedp.Flag("no-pings", true),
		chromedp.Flag("no-zygote", true),
	}

	if execPath := os.Getenv("CHROMEDP_EXEC_PATH"); execPath != "" {
		opts = append(opts, chromedp.ExecPath(execPath))
	}

	return opts
}

func SetupBrowserContext(ctx context.Context, allocCtx context.Context, quiet bool) (context.Context, context.CancelFunc) {
	browserCtx, cancel := chromedp.NewContext(allocCtx, chromedp.WithErrorf(func(format string, v ...interface{}) {
		if !quiet {
			msg := fmt.Sprintf(format, v...)
			if !ShouldIgnoreError(msg) {
				fmt.Fprintf(os.Stderr, "chromedp: %s\n", msg)
			}
		}
	}))
	browserCtx, cancelTimeout := context.WithTimeout(browserCtx, config.RenderTimeout)
	return browserCtx, func() {
		cancelTimeout()
		cancel()
	}
}

func ShouldIgnoreError(msg string) bool {
	ignoredPatterns := []string{
		"unknown IPAddressSpace",
		"parse error: expected string",
		"cookiePart",
	}
	for _, pattern := range ignoredPatterns {
		if strings.Contains(msg, pattern) {
			return true
		}
	}
	return false
}
