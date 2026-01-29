package renderer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ageha734/mmdc/config"
)

type HTMLGenerator struct {
	renderConfig *config.RenderConfig
}

func NewHTMLGenerator(renderConfig *config.RenderConfig) *HTMLGenerator {
	return &HTMLGenerator{renderConfig: renderConfig}
}

func (g *HTMLGenerator) Generate(mermaidCode, svgId string) (string, error) {
	mermaidConfig, err := config.LoadMermaidConfig(g.renderConfig)
	if err != nil {
		return "", err
	}

	configJSON, _ := json.Marshal(mermaidConfig)
	svgId = config.DetermineSVGID(svgId, g.renderConfig)
	cssContent := g.getCSS()

	return g.buildHTMLTemplate(g.renderConfig.BackgroundColor, cssContent, string(configJSON), mermaidCode, svgId), nil
}

func (g *HTMLGenerator) getCSS() string {
	if g.renderConfig.CSSFile == "" {
		return ""
	}
	data, err := os.ReadFile(g.renderConfig.CSSFile)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("<style>%s</style>", string(data))
}

func (g *HTMLGenerator) buildHTMLTemplate(backgroundColor, cssContent, configJSON, mermaidCode, svgId string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<script src="https://cdn.jsdelivr.net/npm/mermaid@11.12.0/dist/mermaid.min.js"></script>
	<style>
		body {
			margin: 0;
			padding: 20px;
			background-color: %s;
		}
		#container {
			display: flex;
			justify-content: center;
			align-items: center;
		}
		%s
	</style>
</head>
<body>
	<div id="container"></div>
	<script>
		(function() {
			var maxWaitTime = 30000;
			var startTime = Date.now();

			function waitForMermaid() {
				if (typeof mermaid !== 'undefined') {
					renderMermaid();
				} else if (Date.now() - startTime < maxWaitTime) {
					setTimeout(waitForMermaid, 100);
				} else {
					console.error('Timeout waiting for mermaid to load');
					document.getElementById('container').setAttribute('data-rendered', 'error');
					document.getElementById('container').innerHTML = '<p style="color: red;">Error: Mermaid script failed to load</p>';
				}
			}

			function renderMermaid() {
				try {
					mermaid.initialize(%s);
					const definition = %s;
					const svgId = %s;
					const container = document.getElementById('container');

					container.setAttribute('data-rendered', 'rendering');

					mermaid.render(svgId, definition).then(function(result) {
						container.innerHTML = result.svg;
						container.setAttribute('data-rendered', 'true');
					}).catch(function(err) {
						console.error('Mermaid render error:', err);
						var errorMsg = err.message || 'Unknown error';
						container.innerHTML = '<p style="color: red;">Error: ' + errorMsg + '</p>';
						container.setAttribute('data-rendered', 'error');
					});
				} catch (err) {
					console.error('Error:', err);
					var errorMsg = err.message || 'Unknown error';
					document.getElementById('container').innerHTML = '<p style="color: red;">Error: ' + errorMsg + '</p>';
					document.getElementById('container').setAttribute('data-rendered', 'error');
				}
			}

			waitForMermaid();
		})();
	</script>
</body>
</html>`, backgroundColor, cssContent, configJSON, g.escapeJS(mermaidCode), g.escapeJS(svgId))
}

func (g *HTMLGenerator) escapeJS(s string) string {
	if s == "" {
		return `""`
	}
	b, _ := json.Marshal(s)
	return string(b)
}
