# Text Extraction Example

This example demonstrates how to extract text with positional information from PDF documents using GoPDF.

## Features

- Extract text from PDF pages
- Get X, Y coordinates for each text element
- Font information (name and size)
- Text bounding boxes
- Statistics (font usage, page coverage)

## Usage

```bash
go run main.go <pdf-file>
```

## Example Output

```
Opening PDF: sample.pdf
PDF has 1 pages

Extracting text from page 1...

Found 25 text elements:

Position details for first 20 elements:
-------------------------------------------------------
[1] Text: "Hello, World!"
    Position: (100.00, 700.00)
    Size: 72.00 x 12.00
    Font: Helvetica, Size: 12.0pt

[2] Text: "This is a sample PDF."
    Position: (100.00, 686.00)
    Size: 108.00 x 12.00
    Font: Helvetica, Size: 12.0pt

...

-------------------------------------------------------
All extracted text:
-------------------------------------------------------
Hello, World!
This is a sample PDF.
...

-------------------------------------------------------

Statistics:
  Total text elements: 25
  Unique fonts: 2
    Helvetica: 20 elements
    Helvetica-Bold: 5 elements
  Text bounding box:
    Bottom-left: (72.00, 100.00)
    Top-right: (540.00, 720.00)
    Dimensions: 468.00 x 620.00 points
```

## Phase 2.5 Implementation

This example demonstrates Phase 2.5 of the GoPDF roadmap: **Text Extraction with Positional Information**.

### What's Implemented:

1. **Text Elements** - Each piece of text with position data
2. **Matrix Transformations** - PDF text matrix handling
3. **Text State Tracking** - Font, size, spacing, leading
4. **Content Stream Parsing** - Parse PDF operators
5. **Text Operators** - BT, ET, Tj, TJ, Tm, Td, Tf, etc.
6. **FlateDecode Support** - Decompress compressed streams

### Critical for Table Extraction:

The positional information extracted here (X, Y coordinates) is **essential** for Phase 2.6 (Table Detection).
Table extraction algorithms need to know where text is located to determine:

- Column boundaries (vertical alignment)
- Row boundaries (horizontal alignment)
- Cell grouping (proximity analysis)
- Table regions (spatial clustering)

## Implementation Details

### Text Extraction Process:

1. Open PDF file with `parser.OpenPDF()`
2. Create `extractor.NewTextExtractor(reader)`
3. Call `ExtractFromPage(pageNum)` to get text elements
4. Process elements with position data

### Text Operators Supported:

- **BT/ET** - Text object delimiters
- **Tf** - Set font and size
- **Tm** - Set text matrix
- **Td/TD** - Move text position
- **T*** - Move to next line
- **Tj** - Show text string
- **TJ** - Show text with positioning
- **Tc/Tw/Tz/TL/Ts** - Text state parameters

### Stream Decoding:

- **FlateDecode** - zlib decompression (most common)
- Other filters can be added in future phases

## Next Steps

After text extraction, the next phase is:

**Phase 2.6 - Table Detection**:
- Use text positions to detect table regions
- Identify ruling lines (lattice mode)
- Analyze whitespace (stream mode)
- Group text into cells
- Build table structure

## Testing

Test with your own PDF files:

```bash
# Simple text document
go run main.go document.pdf

# Invoice with tables
go run main.go invoice.pdf

# Financial report
go run main.go report.pdf
```

## Requirements

- Go 1.25+
- PDF file (unencrypted)

## Limitations (Phase 2.5)

- Font width estimation is approximate (full font metrics in Phase 3)
- Encryption not supported (Phase 4)
- Some advanced encoders not supported (Phase 4)

## Reference

- PDF 1.7 Specification, Section 9.4 (Text Objects)
- PDF 1.7 Specification, Section 9.3 (Text State)
- GoPDF ROADMAP.md, Phase 2.5

## License

MIT License - See LICENSE in project root
