# GxPDF Architecture

Technical architecture overview of the GxPDF PDF library.

## Project Structure

```
github.com/coregx/gxpdf
├── gxpdf.go              # Main public API
├── creator/              # PDF creation API
│   └── forms/            # Interactive form fields
├── export/               # Export formats (CSV, JSON, Excel)
└── internal/             # Private implementation
    ├── application/      # Use cases
    │   ├── extractor/    # Text extraction
    │   ├── detector/     # Table detection
    │   ├── table/        # Table extraction
    │   └── reader/       # PDF reading
    └── infrastructure/   # Technical implementation
        ├── parser/       # PDF parsing
        ├── primitive/    # PDF objects
        ├── encoding/     # Stream codecs
        ├── fonts/        # Font handling
        ├── security/     # Encryption
        └── writer/       # PDF writing
```

## Key Components

### Public API (`/`, `/creator`, `/export`)

User-facing API designed for simplicity:

```go
// Main API
doc, _ := gxpdf.Open("document.pdf")
tables := doc.ExtractTables()

// Creator API
c := creator.New()
page, _ := c.NewPage()
page.AddText("Hello", 100, 700, creator.Helvetica, 12)
c.WriteToFile("output.pdf")

// Export API
exporter := export.NewCSVExporter()
exporter.Export(table, writer)
```

### Parser (`internal/infrastructure/parser`)

Handles PDF file parsing:

- **Lexer** - Tokenizes PDF byte stream
- **Syntax Parser** - Parses PDF objects
- **XRef Parser** - Cross-reference table handling
- **Reader** - Document navigation

### Primitive Objects (`internal/infrastructure/primitive`)

PDF object types:

- `Null`, `Boolean`, `Integer`, `Real`
- `String` (literal and hex)
- `Name`, `Array`, `Dictionary`
- `Stream`, `Reference`

### Encoding (`internal/infrastructure/encoding`)

Stream compression/decompression:

- FlateDecode (zlib)
- ASCII85Decode
- ASCIIHexDecode
- LZWDecode
- DCTDecode (JPEG)

### Extraction (`internal/application`)

Text and table extraction:

- **Text Extractor** - Extract text with positions
- **Table Detector** - Detect table regions
- **Table Extractor** - Extract structured data

## Design Principles

### 1. Simple Public API

The public API hides complexity:

```go
// User sees simple API
doc, _ := gxpdf.Open("file.pdf")
text := doc.Page(0).Text()

// Internal complexity hidden
```

### 2. Internal Privacy

`internal/` enforces API boundaries:

- External code cannot import `internal/`
- Free to refactor without breaking users
- Clear distinction between public and private

### 3. Functional Options

Configuration through options pattern:

```go
reader, _ := parser.OpenPDF(path,
    parser.WithPassword("secret"),
    parser.WithStrictMode(true),
)
```

### 4. Error Handling

Errors with context:

```go
if err != nil {
    return fmt.Errorf("parse xref at %d: %w", offset, err)
}
```

## Testing

- **Unit Tests** - Test individual components
- **Integration Tests** - Test with real PDFs
- **Race Detector** - Verify thread safety

```bash
go test ./...
go test -race ./...
```

## Dependencies

Minimal external dependencies:

- Production: Standard library only
- Testing: `testify`

## Future Improvements

Planned restructuring to feature-based organization:

```
internal/
├── parser/       # All PDF parsing
├── creator/      # All PDF creation
├── extractor/    # All extraction
├── encoding/     # Stream codecs
└── security/     # Encryption
```

This will flatten the hierarchy and organize by feature rather than technical layer.
