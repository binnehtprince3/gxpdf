package tabledetect

import (
	"testing"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test RulingLine

func TestNewRulingLine(t *testing.T) {
	start := extractor.NewPoint(0, 0)
	end := extractor.NewPoint(100, 0)

	line := NewRulingLine(start, end)

	require.NotNil(t, line)
	assert.Equal(t, start, line.Start)
	assert.Equal(t, end, line.End)
	assert.True(t, line.IsHorizontal)
}

func TestRulingLine_Length(t *testing.T) {
	tests := []struct {
		name     string
		start    extractor.Point
		end      extractor.Point
		expected float64
	}{
		{"horizontal line", extractor.NewPoint(0, 0), extractor.NewPoint(100, 0), 100.0},
		{"vertical line", extractor.NewPoint(0, 0), extractor.NewPoint(0, 50), 50.0},
		{"diagonal line", extractor.NewPoint(0, 0), extractor.NewPoint(3, 4), 5.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			line := NewRulingLine(tt.start, tt.end)
			assert.InDelta(t, tt.expected, line.Length(), 0.001)
		})
	}
}

func TestRulingLine_Intersects(t *testing.T) {
	// Horizontal line
	hLine := NewRulingLine(extractor.NewPoint(0, 50), extractor.NewPoint(100, 50))

	// Vertical line
	vLine := NewRulingLine(extractor.NewPoint(50, 0), extractor.NewPoint(50, 100))

	// Should intersect at (50, 50)
	point := hLine.Intersects(vLine)
	require.NotNil(t, point)
	assert.Equal(t, 50.0, point.X)
	assert.Equal(t, 50.0, point.Y)
}

func TestRulingLine_Intersects_NoIntersection(t *testing.T) {
	// Horizontal line
	hLine := NewRulingLine(extractor.NewPoint(0, 50), extractor.NewPoint(100, 50))

	// Vertical line that doesn't intersect
	vLine := NewRulingLine(extractor.NewPoint(150, 0), extractor.NewPoint(150, 100))

	// Should not intersect
	point := hLine.Intersects(vLine)
	assert.Nil(t, point)
}

func TestRulingLine_Intersects_Parallel(t *testing.T) {
	// Two horizontal lines
	line1 := NewRulingLine(extractor.NewPoint(0, 50), extractor.NewPoint(100, 50))
	line2 := NewRulingLine(extractor.NewPoint(0, 60), extractor.NewPoint(100, 60))

	// Should not intersect (parallel)
	point := line1.Intersects(line2)
	assert.Nil(t, point)
}

// Test RulingLineDetector

func TestNewRulingLineDetector(t *testing.T) {
	detector := NewRulingLineDetector()

	require.NotNil(t, detector)
	// minLineLength and tolerance are private fields, tested via behavior
}

func TestRulingLineDetector_WithMinLineLength(t *testing.T) {
	detector := NewRulingLineDetector().WithMinLineLength(20.0)
	require.NotNil(t, detector)
	// minLineLength is a private field, tested via behavior
}

func TestRulingLineDetector_DetectRulingLines(t *testing.T) {
	detector := NewRulingLineDetector()

	// Create graphics elements (lines)
	graphics := []*extractor.GraphicsElement{
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 0),
				extractor.NewPoint(100, 0),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 50),
				extractor.NewPoint(100, 50),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 0),
				extractor.NewPoint(0, 50),
			},
		},
	}

	lines, err := detector.DetectRulingLines(graphics)

	require.NoError(t, err)
	assert.Len(t, lines, 3)
}

func TestRulingLineDetector_FindIntersections(t *testing.T) {
	detector := NewRulingLineDetector()

	// Create two lines that intersect
	lines := []*RulingLine{
		NewRulingLine(extractor.NewPoint(0, 50), extractor.NewPoint(100, 50)),
		NewRulingLine(extractor.NewPoint(50, 0), extractor.NewPoint(50, 100)),
	}

	intersections := detector.FindIntersections(lines)

	require.Len(t, intersections, 1)
	assert.Equal(t, 50.0, intersections[0].X)
	assert.Equal(t, 50.0, intersections[0].Y)
}

// Test Grid and GridBuilder

func TestNewCell(t *testing.T) {
	bounds := extractor.NewRectangle(0, 0, 100, 50)
	cell := NewCell(0, 0, bounds)

	require.NotNil(t, cell)
	assert.Equal(t, 0, cell.Row)
	assert.Equal(t, 0, cell.Column)
	assert.Equal(t, bounds, cell.Bounds)
}

func TestNewGrid(t *testing.T) {
	rows := []float64{0, 50, 100}
	columns := []float64{0, 100, 200}

	grid := NewGrid(rows, columns)

	require.NotNil(t, grid)
	assert.Len(t, grid.Rows, 3)
	assert.Len(t, grid.Columns, 3)
	assert.Equal(t, 2, grid.RowCount())
	assert.Equal(t, 2, grid.ColumnCount())
}

func TestGrid_Bounds(t *testing.T) {
	rows := []float64{0, 50, 100}
	columns := []float64{0, 100, 200}

	grid := NewGrid(rows, columns)
	bounds := grid.Bounds()

	assert.Equal(t, 0.0, bounds.X)
	assert.Equal(t, 0.0, bounds.Y)
	assert.Equal(t, 200.0, bounds.Width)
	assert.Equal(t, 100.0, bounds.Height)
}

func TestGridBuilder_BuildGrid(t *testing.T) {
	builder := NewGridBuilder()

	// Create ruling lines forming a 2x2 grid
	lines := []*RulingLine{
		// Horizontal lines
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(200, 0)),
		NewRulingLine(extractor.NewPoint(0, 50), extractor.NewPoint(200, 50)),
		NewRulingLine(extractor.NewPoint(0, 100), extractor.NewPoint(200, 100)),
		// Vertical lines
		NewRulingLine(extractor.NewPoint(0, 0), extractor.NewPoint(0, 100)),
		NewRulingLine(extractor.NewPoint(100, 0), extractor.NewPoint(100, 100)),
		NewRulingLine(extractor.NewPoint(200, 0), extractor.NewPoint(200, 100)),
	}

	grid, err := builder.BuildGrid(lines)

	require.NoError(t, err)
	require.NotNil(t, grid)
	assert.Equal(t, 2, grid.RowCount())
	assert.Equal(t, 2, grid.ColumnCount())
}

// Test ProjectionProfile

func TestNewProjectionProfile(t *testing.T) {
	bins := []float64{0, 1, 2, 3}
	profile := NewProjectionProfile(bins, 10.0, 0, 40)

	require.NotNil(t, profile)
	assert.Equal(t, 4, profile.BinCount())
	assert.Equal(t, 10.0, profile.BinSize)
}

func TestProjectionProfile_GetDensity(t *testing.T) {
	bins := []float64{0, 1, 2, 3}
	profile := NewProjectionProfile(bins, 10.0, 0, 40)

	// Test getting density at different coordinates
	assert.Equal(t, 0.0, profile.GetDensity(5))
	assert.Equal(t, 1.0, profile.GetDensity(15))
	assert.Equal(t, 2.0, profile.GetDensity(25))

	// Test out of range
	assert.Equal(t, 0.0, profile.GetDensity(-10))
	assert.Equal(t, 0.0, profile.GetDensity(100))
}

func TestProjectionAnalyzer_AnalyzeHorizontal(t *testing.T) {
	analyzer := NewProjectionAnalyzer()

	// Create text elements at different Y positions
	elements := []*extractor.TextElement{
		extractor.NewTextElement("Text1", 0, 0, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text2", 60, 0, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text3", 0, 50, 50, 10, "/F1", 10),
	}

	profile := analyzer.AnalyzeHorizontal(elements)

	require.NotNil(t, profile)
	assert.Greater(t, profile.BinCount(), 0)
}

func TestProjectionAnalyzer_FindGaps(t *testing.T) {
	analyzer := NewProjectionAnalyzer()

	// Create profile with obvious gaps
	bins := []float64{10, 10, 0, 0, 10, 10}
	profile := NewProjectionProfile(bins, 10.0, 0, 60)

	gaps := analyzer.FindGaps(profile)

	require.NotNil(t, gaps)
	// Should find at least one gap
	assert.Greater(t, len(gaps), 0)
}

// Test WhitespaceAnalyzer

func TestNewWhitespaceAnalyzer(t *testing.T) {
	analyzer := NewWhitespaceAnalyzer()

	require.NotNil(t, analyzer)
	assert.Equal(t, 10.0, analyzer.minGapWidth)
	assert.Equal(t, 2.0, analyzer.alignmentTolerance)
}

func TestWhitespaceAnalyzer_DetectColumns(t *testing.T) {
	analyzer := NewWhitespaceAnalyzer()

	// Create text elements in two columns
	elements := []*extractor.TextElement{
		// Column 1
		extractor.NewTextElement("Text1", 0, 0, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text2", 0, 20, 50, 10, "/F1", 10),
		// Column 2
		extractor.NewTextElement("Text3", 100, 0, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text4", 100, 20, 50, 10, "/F1", 10),
	}

	columns := analyzer.DetectColumns(elements)

	require.NotNil(t, columns)
	// Should detect at least 2 column boundaries (left and right edges)
	assert.GreaterOrEqual(t, len(columns), 2)
}

func TestWhitespaceAnalyzer_DetectRows(t *testing.T) {
	analyzer := NewWhitespaceAnalyzer()

	// Create text elements in two rows
	elements := []*extractor.TextElement{
		// Row 1
		extractor.NewTextElement("Text1", 0, 0, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text2", 60, 0, 50, 10, "/F1", 10),
		// Row 2
		extractor.NewTextElement("Text3", 0, 50, 50, 10, "/F1", 10),
		extractor.NewTextElement("Text4", 60, 50, 50, 10, "/F1", 10),
	}

	rows := analyzer.DetectRows(elements)

	require.NotNil(t, rows)
	// Should detect at least 2 row boundaries (top and bottom edges)
	assert.GreaterOrEqual(t, len(rows), 2)
}

// Test TableDetector

func TestNewTableDetector(t *testing.T) {
	detector := NewTableDetector()

	require.NotNil(t, detector)
	// Internal dependencies are private fields, tested via behavior
}

func TestTableDetector_DetectMode_Lattice(t *testing.T) {
	detector := NewTableDetector()

	// Create graphics with ruling lines
	graphics := []*extractor.GraphicsElement{
		// Horizontal lines
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 0),
				extractor.NewPoint(100, 0),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 50),
				extractor.NewPoint(100, 50),
			},
		},
		// Vertical lines
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 0),
				extractor.NewPoint(0, 50),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(100, 0),
				extractor.NewPoint(100, 50),
			},
		},
	}

	mode := detector.DetectMode([]*extractor.TextElement{}, graphics)
	assert.Equal(t, MethodLattice, mode)
}

func TestTableDetector_DetectMode_Stream(t *testing.T) {
	detector := NewTableDetector()

	// No graphics - should use stream mode
	mode := detector.DetectMode([]*extractor.TextElement{}, []*extractor.GraphicsElement{})
	assert.Equal(t, MethodStream, mode)
}

func TestExtractionMethod_String(t *testing.T) {
	tests := []struct {
		method   ExtractionMethod
		expected string
	}{
		{MethodLattice, "Lattice"},
		{MethodStream, "Stream"},
		{MethodAuto, "Auto"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.method.String())
		})
	}
}

func TestNewTableRegion(t *testing.T) {
	bounds := extractor.NewRectangle(0, 0, 200, 100)
	region := NewTableRegion(bounds, MethodLattice)

	require.NotNil(t, region)
	assert.Equal(t, bounds, region.Bounds)
	assert.Equal(t, MethodLattice, region.Method)
	assert.True(t, region.HasRulingLines)
}

func TestTableRegion_String(t *testing.T) {
	bounds := extractor.NewRectangle(0, 0, 200, 100)
	region := NewTableRegion(bounds, MethodLattice)

	str := region.String()
	assert.Contains(t, str, "Lattice")
	assert.Contains(t, str, "200")
	assert.Contains(t, str, "100")
}
