# Contributing to mmdc

Thank you for your interest in contributing to mmdc! This document provides guidelines and instructions for contributing.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Submitting Changes](#submitting-changes)
- [Style Guidelines](#style-guidelines)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

## Getting Started

### Prerequisites

- Go 1.23.2 or later
- Chrome or Chromium browser
- Git

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

   ```bash
   git clone https://github.com/YOUR_USERNAME/mmdc.git
   cd mmdc
   ```

3. Add the upstream remote:

   ```bash
   git remote add upstream https://github.com/ageha734/mmdc.git
   ```

## Development Setup

### Install Dependencies

```bash
go mod download
```

### Build

```bash
go build -o mmdc ./cmd
```

### Run Tests

```bash
go test -v ./...
```

### Install Development Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install mockgen (for generating mocks)
go install go.uber.org/mock/mockgen@latest
```

## Making Changes

### Branch Naming

Create a branch for your changes:

```bash
git checkout -b <type>/<description>
```

Types:

- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Adding or updating tests
- `chore/` - Maintenance tasks

Examples:

- `feature/add-webp-support`
- `fix/markdown-extraction-bug`
- `docs/improve-readme`

### Commit Messages

We follow [Conventional Commits](https://www.conventionalcommits.org/):

```text
<type>(<scope>): <description>
```

Types:

- `feat`: New feature
- `fix`: Bug fix
- `chore`: Maintenance

Examples:

```text
feat(renderer): add WebP output format support
fix(markdown): handle nested code blocks correctly
```

## Submitting Changes

### Pull Request Process

1. Update your branch with the latest upstream changes:

   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. Run all tests and linters:

   ```bash
   go test -v ./...
   golangci-lint run
   ```

3. Push your changes:

   ```bash
   git push origin <your-branch>
   ```

4. Open a Pull Request on GitHub
5. Fill out the PR template completely

### PR Requirements

- [ ] All tests pass
- [ ] Code follows style guidelines
- [ ] Documentation is updated (if applicable)
- [ ] CHANGELOG.md is updated
- [ ] Commit messages follow conventional commits
- [ ] No merge conflicts

## Style Guidelines

### Go Code Style

- Follow the official [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting (automatically run by golangci-lint)
- Use `golangci-lint` for linting:

```bash
golangci-lint run
```

### Code Organization

```text
mmdc/
├── cmd/              # CLI entry point
├── config/           # Configuration
├── domain/           # Domain models
├── handler/          # HTTP/IO handlers
├── infra/            # Infrastructure (external services)
├── pkg/              # Reusable packages
├── test/mock/        # Test mocks
└── usecase/          # Business logic
```

### Naming Conventions

- Use descriptive names
- Acronyms should be consistently cased (e.g., `HTML`, `URL`, `ID`)
- Interface names should describe behavior
- Avoid single-letter variable names except in short loops

### Error Handling

- Always handle errors explicitly
- Use custom error types for domain errors
- Wrap errors with context using `fmt.Errorf("...: %w", err)`

```go
// Good
if err != nil {
    return fmt.Errorf("failed to render diagram: %w", err)
}

// Avoid
if err != nil {
    return err
}
```

## Testing

### Test Structure

- Place tests in the same package as the code being tested
- Use table-driven tests for multiple test cases
- Use meaningful test names

```go
func TestMermaidDiagram_OutputFormat(t *testing.T) {
    tests := []struct {
        name    string
        format  string
        want    OutputFormat
        wantErr bool
    }{
        {"valid svg", "svg", OutputFormatSVG, false},
        {"valid png", "png", OutputFormatPNG, false},
        {"invalid format", "invalid", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // test implementation
        })
    }
}
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -v -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package tests
go test -v ./pkg/markdown/...

# Run tests with race detection
go test -race ./...
```

### Mocks

Generate mocks using mockgen:

```bash
go generate ./...
```

## Documentation

### Code Documentation

- Document all exported types, functions, and methods
- Use complete sentences in comments
- Include examples for complex functionality

```go
// RenderDiagram renders a Mermaid diagram to the specified output format.
// It uses a headless Chrome browser to render the diagram.
//
// Example:
//
//     output, err := RenderDiagram(ctx, "graph TD; A-->B", FormatSVG)
func RenderDiagram(ctx context.Context, code string, format OutputFormat) ([]byte, error)
```

### README Updates

When adding new features:

1. Update the feature list
2. Add usage examples
3. Update the command reference if flags change

### CHANGELOG Updates

Add entries under the "Unreleased" section:

```markdown
## [Unreleased]

### Added
- New feature description

### Changed
- Modified behavior description

### Fixed
- Bug fix description
```

## Questions?

If you have questions about contributing, please:

1. Check existing issues and discussions
2. Open a new discussion or issue
3. Reach out to maintainers

Thank you for contributing!
