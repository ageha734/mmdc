package file_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ageha734/mmdc/infra/file"
)

func TestNewFileReaderWriter(t *testing.T) {
	frw := file.NewFileReaderWriter()
	if frw == nil {
		t.Fatal("NewFileReaderWriter returned nil")
	}
}

func TestFileReaderWriter_Read(t *testing.T) {
	frw := file.NewFileReaderWriter()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	testContent := "test content"

	if err := os.WriteFile(tmpFile, []byte(testContent), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	data, err := frw.Read(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != testContent {
		t.Errorf("expected %q, got %q", testContent, string(data))
	}
}

func TestFileReaderWriter_Read_NonExistent(t *testing.T) {
	frw := file.NewFileReaderWriter()
	_, err := frw.Read("/nonexistent/file.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestFileReaderWriter_ReadString_NonExistent(t *testing.T) {
	frw := file.NewFileReaderWriter()
	_, err := frw.ReadString("/nonexistent/file.txt")
	if err == nil {
		t.Fatal("expected error for non-existent file")
	}
}

func TestFileReaderWriter_Read_Stdin(t *testing.T) {
	testContent := "stdin test content"
	restore := file.SetOsStdinForTesting(strings.NewReader(testContent))
	defer restore()

	frw := file.NewFileReaderWriter()
	data, err := frw.Read("-")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != testContent {
		t.Errorf("expected %q, got %q", testContent, string(data))
	}
}

func TestFileReaderWriter_ReadString_Stdin(t *testing.T) {
	testContent := "stdin test content"
	restore := file.SetOsStdinForTesting(strings.NewReader(testContent))
	defer restore()

	frw := file.NewFileReaderWriter()
	content, err := frw.ReadString("-")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != testContent {
		t.Errorf("expected %q, got %q", testContent, content)
	}
}

func TestFileReaderWriter_ReadString(t *testing.T) {
	frw := file.NewFileReaderWriter()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	testContent := "test content"

	if err := os.WriteFile(tmpFile, []byte(testContent), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	content, err := frw.ReadString(tmpFile)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != testContent {
		t.Errorf("expected %q, got %q", testContent, content)
	}
}

func TestFileReaderWriter_Write(t *testing.T) {
	frw := file.NewFileReaderWriter()
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	testContent := []byte("test content")

	err := frw.Write(tmpFile, testContent, 0o644)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}
	if string(data) != string(testContent) {
		t.Errorf("expected %q, got %q", string(testContent), string(data))
	}
}

func TestFileReaderWriter_MkdirAll(t *testing.T) {
	frw := file.NewFileReaderWriter()
	tmpDir := filepath.Join(t.TempDir(), "nested", "dir", "path")

	err := frw.MkdirAll(tmpDir, 0o755)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	info, err := os.Stat(tmpDir)
	if err != nil {
		t.Fatalf("directory should exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("expected directory")
	}
}

func TestFileExists(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(tmpFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	if !file.FileExists(tmpFile) {
		t.Error("expected file to exist")
	}

	if file.FileExists("/nonexistent/file.txt") {
		t.Error("expected file not to exist")
	}
}

func TestValidateOutputDir(t *testing.T) {
	tests := []struct {
		name    string
		dir     string
		wantErr bool
		setup   func() string
	}{
		{
			name:    "存在するディレクトリ",
			wantErr: false,
			setup:   func() string { return t.TempDir() },
		},
		{
			name:    "存在しないディレクトリ",
			wantErr: true,
			setup:   func() string { return "/nonexistent/dir/path" },
		},
		{
			name:    "カレントディレクトリ",
			dir:     ".",
			wantErr: false,
		},
		{
			name:    "空文字列",
			dir:     "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := tt.dir
			if tt.setup != nil {
				dir = tt.setup()
			}
			err := file.ValidateOutputDir(dir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateOutputDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
