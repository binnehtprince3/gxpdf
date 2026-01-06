package tabledetect

import (
	"fmt"
	"testing"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test interface compliance

func TestInterfaceCompliance(t *testing.T) {
	// Verify that default implementations satisfy interfaces
	var _ RulingLineDetector = (*DefaultRulingLineDetector)(nil)
	var _ WhitespaceAnalyzer = (*DefaultWhitespaceAnalyzer)(nil)
	var _ ProjectionAnalyzer = (*DefaultProjectionAnalyzer)(nil)
	var _ GridBuilder = (*DefaultGridBuilder)(nil)
	var _ TableDetector = (*DefaultTableDetector)(nil)

	// Verify that mock implementations satisfy interfaces
	var _ RulingLineDetector = (*MockRulingLineDetector)(nil)
	var _ WhitespaceAnalyzer = (*MockWhitespaceAnalyzer)(nil)
	var _ ProjectionAnalyzer = (*MockProjectionAnalyzer)(nil)
	var _ GridBuilder = (*MockGridBuilder)(nil)
	var _ TableDetector = (*MockTableDetector)(nil)
}

// Test TableDetector with mocks

func TestTableDetector_WithMockRulingDetector_Lattice(t *testing.T) {
	// Create mock ruling detector that returns lattice lines
	mockRuling := NewMockRulingLineDetector()
	mockRuling.DetectRulingLinesResult = []*RulingLine{
		// Horizontal lines
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(200, 0)),
		NewRulingLine(extractor.NewPoint(0, 100), extractor.NewPoint(200, 100)),
		// Vertical lines
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(0, 100)),
		NewRulingLine(extractor.NewPoint(200, 0), extractor.NewPoint(200, 100)),
	}

	// Create detector with mock
	detector := NewTableDetectorWithDeps(
		mockRuling,
		NewDefaultWhitespaceAnalyzer(),
		NewDefaultGridBuilder(),
	)

	// Test detection
	graphics := []*extractor.GraphicsElement{
		{Type: extractor.GraphicsTypeLine},
	}
	textElements := []*extractor.TextElement{}

	mode := detector.DetectMode(textElements, graphics)
	assert.Equal(t, MethodLattice, mode, "Should detect lattice mode with ruling lines")
}

func TestTableDetector_WithMockRulingDetector_Stream(t *testing.T) {
	// Create mock ruling detector that returns no lines
	mockRuling := NewMockRulingLineDetector()
	mockRuling.DetectRulingLinesResult = []*RulingLine{}

	// Create detector with mock
	detector := NewTableDetectorWithDeps(
		mockRuling,
		NewDefaultWhitespaceAnalyzer(),
		NewDefaultGridBuilder(),
	)

	// Test detection
	graphics := []*extractor.GraphicsElement{}
	textElements := []*extractor.TextElement{}

	mode := detector.DetectMode(textElements, graphics)
	assert.Equal(t, MethodStream, mode, "Should detect stream mode without ruling lines")
}

func TestTableDetector_WithMockWhitespaceAnalyzer(t *testing.T) {
	// Create mock whitespace analyzer with predefined columns/rows
	mockWhitespace := NewMockWhitespaceAnalyzer()
	mockWhitespace.DetectColumnsResult = []float64{0, 100, 200}
	mockWhitespace.DetectRowsResult = []float64{0, 50, 100}

	// Create detector with mock
	detector := NewTableDetectorWithDeps(
		NewDefaultRulingLineDetector(),
		mockWhitespace,
		NewDefaultGridBuilder(),
	)

	// Test stream detection
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Text1", 10, 10, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text2", 110, 10, 50, 10, "/F1", 10),
	}

	tables, err := detector.DetectTablesStream(textElements)
	require.NoError(t, err)
	require.Len(t, tables, 1)

	table := tables[0]
	assert.Equal(t, MethodStream, table.Method)
	assert.Equal(t, mockWhitespace.DetectColumnsResult, table.Columns)
	assert.Equal(t, mockWhitespace.DetectRowsResult, table.Rows)
}

func TestTableDetector_WithMockGridBuilder(t *testing.T) {
	// Create mock grid builder
	mockGrid := NewMockGridBuilder()
	mockGrid.BuildGridResult = NewGrid(
		[]float64{0, 50, 100},
		[]float64{0, 100, 200},
	)

	// Create mock ruling detector
	mockRuling := NewMockRulingLineDetector()
	mockRuling.DetectRulingLinesResult = []*RulingLine{
		// Horizontal
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(200, 0)),
		NewRulingLine(extractor.NewPoint(0, 50), extractor.NewPoint(200, 50)),
		NewRulingLine(extractor.NewPoint(0, 100), extractor.NewPoint(200, 100)),
		// Vertical
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(0, 100)),
		NewRulingLine(extractor.NewPoint(100, 0), extractor.NewPoint(100, 100)),
		NewRulingLine(extractor.NewPoint(200, 0), extractor.NewPoint(200, 100)),
	}

	// Create detector with mocks
	detector := NewTableDetectorWithDeps(
		mockRuling,
		NewDefaultWhitespaceAnalyzer(),
		mockGrid,
	)

	// Test lattice detection
	graphics := []*extractor.GraphicsElement{
		{Type: extractor.GraphicsTypeLine},
	}
	textElements := []*extractor.TextElement{}

	tables, err := detector.DetectTablesLattice(textElements, graphics)
	require.NoError(t, err)
	require.Len(t, tables, 1)

	table := tables[0]
	assert.Equal(t, MethodLattice, table.Method)
	assert.NotNil(t, table.Grid)
	assert.Equal(t, 2, table.Grid.RowCount())
	assert.Equal(t, 2, table.Grid.ColumnCount())
}

func TestTableDetector_WithMockGridBuilder_Error(t *testing.T) {
	// Create mock grid builder that returns error
	mockGrid := NewMockGridBuilder()
	mockGrid.BuildGridError = fmt.Errorf("grid build failed")

	// Create mock ruling detector
	mockRuling := NewMockRulingLineDetector()
	mockRuling.DetectRulingLinesResult = []*RulingLine{
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(200, 0)),
		NewRulingLine(extractor.NewPoint(0, 100), extractor.NewPoint(200, 100)),
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(0, 100)),
		NewRulingLine(extractor.NewPoint(200, 0), extractor.NewPoint(200, 100)),
	}

	// Create mock whitespace analyzer for fallback
	mockWhitespace := NewMockWhitespaceAnalyzer()
	mockWhitespace.DetectColumnsResult = []float64{0, 100}
	mockWhitespace.DetectRowsResult = []float64{0, 50}

	// Create detector with mocks
	detector := NewTableDetectorWithDeps(
		mockRuling,
		mockWhitespace,
		mockGrid,
	)

	// Test lattice detection - should fall back to stream mode
	graphics := []*extractor.GraphicsElement{
		{Type: extractor.GraphicsTypeLine},
	}
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Text", 10, 10, 50, 10, "/F1", 10),
	}

	tables, err := detector.DetectTablesLattice(textElements, graphics)
	require.NoError(t, err, "Should fall back to stream mode on grid build error")
	require.Len(t, tables, 1)

	table := tables[0]
	assert.Equal(t, MethodStream, table.Method, "Should use stream mode as fallback")
}

func TestTableDetector_WithCustomFunctions(t *testing.T) {
	// Test using custom functions in mocks
	rulingCallCount := 0
	whitespaceCallCount := 0

	mockRuling := NewMockRulingLineDetector()
	mockRuling.DetectRulingLinesFunc = func(graphics []*extractor.GraphicsElement) ([]*RulingLine, error) {
		rulingCallCount++
		return []*RulingLine{}, nil
	}

	mockWhitespace := NewMockWhitespaceAnalyzer()
	mockWhitespace.DetectColumnsFunc = func(elements []*extractor.TextElement) []float64 {
		whitespaceCallCount++
		return []float64{0, 100}
	}
	mockWhitespace.DetectRowsFunc = func(elements []*extractor.TextElement) []float64 {
		return []float64{0, 50}
	}

	detector := NewTableDetectorWithDeps(
		mockRuling,
		mockWhitespace,
		NewDefaultGridBuilder(),
	)

	// Test detection
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Text", 10, 10, 50, 10, "/F1", 10),
	}
	graphics := []*extractor.GraphicsElement{}

	_, err := detector.DetectTables(textElements, graphics)
	require.NoError(t, err)

	// Verify mocks were called
	assert.Equal(t, 1, rulingCallCount, "RulingDetector should be called once")
	assert.Equal(t, 1, whitespaceCallCount, "WhitespaceAnalyzer should be called once")
}

// Test that we can easily swap implementations

func TestSwappableImplementations(t *testing.T) {
	// Create detector with default implementations
	detector1 := NewDefaultTableDetector()
	assert.NotNil(t, detector1)

	// Create detector with custom mock implementations
	detector2 := NewTableDetectorWithDeps(
		NewMockRulingLineDetector(),
		NewMockWhitespaceAnalyzer(),
		NewMockGridBuilder(),
	)
	assert.NotNil(t, detector2)

	// Create detector with mixed implementations
	detector3 := NewTableDetectorWithDeps(
		NewDefaultRulingLineDetector(),
		NewMockWhitespaceAnalyzer(),
		NewDefaultGridBuilder(),
	)
	assert.NotNil(t, detector3)

	// All should be valid TableDetector instances
	var _ TableDetector = detector1
	var _ TableDetector = detector2
	var _ TableDetector = detector3
}

// Test backward compatibility

func TestBackwardCompatibility_Constructors(t *testing.T) {
	// Old constructor should still work
	detector1 := NewTableDetector()
	assert.NotNil(t, detector1)

	detector2 := NewRulingLineDetector()
	assert.NotNil(t, detector2)

	detector3 := NewWhitespaceAnalyzer()
	assert.NotNil(t, detector3)

	detector4 := NewGridBuilder()
	assert.NotNil(t, detector4)

	detector5 := NewProjectionAnalyzer()
	assert.NotNil(t, detector5)
}

func TestBackwardCompatibility_BuilderPattern(t *testing.T) {
	// Old builder pattern should still work
	detector := NewTableDetector().
		WithRulingDetector(NewDefaultRulingLineDetector()).
		WithWhitespaceAnalyzer(NewDefaultWhitespaceAnalyzer()).
		WithGridBuilder(NewDefaultGridBuilder())

	assert.NotNil(t, detector)

	// Test detection still works
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Text", 10, 10, 50, 10, "/F1", 10),
	}
	graphics := []*extractor.GraphicsElement{}

	tables, err := detector.DetectTables(textElements, graphics)
	require.NoError(t, err)
	assert.NotNil(t, tables)
}

// Demonstrate benefits of interface-based design

func TestBenefits_EasyTesting(t *testing.T) {
	t.Run("Can test in complete isolation", func(t *testing.T) {
		// Create a detector with all mocked dependencies
		mockRuling := NewMockRulingLineDetector()
		mockRuling.DetectRulingLinesResult = []*RulingLine{}

		mockWhitespace := NewMockWhitespaceAnalyzer()
		mockWhitespace.DetectColumnsResult = []float64{0, 100}
		mockWhitespace.DetectRowsResult = []float64{0, 50}

		mockGrid := NewMockGridBuilder()

		detector := NewTableDetectorWithDeps(mockRuling, mockWhitespace, mockGrid)

		// Test without any real PDF processing
		textElements := []*extractor.TextElement{
			extractor.NewTextElement("Test", 10, 10, 50, 10, "/F1", 10),
		}

		tables, err := detector.DetectTablesStream(textElements)
		require.NoError(t, err)
		assert.Len(t, tables, 1)
	})

	t.Run("Can test error scenarios easily", func(t *testing.T) {
		// Mock that simulates error
		mockRuling := NewMockRulingLineDetector()
		mockRuling.DetectRulingLinesError = fmt.Errorf("detection failed")

		mockWhitespace := NewMockWhitespaceAnalyzer()
		mockWhitespace.DetectColumnsResult = []float64{0, 100}
		mockWhitespace.DetectRowsResult = []float64{0, 50}

		detector := NewTableDetectorWithDeps(
			mockRuling,
			mockWhitespace,
			NewDefaultGridBuilder(),
		)

		// Test error handling - should propagate error
		_, err := detector.DetectTablesLattice(nil, []*extractor.GraphicsElement{})
		assert.Error(t, err, "Should propagate error from ruling detector")
		assert.Contains(t, err.Error(), "detection failed")
	})
}

func TestBenefits_ExtensibilityExample(t *testing.T) {
	// Example: Custom ruling detector with different algorithm
	customRuling := NewMockRulingLineDetector()
	customRuling.DetectRulingLinesFunc = func(graphics []*extractor.GraphicsElement) ([]*RulingLine, error) {
		// Imagine this is a machine learning-based detector
		// or uses a different algorithm (Hough transform, etc.)
		return []*RulingLine{
			// Custom detection logic here
		}, nil
	}

	// Easy to plug in custom implementation
	detector := NewTableDetectorWithDeps(
		customRuling,
		NewDefaultWhitespaceAnalyzer(),
		NewDefaultGridBuilder(),
	)

	assert.NotNil(t, detector)
}
