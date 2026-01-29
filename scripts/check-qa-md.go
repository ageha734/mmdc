package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	e2eDir := "test/e2e"
	qaFile := "test/QA.md"

	testNames, err := extractTestNames(e2eDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting test names: %v\n", err)
		os.Exit(1)
	}

	qaTestNames, err := extractQATestNames(qaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting QA test names: %v\n", err)
		os.Exit(1)
	}

	missingInQA := findMissing(testNames, qaTestNames)
	missingInE2E := findMissing(qaTestNames, testNames)

	hasErrors := false

	if len(missingInQA) > 0 {
		fmt.Println("Tests missing in QA.md:")
		for _, name := range missingInQA {
			fmt.Printf("  - %s\n", name)
		}
		hasErrors = true
	}

	if len(missingInE2E) > 0 {
		fmt.Println("\nTests in QA.md but not in E2E tests:")
		for _, name := range missingInE2E {
			fmt.Printf("  - %s\n", name)
		}
		hasErrors = true
	}

	if hasErrors {
		fmt.Println("\nPlease update QA.md to match the E2E test files.")
		os.Exit(1)
	}

	fmt.Println("QA.md is in sync with E2E tests.")
}

func extractTestNames(dir string) ([]string, error) {
	var testNames []string
	testFuncRegex := regexp.MustCompile(`^func\s+(Test\w+)\s*\(`)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, "_test.go") {
			return nil
		}

		file, err := os.Open(path) //nolint:gosec // G304: path is from trusted source
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			matches := testFuncRegex.FindStringSubmatch(line)
			if len(matches) > 1 {
				testName := matches[1]
				if testName != "TestMain" {
					testNames = append(testNames, testName)
				}
			}
		}

		return scanner.Err()
	})

	return testNames, err
}

func extractQATestNames(qaFile string) ([]string, error) {
	var testNames []string
	testHeaderRegex := regexp.MustCompile(`^###\s+(Test\w+)\s*$`)

	file, err := os.Open(qaFile) //nolint:gosec // G304: path is from trusted source
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := testHeaderRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			testNames = append(testNames, matches[1])
		}
	}

	return testNames, scanner.Err()
}

func findMissing(source, target []string) []string {
	targetSet := make(map[string]bool)
	for _, name := range target {
		targetSet[name] = true
	}

	var missing []string
	for _, name := range source {
		if !targetSet[name] {
			missing = append(missing, name)
		}
	}

	return missing
}
