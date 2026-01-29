package logger_test

import (
	"testing"

	"github.com/ageha734/mmdc/infra/logger"
)

func TestNew(t *testing.T) {
	l := logger.New(false)
	if l == nil {
		t.Fatal("New returned nil")
	}

	l2 := logger.New(true)
	if l2 == nil {
		t.Fatal("New returned nil")
	}
}

func TestNewWithExitFunc(t *testing.T) {
	exitCalled := false
	l := logger.NewWithExitFunc(false, func(code int) {
		exitCalled = true
	})
	if l == nil {
		t.Fatal("NewWithExitFunc returned nil")
	}
	if exitCalled {
		t.Error("exit function should not be called yet")
	}
}

func TestLogger_Info(t *testing.T) {
	tests := []struct {
		name  string
		quiet bool
	}{
		{
			name:  "quiet=false",
			quiet: false,
		},
		{
			name:  "quiet=true",
			quiet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.New(tt.quiet)
			if l == nil {
				t.Fatal("logger should not be nil")
			}
			l.Info("test message")
		})
	}
}

func TestLogger_Success(t *testing.T) {
	tests := []struct {
		name  string
		quiet bool
	}{
		{
			name:  "quiet=false",
			quiet: false,
		},
		{
			name:  "quiet=true",
			quiet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logger.New(tt.quiet)
			if l == nil {
				t.Fatal("logger should not be nil")
			}
			l.Success("test message")
		})
	}
}

func TestLogger_Warn(t *testing.T) {
	l := logger.New(false)
	if l == nil {
		t.Fatal("logger should not be nil")
	}
	l.Warn("test warning")
}

func TestLogger_ErrorExit(t *testing.T) {
	var exitCode int
	l := logger.NewWithExitFunc(false, func(code int) {
		exitCode = code
	})
	if l == nil {
		t.Fatal("logger should not be nil")
	}

	l.ErrorExit("test error")

	if exitCode != 1 {
		t.Errorf("expected exit code 1, got %d", exitCode)
	}
}
