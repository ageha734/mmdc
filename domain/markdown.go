// Package domain defines core domain types for mmdc.
package domain

type ImageInfo struct {
	URL   string
	Title string
	Alt   string
}

type MarkdownFile struct {
	Path     string
	Content  string
	Diagrams []MermaidDiagram
}

func NewMarkdownFile(path, content string, diagrams []MermaidDiagram) *MarkdownFile {
	return &MarkdownFile{
		Path:     path,
		Content:  content,
		Diagrams: diagrams,
	}
}
