package renderer

import (
	"context"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
)

type MockRenderer struct {
	RenderFunc func(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error)
}

func (m *MockRenderer) Render(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error) {
	if m.RenderFunc != nil {
		return m.RenderFunc(ctx, diagram, format, renderConfig)
	}
	switch format {
	case domain.OutputFormatSVG:
		return []byte("<svg>mock rendered</svg>"), nil
	case domain.OutputFormatPNG:
		return []byte("mock png data"), nil
	case domain.OutputFormatPDF:
		return []byte("mock pdf data"), nil
	default:
		return nil, nil
	}
}

func NewMockRenderer() *MockRenderer {
	return &MockRenderer{}
}
