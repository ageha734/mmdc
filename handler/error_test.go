package handler_test

import (
	"strings"
	"testing"

	"github.com/ageha734/mmdc/handler"
)

func TestFormatRenderError(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		err         error
		errorLogs   []string
		consoleLogs []string
		quiet       bool
		checkFunc   func(*testing.T, error)
	}{
		{
			name:   "基本エラー",
			format: "SVG",
			err:    &testError{msg: "test error"},
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error")
				}
				if !strings.Contains(err.Error(), "failed to render SVG") {
					t.Errorf("error message should contain 'failed to render SVG': %v", err)
				}
			},
		},
		{
			name:      "エラーログあり",
			format:    "PNG",
			err:       &testError{msg: "test error"},
			errorLogs: []string{"error1", "error2"},
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error")
				}
				errMsg := err.Error()
				if !strings.Contains(errMsg, "error1") || !strings.Contains(errMsg, "error2") {
					t.Errorf("error message should contain error logs: %v", errMsg)
				}
			},
		},
		{
			name:        "コンソールログあり（quiet=false）",
			format:      "PDF",
			err:         &testError{msg: "test error"},
			consoleLogs: []string{"log1", "log2"},
			quiet:       false,
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error")
				}
				errMsg := err.Error()
				if !strings.Contains(errMsg, "log1") || !strings.Contains(errMsg, "log2") {
					t.Errorf("error message should contain console logs: %v", errMsg)
				}
			},
		},
		{
			name:        "コンソールログあり（quiet=true）",
			format:      "PDF",
			err:         &testError{msg: "test error"},
			consoleLogs: []string{"log1", "log2"},
			quiet:       true,
			checkFunc: func(t *testing.T, err error) {
				if err == nil {
					t.Fatal("expected error")
				}
				errMsg := err.Error()
				if strings.Contains(errMsg, "log1") || strings.Contains(errMsg, "log2") {
					t.Errorf("error message should not contain console logs when quiet: %v", errMsg)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.FormatRenderError(tt.format, tt.err, tt.errorLogs, tt.consoleLogs, tt.quiet)
			tt.checkFunc(t, err)
		})
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
