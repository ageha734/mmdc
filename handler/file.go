package handler

import (
	"fmt"
	"io"
	"os"

	"github.com/ageha734/mmdc/config"
)

type TempFile interface {
	io.WriteCloser
	Name() string
}

var osCreateTemp = func(dir, pattern string) (TempFile, error) {
	return os.CreateTemp(dir, pattern)
}

var osMkdirTemp = os.MkdirTemp

func SetOsCreateTempForTesting(f func(dir, pattern string) (TempFile, error)) func() {
	original := osCreateTemp
	osCreateTemp = f
	return func() {
		osCreateTemp = original
	}
}

func SetOsMkdirTempForTesting(f func(dir, pattern string) (string, error)) func() {
	original := osMkdirTemp
	osMkdirTemp = f
	return func() {
		osMkdirTemp = original
	}
}

func CreateTempHTMLFile(htmlTemplate string) (string, func(), error) {
	tmpFile, err := osCreateTemp("", config.TempFilePattern)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp file: %w", err)
	}

	if _, err := tmpFile.Write([]byte(htmlTemplate)); err != nil {
		_ = tmpFile.Close()
		_ = os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("failed to write temp file: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		_ = os.Remove(tmpFile.Name())
		return "", nil, fmt.Errorf("failed to close temp file: %w", err)
	}

	fileURL := "file://" + tmpFile.Name()
	cleanup := func() {
		_ = os.Remove(tmpFile.Name())
	}

	return fileURL, cleanup, nil
}

func CreateTempUserDataDir() (string, func(), error) {
	userDataDir, err := osMkdirTemp("", config.UserDataDirPattern)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp user data dir: %w", err)
	}
	cleanup := func() {
		_ = os.RemoveAll(userDataDir)
	}
	return userDataDir, cleanup, nil
}
