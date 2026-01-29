//go:build integration

package integration

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/ageha734/mmdc/handler"
	"github.com/chromedp/chromedp"
)

func TestBuildChromeOptions_Default(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	if len(opts) == 0 {
		t.Error("Expected non-empty chrome options")
	}
}

func TestBuildChromeOptions_WithEnvPath(t *testing.T) {
	tmpDir := t.TempDir()

	originalPath := os.Getenv("CHROMEDP_EXEC_PATH")
	os.Setenv("CHROMEDP_EXEC_PATH", "/usr/bin/chromium")
	defer os.Setenv("CHROMEDP_EXEC_PATH", originalPath)

	opts := handler.BuildChromeOptions(tmpDir)
	if len(opts) == 0 {
		t.Error("Expected non-empty chrome options with custom exec path")
	}
}

func TestSetupBrowserContext_QuietMode(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx := context.Background()
	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, true)
	defer cancelBrowser()

	if browserCtx == nil {
		t.Error("Browser context should not be nil")
	}
}

func TestSetupBrowserContext_NonQuietMode(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx := context.Background()
	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, false)
	defer cancelBrowser()

	if browserCtx == nil {
		t.Error("Browser context should not be nil")
	}
}

func TestCreateTempHTMLFile(t *testing.T) {
	htmlContent := `<!DOCTYPE html><html><body>Test</body></html>`

	fileURL, cleanup, err := handler.CreateTempHTMLFile(htmlContent)
	if err != nil {
		t.Fatalf("Failed to create temp HTML file: %v", err)
	}
	defer cleanup()

	if !strings.HasPrefix(fileURL, "file://") {
		t.Errorf("File URL should start with file://, got: %s", fileURL)
	}

	filePath := strings.TrimPrefix(fileURL, "file://")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Temp file should exist")
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read temp file: %v", err)
	}
	if string(content) != htmlContent {
		t.Errorf("Content mismatch: expected %q, got %q", htmlContent, string(content))
	}

	cleanup()

	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("Temp file should be deleted after cleanup")
	}
}

func TestCreateTempUserDataDir(t *testing.T) {
	userDataDir, cleanup, err := handler.CreateTempUserDataDir()
	if err != nil {
		t.Fatalf("Failed to create temp user data dir: %v", err)
	}
	defer cleanup()

	if userDataDir == "" {
		t.Error("User data dir should not be empty")
	}

	if _, err := os.Stat(userDataDir); os.IsNotExist(err) {
		t.Error("User data directory should exist")
	}

	cleanup()

	if _, err := os.Stat(userDataDir); !os.IsNotExist(err) {
		t.Error("User data directory should be deleted after cleanup")
	}
}

func TestShouldIgnoreError(t *testing.T) {
	testCases := []struct {
		name     string
		msg      string
		expected bool
	}{
		{
			name:     "Unknown IPAddressSpace",
			msg:      "unknown IPAddressSpace in response",
			expected: true,
		},
		{
			name:     "Parse error",
			msg:      "parse error: expected string",
			expected: true,
		},
		{
			name:     "CookiePart",
			msg:      "error with cookiePart",
			expected: true,
		},
		{
			name:     "Random error",
			msg:      "some other error",
			expected: false,
		},
		{
			name:     "Empty string",
			msg:      "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := handler.ShouldIgnoreError(tc.msg)
			if result != tc.expected {
				t.Errorf("Expected %v for message %q, got %v", tc.expected, tc.msg, result)
			}
		})
	}
}

func TestBrowserContext_ActualBrowserLaunch(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx := context.Background()
	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, true)
	defer cancelBrowser()

	err := chromedp.Run(browserCtx,
		chromedp.Navigate("about:blank"),
	)
	if err != nil {
		t.Fatalf("Failed to navigate with browser context: %v", err)
	}
}

func TestBrowserContext_WithTimeout(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, true)
	defer cancelBrowser()

	if browserCtx == nil {
		t.Error("Browser context should be created even with short timeout")
	}
}

func TestBrowserContext_RenderSimplePage(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	ctx := context.Background()
	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, true)
	defer cancelBrowser()

	htmlContent := `<!DOCTYPE html>
<html>
<body>
<div id="test">Hello World</div>
</body>
</html>`

	fileURL, cleanup, err := handler.CreateTempHTMLFile(htmlContent)
	if err != nil {
		t.Fatalf("Failed to create temp HTML file: %v", err)
	}
	defer cleanup()

	var content string
	err = chromedp.Run(browserCtx,
		chromedp.Navigate(fileURL),
		chromedp.WaitVisible("#test", chromedp.ByQuery),
		chromedp.Text("#test", &content, chromedp.ByQuery),
	)
	if err != nil {
		t.Fatalf("Failed to execute chromedp actions: %v", err)
	}

	if content != "Hello World" {
		t.Errorf("Expected 'Hello World', got %q", content)
	}
}

func TestBrowserContext_MultipleContexts(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	opts1 := handler.BuildChromeOptions(tmpDir1)
	allocCtx1, cancelAlloc1 := chromedp.NewExecAllocator(context.Background(), opts1...)
	defer cancelAlloc1()

	ctx1, cancelBrowser1 := handler.SetupBrowserContext(context.Background(), allocCtx1, true)
	defer cancelBrowser1()

	opts2 := handler.BuildChromeOptions(tmpDir2)
	allocCtx2, cancelAlloc2 := chromedp.NewExecAllocator(context.Background(), opts2...)
	defer cancelAlloc2()

	ctx2, cancelBrowser2 := handler.SetupBrowserContext(context.Background(), allocCtx2, true)
	defer cancelBrowser2()

	err1 := chromedp.Run(ctx1, chromedp.Navigate("about:blank"))
	if err1 != nil {
		t.Errorf("First context failed: %v", err1)
	}

	err2 := chromedp.Run(ctx2, chromedp.Navigate("about:blank"))
	if err2 != nil {
		t.Errorf("Second context failed: %v", err2)
	}
}

func TestBrowserContext_ErrorHandling(t *testing.T) {
	tmpDir := t.TempDir()

	opts := handler.BuildChromeOptions(tmpDir)
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancelAlloc()

	browserCtx, cancelBrowser := handler.SetupBrowserContext(context.Background(), allocCtx, true)
	defer cancelBrowser()

	var result string
	ctx, cancel := context.WithTimeout(browserCtx, 2*time.Second)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("about:blank"),
		chromedp.Text("#non-existent-element", &result, chromedp.ByQuery),
	)

	if err == nil {
		t.Error("Expected error when trying to get text from non-existent element")
	}
}
