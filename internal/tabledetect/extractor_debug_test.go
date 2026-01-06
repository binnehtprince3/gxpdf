package tabledetect

import (
	"fmt"
	"testing"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCellExtractionDebug is a diagnostic test to reproduce the empty cells bug
// reported in ISSUE_TABLE_EXTRACTION_EMPTY_CELLS.md
func TestCellExtractionDebug(t *testing.T) {
	// Create a simple 2x2 table with text elements
	// Table structure:
	//   X: 0-100, 100-200
	//   Y: 200-220, 180-200
	//
	//   +------------+------------+
	//   | "Header1"  | "Header2"  | Y: 200-220
	//   +------------+------------+
	//   | "Cell1"    | "Cell2"    | Y: 180-200
	//   +------------+------------+
	//     X: 0-100     X: 100-200

	textElements := []*extractor.TextElement{
		// Row 0, Col 0: "Header1" at (10, 205)
		extractor.NewTextElement("Header1", 10, 205, 60, 10, "/F1", 12),
		// Row 0, Col 1: "Header2" at (110, 205)
		extractor.NewTextElement("Header2", 110, 205, 60, 10, "/F1", 12),
		// Row 1, Col 0: "Cell1" at (10, 185)
		extractor.NewTextElement("Cell1", 10, 185, 40, 10, "/F1", 12),
		// Row 1, Col 1: "Cell2" at (110, 185)
		extractor.NewTextElement("Cell2", 110, 185, 40, 10, "/F1", 12),
	}

	fmt.Println("\n=== DEBUG: Text Elements ===")
	for i, elem := range textElements {
		fmt.Printf("[%d] %s\n", i, elem.String())
		fmt.Printf("    Center: (%.2f, %.2f)\n", elem.CenterX(), elem.CenterY())
	}

	// Create a TableRegion manually (Stream mode)
	// IMPORTANT: WhitespaceAnalyzer.DetectRows() returns rows sorted ASCENDING (bottom to top)
	// So Rows array = [bottom_edge, ..., top_edge] = [low_Y, ..., high_Y]
	region := &TableRegion{
		Bounds: extractor.NewRectangle(0, 180, 200, 40), // Full table bounds
		Method: MethodStream,
		Rows: []float64{
			180, // Bottom of table (lowest Y)
			200, // Middle row boundary
			220, // Top of table (highest Y)
		},
		Columns: []float64{
			0,   // Left of col 0
			100, // Right of col 0 / Left of col 1
			200, // Right of col 1
		},
	}

	fmt.Println("\n=== DEBUG: TableRegion ===")
	fmt.Printf("Bounds: %s\n", region.Bounds.String())
	fmt.Printf("Method: %s\n", region.Method)
	fmt.Printf("Rows: %v\n", region.Rows)
	fmt.Printf("Columns: %v\n", region.Columns)

	// Create TableExtractor
	tableExtractor := NewTableExtractor(textElements)

	// Extract table
	tbl, err := tableExtractor.ExtractTable(region)
	require.NoError(t, err, "failed to extract table")
	require.NotNil(t, tbl)

	// Verify table structure
	assert.Equal(t, 2, tbl.RowCount, "table should have 2 rows")
	assert.Equal(t, 2, tbl.ColCount, "table should have 2 columns")

	fmt.Println("\n=== DEBUG: Extracted Table ===")
	fmt.Printf("Rows: %d, Columns: %d\n", tbl.RowCount, tbl.ColCount)

	// Check cell contents
	for r := 0; r < tbl.RowCount; r++ {
		for c := 0; c < tbl.ColCount; c++ {
			cell := tbl.GetCell(r, c)
			fmt.Printf("Cell[%d,%d]: %q (bounds: %s)\n", r, c, cell.Text, cell.Bounds.String())
		}
	}

	// CRITICAL ASSERTIONS - These should pass but currently fail
	cell00 := tbl.GetCell(0, 0)
	cell01 := tbl.GetCell(0, 1)
	cell10 := tbl.GetCell(1, 0)
	cell11 := tbl.GetCell(1, 1)

	// Check that cells are not empty
	assert.NotEmpty(t, cell00.Text, "Cell[0,0] should contain 'Header1'")
	assert.NotEmpty(t, cell01.Text, "Cell[0,1] should contain 'Header2'")
	assert.NotEmpty(t, cell10.Text, "Cell[1,0] should contain 'Cell1'")
	assert.NotEmpty(t, cell11.Text, "Cell[1,1] should contain 'Cell2'")

	// Check exact content
	assert.Equal(t, "Header1", cell00.Text, "Cell[0,0] content mismatch")
	assert.Equal(t, "Header2", cell01.Text, "Cell[0,1] content mismatch")
	assert.Equal(t, "Cell1", cell10.Text, "Cell[1,0] content mismatch")
	assert.Equal(t, "Cell2", cell11.Text, "Cell[1,1] content mismatch")
}

// TestCellExtractorDirectly tests CellExtractor in isolation
func TestCellExtractorDirectly(t *testing.T) {
	// Create text elements at known positions
	textElements := []*extractor.TextElement{
		// Text at (10, 100) with size (50, 10)
		extractor.NewTextElement("TestText", 10, 100, 50, 10, "/F1", 12),
	}

	fmt.Println("\n=== DEBUG: Direct CellExtractor Test ===")
	fmt.Printf("Text element: %s\n", textElements[0].String())
	fmt.Printf("Center: (%.2f, %.2f)\n", textElements[0].CenterX(), textElements[0].CenterY())

	// Create CellExtractor
	cellExtractor := extractor.NewCellExtractor(textElements)

	// Define cell bounds that should contain the text element
	// Text is at (10, 100) with width 50, height 10
	// So it spans: X: 10-60, Y: 100-110
	// Center is at: X: 35, Y: 105
	cellBounds := extractor.NewRectangle(0, 95, 100, 20) // X: 0-100, Y: 95-115

	fmt.Printf("Cell bounds: %s\n", cellBounds.String())
	fmt.Printf("Bounds contains center? %v\n", cellBounds.Contains(textElements[0].CenterX(), textElements[0].CenterY()))

	// Extract content
	content := cellExtractor.ExtractCellContent(cellBounds)

	fmt.Printf("Extracted content: %q\n", content)

	// Verify
	assert.Equal(t, "TestText", content, "CellExtractor should extract text from bounds")
}
