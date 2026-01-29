package handler_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ageha734/mmdc/handler"
)

type mockTempFile struct {
	name       string
	writeErr   error
	closeErr   error
	writeCalls int
	closeCalls int
}

func (m *mockTempFile) Write(p []byte) (int, error) {
	m.writeCalls++
	if m.writeErr != nil {
		return 0, m.writeErr
	}
	return len(p), nil
}

func (m *mockTempFile) Close() error {
	m.closeCalls++
	return m.closeErr
}

func (m *mockTempFile) Name() string {
	return m.name
}

func TestCreateTempHTMLFile(t *testing.T) {
	htmlTemplate := "<html><body>Test</body></html>"
	fileURL, cleanup, err := handler.CreateTempHTMLFile(htmlTemplate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cleanup()

	if fileURL == "" {
		t.Fatal("fileURL should not be empty")
	}
	if !strings.HasPrefix(fileURL, "file://") {
		t.Errorf("fileURL should start with 'file://', got %q", fileURL)
	}

	filePath := strings.TrimPrefix(fileURL, "file://")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Errorf("temp file should exist: %v", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read temp file: %v", err)
	}
	if string(content) != htmlTemplate {
		t.Errorf("file content mismatch: expected %q, got %q", htmlTemplate, string(content))
	}

	cleanup()
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		t.Error("temp file should be deleted after cleanup")
	}
}

func TestCreateTempUserDataDir(t *testing.T) {
	userDataDir, cleanup, err := handler.CreateTempUserDataDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cleanup()

	if userDataDir == "" {
		t.Fatal("userDataDir should not be empty")
	}

	info, err := os.Stat(userDataDir)
	if err != nil {
		t.Fatalf("userDataDir should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("userDataDir should be a directory")
	}

	cleanup()
	if _, err := os.Stat(userDataDir); !os.IsNotExist(err) {
		t.Error("userDataDir should be deleted after cleanup")
	}
}

func TestCreateTempHTMLFile_FilePattern(t *testing.T) {
	htmlTemplate := "test"
	fileURL, cleanup, err := handler.CreateTempHTMLFile(htmlTemplate)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer cleanup()

	filePath := strings.TrimPrefix(fileURL, "file://")
	fileName := filepath.Base(filePath)
	if !strings.HasPrefix(fileName, "mmdc-") {
		t.Errorf("file name should start with 'mmdc-', got %q", fileName)
	}
	if !strings.HasSuffix(fileName, ".html") {
		t.Errorf("file name should end with '.html', got %q", fileName)
	}
}

func TestCreateTempHTMLFile_CreateTempError(t *testing.T) {
	restore := handler.SetOsCreateTempForTesting(func(dir, pattern string) (handler.TempFile, error) {
		return nil, fmt.Errorf("mock create temp error")
	})
	defer restore()

	_, _, err := handler.CreateTempHTMLFile("test")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to create temp file") {
		t.Errorf("expected 'failed to create temp file' error, got %v", err)
	}
}

func TestCreateTempHTMLFile_WriteError(t *testing.T) {
	tmpDir := t.TempDir()
	mockFile := &mockTempFile{
		name:     filepath.Join(tmpDir, "mock-temp.html"),
		writeErr: fmt.Errorf("mock write error"),
	}

	restore := handler.SetOsCreateTempForTesting(func(dir, pattern string) (handler.TempFile, error) {
		return mockFile, nil
	})
	defer restore()

	_, _, err := handler.CreateTempHTMLFile("test")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to write temp file") {
		t.Errorf("expected 'failed to write temp file' error, got %v", err)
	}
	if mockFile.closeCalls != 1 {
		t.Errorf("expected Close to be called once, got %d", mockFile.closeCalls)
	}
}

func TestCreateTempHTMLFile_CloseError(t *testing.T) {
	tmpDir := t.TempDir()
	mockFile := &mockTempFile{
		name:     filepath.Join(tmpDir, "mock-temp.html"),
		closeErr: fmt.Errorf("mock close error"),
	}

	restore := handler.SetOsCreateTempForTesting(func(dir, pattern string) (handler.TempFile, error) {
		return mockFile, nil
	})
	defer restore()

	_, _, err := handler.CreateTempHTMLFile("test")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to close temp file") {
		t.Errorf("expected 'failed to close temp file' error, got %v", err)
	}
}

func TestCreateTempUserDataDir_MkdirTempError(t *testing.T) {
	restore := handler.SetOsMkdirTempForTesting(func(dir, pattern string) (string, error) {
		return "", fmt.Errorf("mock mkdir temp error")
	})
	defer restore()

	_, _, err := handler.CreateTempUserDataDir()
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "failed to create temp user data dir") {
		t.Errorf("expected 'failed to create temp user data dir' error, got %v", err)
	}
}
