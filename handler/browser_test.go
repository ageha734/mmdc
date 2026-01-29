package handler_test

import (
	"context"
	"os"
	"testing"

	"github.com/ageha734/mmdc/handler"
	"github.com/chromedp/chromedp"
)

func TestBuildChromeOptions(t *testing.T) {
	userDataDir := t.TempDir()
	opts := handler.BuildChromeOptions(userDataDir)

	if len(opts) == 0 {
		t.Fatal("expected options to be returned")
	}
}

func TestBuildChromeOptions_WithExecPath(t *testing.T) {
	originalPath := os.Getenv("CHROMEDP_EXEC_PATH")
	defer func() {
		if originalPath != "" {
			os.Setenv("CHROMEDP_EXEC_PATH", originalPath)
		} else {
			os.Unsetenv("CHROMEDP_EXEC_PATH")
		}
	}()

	testPath := "/test/chrome/path"
	os.Setenv("CHROMEDP_EXEC_PATH", testPath)

	userDataDir := t.TempDir()
	opts := handler.BuildChromeOptions(userDataDir)

	if len(opts) == 0 {
		t.Error("expected options to be returned")
	}
}

func TestShouldIgnoreError(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		expected bool
	}{
		{"unknown IPAddressSpace", "unknown IPAddressSpace error", true},
		{"parse error", "parse error: expected string", true},
		{"cookiePart", "cookiePart error", true},
		{"other error", "some other error", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.ShouldIgnoreError(tt.msg)
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestSetupBrowserContext(t *testing.T) {
	ctx := context.Background()
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, chromedp.Headless)
	defer cancelAlloc()

	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, false)
	if browserCtx == nil {
		t.Fatal("expected browserCtx to be non-nil")
	}
	defer cancelBrowser()

	if browserCtx.Err() != nil {
		t.Errorf("browserCtx should not be cancelled: %v", browserCtx.Err())
	}
}

func TestSetupBrowserContext_Quiet(t *testing.T) {
	ctx := context.Background()
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, chromedp.Headless)
	defer cancelAlloc()

	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, true)
	if browserCtx == nil {
		t.Fatal("expected browserCtx to be non-nil")
	}
	defer cancelBrowser()
}

func TestSetupBrowserContext_Cancel(t *testing.T) {
	ctx := context.Background()
	allocCtx, cancelAlloc := chromedp.NewExecAllocator(ctx, chromedp.Headless)
	defer cancelAlloc()

	browserCtx, cancelBrowser := handler.SetupBrowserContext(ctx, allocCtx, false)
	cancelBrowser()

	if browserCtx.Err() == nil {
		t.Error("browserCtx should be cancelled after cancelBrowser()")
	}
}
