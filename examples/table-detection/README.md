# Table Detection Example

This example demonstrates Phase 2.6 of the GoPDF project - **Table Detection** using algorithms inspired by tabula-java.

## Overview

The table detection system can identify table regions in PDF documents using two modes:

1. **Lattice Mode**: Detects tables with visible ruling lines (borders and grid lines)
2. **Stream Mode**: Detects tables without visible lines using whitespace analysis
3. **Auto Mode**: Automatically selects the best mode for each page

## Features

- Graphics operator parsing to extract lines and rectangles
- Ruling line detection and merging for lattice mode
- Projection profile analysis for stream mode
- Whitespace analysis to detect column and row boundaries
- Grid building from intersecting ruling lines
- Auto-detection of best extraction strategy

## Usage

```bash
# Build the example
go build -o table-detection

# Run on a PDF file
./table-detection document.pdf
```

## Example Output

```
Opening PDF: invoice.pdf
PDF has 3 pages

=== Page 1 ===
Extracting text elements...
Found 145 text elements
Extracting graphics elements...
Found 24 graphics elements
Detecting tables...
Found 2 table(s):

Table 1:
  Method: Lattice
  Bounds: Rectangle{x=50.00, y=100.00, w=500.00, h=300.00}
  Rows: 8
  Columns: 5
  Grid structure:
    Row coordinates: [100 125 150 175 200 225 250 275 300]
    Column coordinates: [50 150 250 350 450 550]

Table 2:
  Method: Stream
  Bounds: Rectangle{x=50.00, y=450.00, w=500.00, h=150.00}
  Rows: 4
  Columns: 3
  Stream mode structure:
    Row coordinates: [450 485 520 555 600]
    Column coordinates: [50 200 350 550]

Auto-detected mode:
  Recommended mode: Lattice
```

## Algorithm Details

### Lattice Mode (With Ruling Lines)

1. **Extract graphics operators** from content stream
2. **Detect ruling lines**:
   - Find horizontal and vertical lines
   - Merge collinear lines
   - Filter by minimum length
3. **Build grid**:
   - Find line intersections
   - Create cell grid structure
   - Validate grid dimensions
4. **Create table region** with grid coordinates

### Stream Mode (Without Ruling Lines)

1. **Extract text elements** with X,Y positions
2. **Analyze projection profiles**:
   - Create horizontal profile (text density by Y)
   - Create vertical profile (text density by X)
3. **Find gaps** in profiles:
   - Low density = whitespace
   - Whitespace = potential boundaries
4. **Detect alignments**:
   - Text edges that align vertically (columns)
   - Text edges that align horizontally (rows)
5. **Create table region** with row/column coordinates

### Auto-Detection

The system automatically chooses the best mode:
- If ruling lines detected (≥2 horizontal + ≥2 vertical) → **Lattice Mode**
- Otherwise → **Stream Mode**

## Implementation Files

### Core Detection
- `internal/application/detector/table_detector.go` - Main table detector
- `internal/application/detector/ruling_line.go` - Lattice mode detection
- `internal/application/detector/grid_builder.go` - Grid structure builder
- `internal/application/detector/whitespace_analyzer.go` - Stream mode detection
- `internal/application/detector/projection_profile.go` - Projection analysis

### Graphics Parsing
- `internal/application/extractor/graphics_parser.go` - Graphics operator parser

### Tests
- `internal/application/detector/detector_test.go` - Comprehensive tests
- `internal/application/extractor/graphics_parser_test.go` - Graphics tests

## Test Coverage

- Detector package: 49.0% coverage
- All tests passing (26 test cases)
- 0 linter issues (except expected duplication)

## Next Steps (Phase 2.7)

Once table regions are detected, Phase 2.7 will:
1. Extract text content from each cell
2. Handle merged cells and multi-line content
3. Export to CSV, JSON, and Excel formats

## References

- **tabula-java**: Reference implementation
  - `NurminenDetectionAlgorithm.java` - Stream mode inspiration
  - `SpreadsheetExtractionAlgorithm.java` - Lattice mode inspiration
  - `BasicExtractionAlgorithm.java` - Whitespace analysis

- **Academic**: Nurminen's master's thesis on table detection
  - http://dspace.cc.tut.fi/dpub/bitstream/handle/123456789/21520/Nurminen.pdf

- **PDF Specification**: ISO 32000-1:2008
  - Section 8: Graphics (operators, state, paths)
  - Section 9: Text (positioning, rendering)

## License

Part of the GoPDF project - MIT License
