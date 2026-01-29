//go:build ignore

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Usage: merge-coverage.go <file1> [file2] ...")
		os.Exit(1)
	}

	output, err := os.Create("coverage-all.out")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating output file: %v\n", err)
		os.Exit(1)
	}
	defer output.Close()

	fmt.Fprintln(output, "mode: atomic")

	for _, filename := range os.Args[1:] {
		if err := processFile(filename, output); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
	}
}

func processFile(filename string, output *os.File) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("cannot open %s: %w", filename, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		fmt.Fprintln(output, line)
	}

	return scanner.Err()
}
