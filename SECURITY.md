# Security Policy

## Supported Versions

| Version | Supported          |
| ------- | ------------------ |
| 0.1.x   | :white_check_mark: |

Security updates are provided for the latest minor version.

## Reporting a Vulnerability

We take security seriously. If you discover a security vulnerability, please report it responsibly.

### How to Report

**DO NOT** open a public GitHub issue for security vulnerabilities.

Instead, report security issues via:

1. **GitHub Security Advisory** (preferred):
   https://github.com/coregx/gxpdf/security/advisories/new

2. **Email**: Create a private GitHub issue or contact via discussions

### What to Include

- Description of the vulnerability
- Steps to reproduce (include malicious PDF if applicable)
- Affected versions
- Potential impact (RCE, DoS, memory exhaustion, etc.)
- Suggested fix (if available)

### Response Timeline

- **Initial Response**: Within 72 hours
- **Assessment**: Within 1 week
- **Fix & Disclosure**: Coordinated with reporter

## Security Considerations

GxPDF parses untrusted PDF files. Users should be aware of these risks:

### Malicious PDF Files

**Risk**: Crafted PDFs can cause crashes, memory exhaustion, or unexpected behavior.

**Mitigations**:
- Recursion depth limits
- Cycle detection in object references
- Size limits on decompression
- Input validation

**User Recommendations**:
```go
// Use timeout for untrusted PDFs
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

reader := pdf.NewReader("untrusted.pdf")
doc, err := reader.ReadWithContext(ctx)
```

### Stream Decompression (Zip Bombs)

**Risk**: Compressed streams can expand to large sizes.

**Mitigations**:
- Maximum decompressed size limits
- Compression ratio limits
- Streaming decompression

### Integer Overflow

**Risk**: Large values can cause integer overflow.

**Mitigations**:
- Safe integer arithmetic
- Bounds checking
- Validation of object numbers

## Best Practices for Users

### Input Validation

```go
// Validate file size
fi, _ := os.Stat(pdfPath)
if fi.Size() > maxPDFSize {
    return errors.New("PDF too large")
}

// Use timeout
ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
defer cancel()
```

### Resource Limits

```go
// Limit pages
const maxPages = 100
pages := min(doc.PageCount(), maxPages)

// Limit content size
const maxTextSize = 10 * 1024 * 1024
text, _ := page.ExtractText()
if len(text) > maxTextSize {
    return errors.New("text too large")
}
```

### Error Handling

```go
// Always check errors
doc, err := reader.Read()
if err != nil {
    return fmt.Errorf("PDF parsing failed: %w", err)
}
```

## Security Testing

- Unit tests with edge cases
- Fuzz testing for parser
- Race detector (0 data races)
- Static analysis with gosec
- golangci-lint security checks

## Dependencies

GxPDF minimizes external dependencies:

- Standard library for production code
- `testify` for testing only
- Dependabot enabled for updates

## Security Contact

- **GitHub Security Advisory**: https://github.com/coregx/gxpdf/security/advisories/new
- **Issues** (non-sensitive): https://github.com/coregx/gxpdf/issues

---

**Thank you for helping keep GxPDF secure!**
