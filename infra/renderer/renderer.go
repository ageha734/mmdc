package renderer

import (
	"context"

	"github.com/ageha734/mmdc/config"
	"github.com/ageha734/mmdc/domain"
)

type Renderer interface {
	Render(ctx context.Context, diagram *domain.MermaidDiagram, format domain.OutputFormat, renderConfig *config.RenderConfig) ([]byte, error)
}

type RenderOptions struct {
	Quiet bool
}
