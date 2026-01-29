package markdown

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ageha734/mmdc/domain"
)

type MermaidMatch struct {
	FullMatch string
	Code      string
	Index     int
}

var (
	mermaidCodeBlockRegex = regexp.MustCompile("(?s)```mermaid\\s*\n(.*?)```")
	mermaidFenceRegex     = regexp.MustCompile("(?s):::mermaid\\s*\n(.*?):::")
	mermaidReplaceRegex   = regexp.MustCompile("(?s)```mermaid\\s*\n.*?```|:::mermaid\\s*\n.*?:::")
)

func ExtractAllMermaidFromMarkdown(content string) []MermaidMatch {
	var matches []MermaidMatch

	matches = append(matches, extractMatches(mermaidCodeBlockRegex, content, 0)...)
	codeBlockCount := len(matches)
	matches = append(matches, extractMatches(mermaidFenceRegex, content, codeBlockCount)...)

	return matches
}

func extractMatches(re *regexp.Regexp, content string, startIndex int) []MermaidMatch {
	var matches []MermaidMatch
	allMatches := re.FindAllStringSubmatch(content, -1)
	for i, match := range allMatches {
		if len(match) >= 2 {
			matches = append(matches, MermaidMatch{
				FullMatch: match[0],
				Code:      match[1],
				Index:     startIndex + i,
			})
		}
	}
	return matches
}

func GenerateOutputFileName(outputPath string, index int, format domain.OutputFormat) string {
	absOutputPath, err := filepath.Abs(outputPath)
	if err != nil {
		absOutputPath = outputPath
	}
	baseName := strings.TrimSuffix(filepath.Base(absOutputPath), filepath.Ext(absOutputPath))
	ext := filepath.Ext(absOutputPath)
	if ext == ".md" || ext == ".markdown" {
		ext = "." + format.String()
	}
	return filepath.Join(filepath.Dir(absOutputPath), fmt.Sprintf("%s-%d%s", baseName, index+1, ext))
}

func CalculateRelativePath(outputPath, outputFile string) string {
	relPath, _ := filepath.Rel(filepath.Dir(outputPath), outputFile)
	if !strings.HasPrefix(relPath, ".") {
		relPath = "./" + relPath
	}
	return relPath
}

func ReplaceMermaidWithImages(content string, imageInfos []domain.ImageInfo) string {
	imageIndex := 0
	return mermaidReplaceRegex.ReplaceAllStringFunc(content, func(match string) string {
		if imageIndex < len(imageInfos) {
			img := imageInfos[imageIndex]
			imageIndex++
			return fmt.Sprintf("![%s](%s)", img.Alt, img.URL)
		}
		return match
	})
}
