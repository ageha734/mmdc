package file

import (
	"fmt"
	"io"
	"os"
)

var osStdin io.Reader = os.Stdin

func SetOsStdinForTesting(r io.Reader) func() {
	original := osStdin
	osStdin = r
	return func() {
		osStdin = original
	}
}

type FileReader interface {
	Read(path string) ([]byte, error)
	ReadString(path string) (string, error)
}

type FileWriter interface {
	Write(path string, data []byte, perm os.FileMode) error
	MkdirAll(path string, perm os.FileMode) error
}

type FileReaderWriter struct{}

func NewFileReaderWriter() *FileReaderWriter {
	return &FileReaderWriter{}
}

func (f *FileReaderWriter) Read(path string) ([]byte, error) {
	if path == "-" {
		return io.ReadAll(osStdin)
	}
	return os.ReadFile(path) //nolint:gosec // G304: path is trusted input
}

func (f *FileReaderWriter) ReadString(path string) (string, error) {
	data, err := f.Read(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (f *FileReaderWriter) Write(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (f *FileReaderWriter) MkdirAll(path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func ValidateOutputDir(dir string) error {
	if dir != "." && dir != "" {
		if !FileExists(dir) {
			return fmt.Errorf(`output directory "%s/" doesn't exist`, dir)
		}
	}
	return nil
}
