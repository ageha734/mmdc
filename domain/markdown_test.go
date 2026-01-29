package domain_test

import (
	"testing"

	"github.com/ageha734/mmdc/domain"
)

func TestNewMarkdownFile(t *testing.T) {
	diagrams := []domain.MermaidDiagram{
		{Code: "graph TD\nA-->B", Index: 0},
		{Code: "sequenceDiagram\nA->>B: Hello", Index: 1},
	}
	mdFile := domain.NewMarkdownFile("test.md", "content", diagrams)
	if mdFile == nil {
		t.Fatal("NewMarkdownFile returned nil")
	}
	if mdFile.Path != "test.md" {
		t.Errorf("expected Path 'test.md', got %q", mdFile.Path)
	}
	if mdFile.Content != "content" {
		t.Errorf("expected Content 'content', got %q", mdFile.Content)
	}
	if len(mdFile.Diagrams) != 2 {
		t.Errorf("expected 2 diagrams, got %d", len(mdFile.Diagrams))
	}
}

func TestImageInfo(t *testing.T) {
	img := domain.ImageInfo{
		URL:   "./diagram-1.svg",
		Title: "Diagram 1",
		Alt:   "diagram 1",
	}
	if img.URL != "./diagram-1.svg" {
		t.Errorf("expected URL './diagram-1.svg', got %q", img.URL)
	}
	if img.Title != "Diagram 1" {
		t.Errorf("expected Title 'Diagram 1', got %q", img.Title)
	}
	if img.Alt != "diagram 1" {
		t.Errorf("expected Alt 'diagram 1', got %q", img.Alt)
	}
}
