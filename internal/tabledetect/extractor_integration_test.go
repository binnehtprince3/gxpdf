package tabledetect

import (
	"testing"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestStreamTableExtraction_PositiveHeights verifies that all cell bounds have positive dimensions
func TestStreamTableExtraction_PositiveHeights(t *testing.T) {
	// Create test data with rows sorted ascending (as WhitespaceAnalyzer produces)
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Top Left", 10, 205, 60, 10, "/F1", 12),
		extractor.NewTextElement("Top Right", 110, 205, 60, 10, "/F1", 12),
		extractor.NewTextElement("Bottom Left", 10, 185, 60, 10, "/F1", 12),
		extractor.NewTextElement("Bottom Right", 110, 185, 60, 10, "/F1", 12),
	}

	region := &TableRegion{
		Bounds:  extractor.NewRectangle(0, 180, 200, 40),
		Method:  MethodStream,
		Rows:    []float64{180, 200, 220}, // ASCENDING: bottom to top
		Columns: []float64{0, 100, 200},   // ASCENDING: left to right
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	// Verify all cells have positive dimensions
	for r := 0; r < tbl.RowCount; r++ {
		for c := 0; c < tbl.ColCount; c++ {
			cell := tbl.GetCell(r, c)
			assert.Greater(t, cell.Bounds.Width, 0.0,
				"Cell[%d,%d] should have positive width", r, c)
			assert.Greater(t, cell.Bounds.Height, 0.0,
				"Cell[%d,%d] should have positive height", r, c)
		}
	}
}

// TestStreamTableExtraction_TextElementMatching verifies text elements are correctly assigned to cells
func TestStreamTableExtraction_TextElementMatching(t *testing.T) {
	// Create a 3x3 table with specific text in each cell
	textElements := []*extractor.TextElement{
		// Row 0 (Y: 200-220)
		extractor.NewTextElement("A1", 10, 205, 20, 10, "/F1", 12),  // Col 0 (X: 0-50)
		extractor.NewTextElement("B1", 60, 205, 20, 10, "/F1", 12),  // Col 1 (X: 50-100)
		extractor.NewTextElement("C1", 110, 205, 20, 10, "/F1", 12), // Col 2 (X: 100-150)
		// Row 1 (Y: 180-200)
		extractor.NewTextElement("A2", 10, 185, 20, 10, "/F1", 12),  // Col 0
		extractor.NewTextElement("B2", 60, 185, 20, 10, "/F1", 12),  // Col 1
		extractor.NewTextElement("C2", 110, 185, 20, 10, "/F1", 12), // Col 2
		// Row 2 (Y: 160-180)
		extractor.NewTextElement("A3", 10, 165, 20, 10, "/F1", 12),  // Col 0
		extractor.NewTextElement("B3", 60, 165, 20, 10, "/F1", 12),  // Col 1
		extractor.NewTextElement("C3", 110, 165, 20, 10, "/F1", 12), // Col 2
	}

	region := &TableRegion{
		Bounds:  extractor.NewRectangle(0, 160, 150, 60),
		Method:  MethodStream,
		Rows:    []float64{160, 180, 200, 220}, // 3 rows
		Columns: []float64{0, 50, 100, 150},    // 3 columns
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	// Verify table structure
	assert.Equal(t, 3, tbl.RowCount, "should have 3 rows")
	assert.Equal(t, 3, tbl.ColCount, "should have 3 columns")

	// Verify content of each cell
	expectedContent := [][]string{
		{"A1", "B1", "C1"}, // Row 0 (top)
		{"A2", "B2", "C2"}, // Row 1 (middle)
		{"A3", "B3", "C3"}, // Row 2 (bottom)
	}

	for r := 0; r < 3; r++ {
		for c := 0; c < 3; c++ {
			cell := tbl.GetCell(r, c)
			assert.Equal(t, expectedContent[r][c], cell.Text,
				"Cell[%d,%d] content mismatch", r, c)
		}
	}
}

// TestStreamTableExtraction_CoordinateSystem verifies proper handling of PDF coordinate system
func TestStreamTableExtraction_CoordinateSystem(t *testing.T) {
	// PDF coordinates: Y increases upward
	// Row 0 should be at the TOP (highest Y)
	// Row n should be at the BOTTOM (lowest Y)

	textElements := []*extractor.TextElement{
		// Top row (highest Y)
		extractor.NewTextElement("TopText", 10, 300, 60, 10, "/F1", 12),
		// Bottom row (lowest Y)
		extractor.NewTextElement("BottomText", 10, 100, 60, 10, "/F1", 12),
	}

	region := &TableRegion{
		Bounds:  extractor.NewRectangle(0, 90, 100, 220),
		Method:  MethodStream,
		Rows:    []float64{90, 200, 310}, // Ascending: bottom to top
		Columns: []float64{0, 100},
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	// Row 0 should contain "TopText" (from highest Y coordinate)
	cell00 := tbl.GetCell(0, 0)
	assert.Equal(t, "TopText", cell00.Text, "Row 0 should be the top row")

	// Row 1 should contain "BottomText" (from lowest Y coordinate)
	cell10 := tbl.GetCell(1, 0)
	assert.Equal(t, "BottomText", cell10.Text, "Row 1 should be the bottom row")

	// Verify Y coordinates
	assert.Greater(t, cell00.Bounds.Y, cell10.Bounds.Y,
		"Top row should have higher Y coordinate than bottom row")
}

// TestStreamTableExtraction_EmptyCells verifies handling of cells without text
func TestStreamTableExtraction_EmptyCells(t *testing.T) {
	// Create a 2x2 table with text only in diagonal cells
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("TopLeft", 10, 205, 40, 10, "/F1", 12),
		extractor.NewTextElement("BottomRight", 110, 185, 40, 10, "/F1", 12),
		// No text in TopRight and BottomLeft cells
	}

	region := &TableRegion{
		Bounds:  extractor.NewRectangle(0, 180, 200, 40),
		Method:  MethodStream,
		Rows:    []float64{180, 200, 220},
		Columns: []float64{0, 100, 200},
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	// Verify filled cells
	assert.Equal(t, "TopLeft", tbl.GetCell(0, 0).Text)
	assert.Equal(t, "BottomRight", tbl.GetCell(1, 1).Text)

	// Verify empty cells
	assert.Empty(t, tbl.GetCell(0, 1).Text, "TopRight should be empty")
	assert.Empty(t, tbl.GetCell(1, 0).Text, "BottomLeft should be empty")
}

// TestStreamTableExtraction_MultilineCell verifies handling of multi-line cell content
func TestStreamTableExtraction_MultilineCell(t *testing.T) {
	// Create a cell with two lines of text
	textElements := []*extractor.TextElement{
		// Both text elements in the same cell (X: 0-100, Y: 180-220)
		extractor.NewTextElement("Line1", 10, 205, 40, 10, "/F1", 12), // Upper line
		extractor.NewTextElement("Line2", 10, 190, 40, 10, "/F1", 12), // Lower line
	}

	region := &TableRegion{
		Bounds:  extractor.NewRectangle(0, 180, 100, 40),
		Method:  MethodStream,
		Rows:    []float64{180, 220},
		Columns: []float64{0, 100},
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	cell := tbl.GetCell(0, 0)
	// Cell should contain both lines (joined with newline)
	assert.Contains(t, cell.Text, "Line1")
	assert.Contains(t, cell.Text, "Line2")
}

// TestStreamTableExtraction_LatticeMode verifies that lattice mode still works
func TestStreamTableExtraction_LatticeMode(t *testing.T) {
	// This test ensures we didn't break lattice mode extraction
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Cell1", 10, 105, 40, 10, "/F1", 12),
	}

	// Create a simple grid
	cell := NewCell(0, 0, extractor.NewRectangle(0, 100, 100, 20))
	grid := &Grid{
		Rows:    []float64{100, 120},
		Columns: []float64{0, 100},
		Cells:   [][]*Cell{{cell}},
	}

	region := &TableRegion{
		Bounds: extractor.NewRectangle(0, 100, 100, 20),
		Method: MethodLattice,
		Grid:   grid,
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	// Verify lattice mode extraction still works
	assert.Equal(t, 1, tbl.RowCount)
	assert.Equal(t, 1, tbl.ColCount)
	assert.Equal(t, "Cell1", tbl.GetCell(0, 0).Text)
	assert.Equal(t, "Lattice", tbl.Method)
}

// TestStreamTableExtraction_CyrillicContent verifies UTF-8/Cyrillic text handling
func TestStreamTableExtraction_CyrillicContent(t *testing.T) {
	// Create table with Cyrillic text
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Привет", 10, 205, 60, 10, "/F1", 12),  // Russian "Hello"
		extractor.NewTextElement("Мир", 110, 205, 40, 10, "/F1", 12),    // Russian "World"
		extractor.NewTextElement("Выписка", 10, 185, 80, 10, "/F1", 12), // Russian "Statement"
		extractor.NewTextElement("РОССИЯ", 110, 185, 60, 10, "/F1", 12), // Russian "RUSSIA"
	}

	region := &TableRegion{
		Bounds:  extractor.NewRectangle(0, 180, 200, 40),
		Method:  MethodStream,
		Rows:    []float64{180, 200, 220},
		Columns: []float64{0, 100, 200},
	}

	te := NewTableExtractor(textElements)
	tbl, err := te.ExtractTable(region)
	require.NoError(t, err)

	// Verify Cyrillic content is preserved
	assert.Equal(t, "Привет", tbl.GetCell(0, 0).Text)
	assert.Equal(t, "Мир", tbl.GetCell(0, 1).Text)
	assert.Equal(t, "Выписка", tbl.GetCell(1, 0).Text)
	assert.Equal(t, "РОССИЯ", tbl.GetCell(1, 1).Text)
}
