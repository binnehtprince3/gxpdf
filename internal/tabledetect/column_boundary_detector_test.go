package tabledetect

import (
	"testing"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/stretchr/testify/assert"
)

func TestColumnBoundaryDetector_DetectBoundaries_SimpleTable(t *testing.T) {
	// Simulate simple 3-column table
	// Column 1: X=50-100, Column 2: X=150-200, Column 3: X=250-300
	elements := []*extractor.TextElement{
		// Row 1
		newTextElement("A1", 50, 100, 50, 10),
		newTextElement("B1", 150, 100, 50, 10),
		newTextElement("C1", 250, 100, 50, 10),
		// Row 2
		newTextElement("A2", 50, 90, 50, 10),
		newTextElement("B2", 150, 90, 50, 10),
		newTextElement("C2", 250, 90, 50, 10),
		// Row 3
		newTextElement("A3", 50, 80, 50, 10),
		newTextElement("B3", 150, 80, 50, 10),
		newTextElement("C3", 250, 80, 50, 10),
	}

	detector := NewColumnBoundaryDetector()
	boundaries := detector.DetectBoundaries(elements)

	// Should detect boundaries around X=50, 100, 150, 200, 250, 300
	// After clustering and filtering, expect ~6 boundaries (3 columns * 2 edges)
	assert.GreaterOrEqual(t, len(boundaries), 3, "Should detect at least 3 boundaries")
	assert.LessOrEqual(t, len(boundaries), 6, "Should not have more than 6 boundaries")

	// Boundaries should be sorted
	for i := 1; i < len(boundaries); i++ {
		assert.Greater(t, boundaries[i], boundaries[i-1], "Boundaries should be sorted")
	}
}

func TestColumnBoundaryDetector_DetectColumnCount_ThreeColumns(t *testing.T) {
	// Simulate 3-column table (like simple case above)
	elements := []*extractor.TextElement{
		newTextElement("A1", 50, 100, 50, 10),
		newTextElement("B1", 150, 100, 50, 10),
		newTextElement("C1", 250, 100, 50, 10),
		newTextElement("A2", 50, 90, 50, 10),
		newTextElement("B2", 150, 90, 50, 10),
		newTextElement("C2", 250, 90, 50, 10),
	}

	detector := NewColumnBoundaryDetector()
	colCount := detector.DetectColumnCount(elements)

	// Should detect 3 columns
	assert.Equal(t, 3, colCount, "Should detect 3 columns")
}

func TestColumnBoundaryDetector_DetectColumnCount_SevenColumns_VTB(t *testing.T) {
	// Simulate VTB table with 7 columns
	// Expected structure from screenshot:
	// Col1: Date/Time (X~50)
	// Col2: Processing Date (X~120)
	// Col3: Amount in operation currency (X~190)
	// Col4: Empty (X~260)
	// Col5: Empty (X~300)
	// Col6: Commission (X~400)
	// Col7: Description (X~470)

	colStarts := []float64{50, 120, 190, 260, 300, 400, 470}
	colWidths := []float64{60, 60, 60, 30, 30, 60, 100}

	elements := []*extractor.TextElement{}

	// Create 10 rows
	for row := 0; row < 10; row++ {
		y := 200.0 - float64(row)*15.0

		// Add elements for each column (skip empty columns sometimes)
		for col := 0; col < len(colStarts); col++ {
			// Skip columns 3 and 4 (empty in VTB)
			if col == 3 || col == 4 {
				continue
			}

			text := "Data"
			elements = append(elements, newTextElement(
				text,
				colStarts[col],
				y,
				colWidths[col],
				10,
			))
		}
	}

	detector := NewColumnBoundaryDetector()
	colCount := detector.DetectColumnCount(elements)

	// Should detect at least 3 columns (simplified simulation)
	// Real VTB data will be tested with actual PDF
	assert.GreaterOrEqual(t, colCount, 3, "Should detect at least 3 columns")
	assert.LessOrEqual(t, colCount, 7, "Should detect at most 7 columns")
}

func TestColumnBoundaryDetector_AssignToColumns(t *testing.T) {
	// Simple 2-column table
	elements := []*extractor.TextElement{
		newTextElement("A1", 50, 100, 50, 10),  // Column 0
		newTextElement("B1", 150, 100, 50, 10), // Column 1
		newTextElement("A2", 50, 90, 50, 10),   // Column 0
		newTextElement("B2", 150, 90, 50, 10),  // Column 1
	}

	detector := NewColumnBoundaryDetector()
	boundaries := detector.DetectBoundaries(elements)
	columnMap := detector.AssignToColumns(elements, boundaries)

	// Should have assigned elements to columns
	assert.Greater(t, len(columnMap), 0, "Should have at least 1 column")

	// Check that all elements were assigned
	totalAssigned := 0
	for _, elems := range columnMap {
		totalAssigned += len(elems)
	}
	assert.Equal(t, len(elements), totalAssigned, "All elements should be assigned")
}

func TestColumnBoundaryDetector_ClusterEdges(t *testing.T) {
	// Edges with 3 clusters: [49, 50, 51], [99, 100, 101], [199, 200, 201]
	edges := []float64{50, 51, 49, 100, 101, 99, 200, 199, 201}

	detector := NewColumnBoundaryDetector()
	clusters := detector.clusterEdges(edges)

	// Should have 3 clusters
	assert.Equal(t, 3, len(clusters), "Should detect 3 clusters")

	// Each cluster should have 3 edges
	for i, cluster := range clusters {
		assert.Equal(t, 3, len(cluster.edges), "Cluster %d should have 3 edges", i)
		assert.Equal(t, 3, cluster.support, "Cluster %d should have support=3", i)
	}

	// Centers should be approximately 50, 100, 200
	assert.InDelta(t, 50.0, clusters[0].center, 1.0, "First cluster center should be ~50")
	assert.InDelta(t, 100.0, clusters[1].center, 1.0, "Second cluster center should be ~100")
	assert.InDelta(t, 200.0, clusters[2].center, 1.0, "Third cluster center should be ~200")
}

func TestColumnBoundaryDetector_AnalyzeTableStructure(t *testing.T) {
	elements := []*extractor.TextElement{
		newTextElement("A", 50, 100, 30, 10),
		newTextElement("B", 150, 100, 40, 10),
		newTextElement("C", 250, 100, 50, 10),
	}

	detector := NewColumnBoundaryDetector()
	analysis := detector.AnalyzeTableStructure(elements)

	assert.NotNil(t, analysis, "Analysis should not be nil")
	assert.Equal(t, 3, analysis.ElementCount, "Should have 3 elements")
	assert.Greater(t, analysis.ColumnCount, 0, "Should detect at least 1 column")
	assert.Greater(t, len(analysis.Boundaries), 0, "Should have at least 1 boundary")
	assert.Equal(t, 50.0, analysis.MinX, "MinX should be 50")
	assert.Equal(t, 250.0, analysis.MaxX, "MaxX should be 250")
	assert.InDelta(t, 40.0, analysis.AvgWidth, 5.0, "AvgWidth should be ~40")
}

func TestColumnBoundaryDetector_EmptyInput(t *testing.T) {
	detector := NewColumnBoundaryDetector()
	boundaries := detector.DetectBoundaries([]*extractor.TextElement{})
	assert.Empty(t, boundaries, "Empty input should return empty boundaries")

	colCount := detector.DetectColumnCount([]*extractor.TextElement{})
	assert.Equal(t, 1, colCount, "Empty input should return 1 column (default)")
}

func TestColumnBoundaryDetector_SingleElement(t *testing.T) {
	elements := []*extractor.TextElement{
		newTextElement("A", 50, 100, 30, 10),
	}

	detector := NewColumnBoundaryDetector()
	boundaries := detector.DetectBoundaries(elements)
	assert.GreaterOrEqual(t, len(boundaries), 1, "Single element should produce at least 1 boundary")

	colCount := detector.DetectColumnCount(elements)
	assert.Equal(t, 1, colCount, "Single element should be 1 column")
}

// Helper function to create test elements
func newTextElement(text string, x, y, width, fontSize float64) *extractor.TextElement {
	return &extractor.TextElement{
		Text:     text,
		X:        x,
		Y:        y,
		Width:    width,
		Height:   fontSize,
		FontSize: fontSize,
		FontName: "Arial",
	}
}
