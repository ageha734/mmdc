# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [1.0.0] - 2024-XX-XX

### Added

- Initial release of mmdc
- Support for rendering Mermaid diagrams to SVG, PNG, and PDF formats
- Markdown file processing with embedded Mermaid diagram extraction
- Theme support: default, forest, dark, and neutral
- Custom CSS styling support
- Mermaid configuration file support
- stdin/stdout support for pipeline integration
- Command-line options:
  - `-i, --input`: Input file path
  - `-o, --output`: Output file path
  - `-e, --outputFormat`: Output format selection
  - `-t, --theme`: Theme selection
  - `-w, --width`: Page width
  - `-H, --height`: Page height
  - `-b, --backgroundColor`: Background color
  - `-s, --scale`: Scale factor
  - `-q, --quiet`: Suppress log output
  - `-c, --configFile`: Mermaid config JSON file
  - `-C, --cssFile`: Custom CSS file
  - `-I, --svgId`: SVG element ID
  - `-f, --pdfFit`: PDF fit to chart
  - `-a, --artefacts`: Artefacts output path
- Cross-platform support: Linux, macOS, Windows
- Headless Chrome rendering via chromedp
- Comprehensive error handling and validation
- Color-coded console output
- CI/CD workflows with GitHub Actions
- Security scanning with CodeQL, govulncheck, Trivy, and gosec
- Dependabot configuration for dependency updates

### Technical Details

- Built with Go 1.23.2
- Uses chromedp for headless browser automation
- Uses Cobra for CLI framework
- Clean architecture with separated concerns:
  - `cmd/`: CLI entry point
  - `config/`: Configuration constants
  - `domain/`: Domain models
  - `handler/`: Browser and file handling
  - `infra/`: Infrastructure (renderer, logger, file I/O)
  - `pkg/`: Utility packages
  - `usecase/`: Business logic

[Unreleased]: https://github.com/ageha734/mmdc/compare/v1.0.0...HEAD
[1.0.0]: https://github.com/ageha734/mmdc/releases/tag/v1.0.0
