# Contributing to GxPDF

Thank you for considering contributing to GxPDF!

The following guidelines help maintain code quality and consistency. Use your best judgment, and feel free to propose changes to this document.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Coding Standards](#coding-standards)
- [Architecture Guidelines](#architecture-guidelines)

## Code of Conduct

This project follows a professional Code of Conduct. By participating, you agree to:

- Be respectful and professional
- Provide constructive feedback
- Focus on the code, not the person

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gxpdf.git
   cd gxpdf
   ```
3. **Add upstream remote**:
   ```bash
   git remote add upstream https://github.com/coregx/gxpdf.git
   ```

## Development Setup

### Prerequisites

- **Go 1.25 or later** (required)
- **Git**
- **golangci-lint** (recommended)

### Install Tools

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Download dependencies
go mod download
```

### Verify Setup

```bash
go test ./...
golangci-lint run
```

## Making Changes

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

**Branch naming conventions**:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Adding tests

### 2. Make Your Changes

- Follow the [Coding Standards](#coding-standards)
- Follow the [Architecture Guidelines](#architecture-guidelines)
- Write tests for new functionality
- Update documentation as needed

### 3. Run Checks

Before committing:

```bash
go fmt ./...
go test ./...
go test -race ./...
golangci-lint run
```

## Testing

### Writing Tests

Use **table-driven tests** (Go best practice):

```go
func TestRectangle_Dimensions(t *testing.T) {
    tests := []struct {
        name  string
        rect  Rectangle
        wantW float64
        wantH float64
    }{
        {
            name:  "A4 size",
            rect:  MustRectangle(0, 0, 595, 842),
            wantW: 595,
            wantH: 842,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.wantW, tt.rect.Width())
            assert.Equal(t, tt.wantH, tt.rect.Height())
        })
    }
}
```

### Running Tests

```bash
# All tests
go test ./...

# With coverage
go test -cover ./...

# With race detector
go test -race ./...

# Specific package
go test ./internal/infrastructure/parser/...
```

## Submitting Changes

### 1. Commit Your Changes

Follow **Conventional Commits** format:

```
<type>(<scope>): <subject>

<body>
```

**Types**:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation only
- `style:` - Code style (formatting)
- `refactor:` - Code refactoring
- `test:` - Adding tests
- `chore:` - Maintenance

**Examples**:
```bash
git commit -m "feat(parser): add support for object streams"
git commit -m "fix(writer): handle empty pages correctly"
git commit -m "docs(readme): update installation instructions"
```

### 2. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub with:
- Clear description of changes
- Link to related issues
- Test results

## Coding Standards

### Go Style

Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).

### Naming Conventions

```go
// Types (PascalCase)
type Document struct { ... }
type PdfReader struct { ... }

// Interfaces (-er suffix when possible)
type Parser interface { ... }
type Encoder interface { ... }

// Private fields (camelCase)
type Document struct {
    id      DocumentID
    version Version
}
```

### Error Handling

```go
// GOOD: Wrap errors with context
if err := parser.Parse(); err != nil {
    return fmt.Errorf("parse PDF at offset %d: %w", offset, err)
}

// BAD: Lose context
if err := parser.Parse(); err != nil {
    return err
}
```

### Comments

```go
// GOOD: Document exported functions
// Parse parses a PDF file from the given reader.
// It returns a Document or an error if parsing fails.
func Parse(r io.Reader) (*Document, error) { ... }
```

## Architecture Guidelines

GxPDF follows **Domain-Driven Design (DDD)** principles.

### Layer Structure

```
internal/
├── domain/              # Pure business logic (NO external deps)
├── infrastructure/      # Technical implementation
└── application/         # Use cases (orchestrates domain)
```

### Dependency Rules

1. **domain/** → NO dependencies
2. **application/** → depends on **domain/**
3. **infrastructure/** → depends on **domain/**
4. **pkg/** → depends on **application/** + **domain/**

### Rich Domain Model

**Prefer behavior over data**:

```go
// BAD: Anemic model
type Page struct {
    Width  float64
    Height float64
}

func (p *Page) GetWidth() float64 { return p.Width }

// GOOD: Rich model
type Page struct {
    dimensions Rectangle
    content    ContentStream
}

func (p *Page) AddText(text string, pos Position, font *Font) error {
    return p.content.AppendText(text, pos, font)
}
```

## What to Contribute

### Good First Issues

Look for issues labeled `good first issue`:
- Documentation improvements
- Adding tests
- Small bug fixes
- Adding examples

### Areas That Need Help

- Parser improvements
- New stream encoders
- Documentation
- Test coverage
- Performance optimization

## Questions?

- **Issues**: [GitHub Issues](https://github.com/coregx/gxpdf/issues)
- **Discussions**: [GitHub Discussions](https://github.com/coregx/gxpdf/discussions)

## License

By contributing to GxPDF, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to GxPDF!**
