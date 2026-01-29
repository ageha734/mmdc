# mmdc

[![CI](https://github.com/ageha734/mmdc/actions/workflows/ci.yml/badge.svg)](https://github.com/ageha734/mmdc/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/ageha734/mmdc)](https://goreportcard.com/report/github.com/ageha734/mmdc)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/ageha734/mmdc.svg)](https://pkg.go.dev/github.com/ageha734/mmdc)

[日本語](README.ja.md)

A Go-based CLI tool for rendering Mermaid diagrams to SVG, PNG, and PDF formats. This tool uses chromedp for headless browser rendering, providing a native Go alternative to the official [@mermaid-js/mermaid-cli](https://github.com/mermaid-js/mermaid-cli).

## Features

- 🎨 **Multiple Output Formats**: SVG, PNG, and PDF
- 📝 **Markdown Support**: Extract and render embedded Mermaid diagrams from Markdown files
- 🎭 **Theme Support**: default, forest, dark, and neutral themes
- ⚡ **Fast Rendering**: Uses chromedp for efficient headless browser rendering
- 📦 **Single Binary**: No Node.js or npm required
- 🔧 **Highly Configurable**: Customize dimensions, colors, scale, and more

## Installation

### Using Go

```bash
go install github.com/ageha734/mmdc/cmd@latest
```

### Using Homebrew (macOS/Linux)

```bash
brew tap ageha734/mmdc
brew install mmdc
```

### Using curl (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/ageha734/mmdc/master/install.sh | bash
```

### Using PowerShell (Windows)

```powershell
irm https://raw.githubusercontent.com/ageha734/mmdc/master/install.ps1 | iex
```

### Download Binary

Download the latest binary from the [Releases](https://github.com/ageha734/mmdc/releases) page.

## Quick Start

### Render a Mermaid file to SVG

```bash
mmdc -i diagram.mmd -o diagram.svg
```

### Render to PNG with custom dimensions

```bash
mmdc -i diagram.mmd -o diagram.png -w 1200 -H 800
```

### Render Markdown with embedded diagrams

```bash
mmdc -i document.md -a ./images/
```

### Read from stdin

```bash
echo "graph TD; A-->B" | mmdc -i - -o output.svg
```

## Command Reference

```text
mmdc - Mermaid CLI for Go

Usage:
  mmdc [flags]
  mmdc [command]

Flags:
  -i, --input string              Input mermaid file or markdown file
  -o, --output string             Output file path (default: input + .svg)
  -e, --outputFormat string       Output format: svg, png, or pdf
  -t, --theme string              Theme: default, forest, dark, or neutral (default "default")
  -w, --width int                 Page width (default 800)
  -H, --height int                Page height (default 600)
  -b, --backgroundColor string    Background color (default "white")
  -s, --scale int                 Scale factor (default 1)
  -q, --quiet                     Suppress log output
  -c, --configFile string         Mermaid config JSON file
  -C, --cssFile string            CSS file for custom styling
  -I, --svgId string              SVG element ID
  -f, --pdfFit                    Scale PDF to fit chart
  -a, --artefacts string          Output artefacts path (for Markdown input)
  -p, --puppeteerConfigFile       Puppeteer config JSON file
      --iconPacks strings         Icon packs to use
  -v, --version                   Version information

Commands:
  render      Render Mermaid diagram
  help        Help about any command
```

## Examples

### Basic Usage

```bash
# Convert Mermaid to SVG
mmdc -i flowchart.mmd

# Convert to PNG with dark theme
mmdc -i sequence.mmd -o sequence.png -t dark

# Convert to PDF
mmdc -i diagram.mmd -o diagram.pdf -e pdf
```

### Markdown Processing

```bash
# Extract diagrams from Markdown and save to images directory
mmdc -i README.md -a ./images/

# Process Markdown with custom output format
mmdc -i docs.md -a ./diagrams/ -e png
```

### Advanced Configuration

```bash
# Use custom Mermaid configuration
mmdc -i diagram.mmd -c mermaid-config.json

# Apply custom CSS styling
mmdc -i diagram.mmd -C custom.css -o styled.svg

# High-resolution PNG output
mmdc -i diagram.mmd -o hires.png -s 2 -w 1600 -H 1200
```

### Pipeline Usage

```bash
# Generate diagram from script output
cat <<EOF | mmdc -i - -o pipeline.svg
graph LR
    A[Start] --> B{Decision}
    B -->|Yes| C[Process]
    B -->|No| D[End]
EOF
```

## Configuration File

You can use a JSON configuration file for Mermaid settings:

```json
{
  "theme": "forest",
  "themeVariables": {
    "primaryColor": "#ff6b6b",
    "secondaryColor": "#4ecdc4"
  },
  "flowchart": {
    "curve": "basis"
  },
  "sequence": {
    "mirrorActors": false
  }
}
```

Use with:

```bash
mmdc -i diagram.mmd -c config.json
```

## Requirements

- Go 1.23.2 or later (for building from source)
- Chrome/Chromium browser (automatically detected)

## Architecture

```text
mmdc/
├── cmd/              # CLI entry point (Cobra)
├── config/           # Configuration constants
├── domain/           # Domain models (MermaidDiagram, OutputFormat)
├── handler/          # Browser and file handling
├── infra/            # Infrastructure (renderer, logger, file I/O)
├── pkg/              # Utility packages (markdown, validation, path)
├── test/mock/        # Test mocks
└── usecase/          # Business logic
```

## Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) before submitting a pull request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Security

For security issues, please see our [Security Policy](SECURITY.md).

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.

## Acknowledgments

- [mermaid-js/mermaid](https://github.com/mermaid-js/mermaid) - The amazing diagramming library
- [chromedp/chromedp](https://github.com/chromedp/chromedp) - Chrome DevTools Protocol for Go
- [spf13/cobra](https://github.com/spf13/cobra) - CLI framework

## Related Projects

- [@mermaid-js/mermaid-cli](https://github.com/mermaid-js/mermaid-cli) - Official Node.js-based Mermaid CLI
