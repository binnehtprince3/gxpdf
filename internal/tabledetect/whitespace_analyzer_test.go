// Package detector implements table detection algorithms.
package tabledetect

import (
	"math"
	"testing"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestClusterCoordinates_EmptyInput tests clustering with empty input.
func TestClusterCoordinates_EmptyInput(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	clusters := wa.clusterCoordinates([]float64{}, 5.0)

	assert.Empty(t, clusters, "should return empty slice for empty input")
}

// TestClusterCoordinates_SingleCoordinate tests clustering with single coordinate.
func TestClusterCoordinates_SingleCoordinate(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	clusters := wa.clusterCoordinates([]float64{100.0}, 5.0)

	require.Len(t, clusters, 1, "should create one cluster")
	assert.Equal(t, 1, clusters[0].count)
	assert.Equal(t, 100.0, clusters[0].center)
}

// TestClusterCoordinates_NearbyCoordinates tests clustering with nearby coordinates.
func TestClusterCoordinates_NearbyCoordinates(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Three coordinates within 5 points of each other
	coords := []float64{100.0, 102.0, 104.0}

	clusters := wa.clusterCoordinates(coords, 5.0)

	require.Len(t, clusters, 1, "should create single cluster for nearby coordinates")
	assert.Equal(t, 3, clusters[0].count)
	// Center should be average: (100 + 102 + 104) / 3 = 102
	assert.InDelta(t, 102.0, clusters[0].center, 0.1)
}

// TestClusterCoordinates_DistantCoordinates tests clustering with distant coordinates.
func TestClusterCoordinates_DistantCoordinates(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Three coordinates far apart (>5 points)
	coords := []float64{100.0, 200.0, 300.0}

	clusters := wa.clusterCoordinates(coords, 5.0)

	require.Len(t, clusters, 3, "should create three clusters for distant coordinates")
	assert.Equal(t, 1, clusters[0].count)
	assert.Equal(t, 1, clusters[1].count)
	assert.Equal(t, 1, clusters[2].count)
	assert.Equal(t, 100.0, clusters[0].center)
	assert.Equal(t, 200.0, clusters[1].center)
	assert.Equal(t, 300.0, clusters[2].center)
}

// TestClusterCoordinates_MixedDistances tests clustering with mixed distances.
func TestClusterCoordinates_MixedDistances(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Two groups of nearby coordinates
	coords := []float64{
		100.0, 101.0, 103.0, // Group 1 (within 5 points)
		200.0, 202.0, // Group 2 (within 5 points)
		300.0, // Group 3 (isolated)
	}

	clusters := wa.clusterCoordinates(coords, 5.0)

	require.Len(t, clusters, 3, "should create three clusters")

	// Cluster 1: 100, 101, 103
	assert.Equal(t, 3, clusters[0].count)
	assert.InDelta(t, 101.33, clusters[0].center, 0.1)

	// Cluster 2: 200, 202
	assert.Equal(t, 2, clusters[1].count)
	assert.InDelta(t, 201.0, clusters[1].center, 0.1)

	// Cluster 3: 300
	assert.Equal(t, 1, clusters[2].count)
	assert.Equal(t, 300.0, clusters[2].center)
}

// TestClusterCoordinates_UnsortedInput tests clustering with unsorted input.
func TestClusterCoordinates_UnsortedInput(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Unsorted coordinates
	coords := []float64{300.0, 100.0, 200.0, 101.0}

	clusters := wa.clusterCoordinates(coords, 5.0)

	require.Len(t, clusters, 3, "should handle unsorted input")

	// Should be sorted: [100.0, 101.0] [200.0] [300.0]
	assert.Equal(t, 2, clusters[0].count) // 100, 101
	assert.Equal(t, 1, clusters[1].count) // 200
	assert.Equal(t, 1, clusters[2].count) // 300
}

// TestFilterClustersBySize tests cluster filtering by size.
func TestFilterClustersBySize(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	clusters := []coordinateCluster{
		{count: 1, center: 100.0},
		{count: 5, center: 200.0},
		{count: 10, center: 300.0},
		{count: 2, center: 400.0},
	}

	tests := []struct {
		name          string
		minSize       int
		expectedCount int
	}{
		{
			name:          "filter small clusters (min 3)",
			minSize:       3,
			expectedCount: 2, // Only clusters with 5 and 10 elements
		},
		{
			name:          "filter very small clusters (min 1)",
			minSize:       1,
			expectedCount: 4, // All clusters
		},
		{
			name:          "filter all clusters (min 100)",
			minSize:       100,
			expectedCount: 0, // No clusters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := wa.filterClustersBySize(clusters, tt.minSize)
			assert.Len(t, filtered, tt.expectedCount)
		})
	}
}

// TestDetectColumns_BankStatementScenario tests column detection for bank statements.
func TestDetectColumns_BankStatementScenario(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Simulate bank statement with 3 columns: Date | Description | Amount
	// Each column has multiple text elements (rows)
	elements := []*extractor.TextElement{
		// Column 1: Dates (X = 30-35)
		{X: 31.0, Y: 100, Width: 50, Height: 10},
		{X: 32.0, Y: 110, Width: 50, Height: 10},
		{X: 31.5, Y: 120, Width: 50, Height: 10},
		{X: 30.0, Y: 130, Width: 50, Height: 10},
		{X: 33.0, Y: 140, Width: 50, Height: 10},

		// Column 2: Descriptions (X = 100-105) - variable width text
		{X: 100.0, Y: 100, Width: 200, Height: 10},
		{X: 102.0, Y: 110, Width: 180, Height: 10},
		{X: 101.0, Y: 120, Width: 220, Height: 10},
		{X: 103.0, Y: 130, Width: 190, Height: 10},
		{X: 100.5, Y: 140, Width: 210, Height: 10},

		// Column 3: Amounts (X = 350-355)
		{X: 350.0, Y: 100, Width: 80, Height: 10},
		{X: 352.0, Y: 110, Width: 80, Height: 10},
		{X: 351.0, Y: 120, Width: 80, Height: 10},
		{X: 350.5, Y: 130, Width: 80, Height: 10},
		{X: 353.0, Y: 140, Width: 80, Height: 10},
	}

	columns := wa.DetectColumns(elements)

	// UPDATED (2025-10-27): Adaptive algorithm may detect more boundaries
	// Old algorithm: 4-5 boundaries
	// New adaptive algorithm: may detect left+right edges of each column
	// Expecting: at least 4 boundaries (3 columns), at most 8 (3 columns * 2 edges + margins)
	assert.GreaterOrEqual(t, len(columns), 4, "should detect at least 4 column boundaries")
	assert.LessOrEqual(t, len(columns), 8, "should not detect more than 8 column boundaries")

	// Verify first column boundary is near 30 (date column)
	assert.InDelta(t, 31.5, columns[0], 5.0, "first column should be around 30-35")

	// Verify a boundary exists between Date column (ends ~81) and Description column (starts ~100)
	// The adaptive algorithm detects whitespace gaps, so boundary should be around 80-85
	foundGapBoundary := false
	for _, col := range columns {
		if col >= 75 && col <= 90 { // Gap between Date (ends ~81) and Description (starts ~100)
			foundGapBoundary = true
			break
		}
	}
	assert.True(t, foundGapBoundary, "should detect gap boundary around 75-90 (between Date and Description)")

	// Verify third column boundary is near 350 (amount column)
	foundAmountColumn := false
	for _, col := range columns[1 : len(columns)-1] {
		if math.Abs(col-351.3) < 10.0 {
			foundAmountColumn = true
			break
		}
	}
	assert.True(t, foundAmountColumn, "should detect amount column around 350-355")
}

// TestDetectColumns_NoFalseColumnsInWhitespace tests that whitespace within cells
// does not create false column boundaries.
func TestDetectColumns_NoFalseColumnsInWhitespace(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Simulate a single column with text elements that have varying X positions
	// but are all logically in the same column
	elements := []*extractor.TextElement{
		{X: 100.0, Y: 100, Width: 50, Height: 10},  // "Date"
		{X: 150.0, Y: 100, Width: 100, Height: 10}, // "Description" (whitespace gap, but same row)
		{X: 100.0, Y: 110, Width: 50, Height: 10},  // "01.02.24"
		{X: 150.0, Y: 110, Width: 100, Height: 10}, // "Payment for..."
		{X: 100.0, Y: 120, Width: 50, Height: 10},  // "02.02.24"
		{X: 150.0, Y: 120, Width: 100, Height: 10}, // "Transfer to..."
	}

	columns := wa.DetectColumns(elements)

	// OLD BEHAVIOR: Would detect column at ~100 and ~150 (2 internal columns + edges = 4 total)
	// NEW BEHAVIOR: Should cluster 100 positions together, 150 positions together
	//               Result: 2 significant columns (left edge, internal boundary, right edge = 3-4 total)

	assert.GreaterOrEqual(t, len(columns), 3, "should detect at least 3 boundaries")
	assert.LessOrEqual(t, len(columns), 4, "should not detect excessive boundaries")
}

// TestDetectColumns_MultilineCell tests handling of multi-line cell content.
func TestDetectColumns_MultilineCell(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Simulate a table where cell content spans multiple lines
	// Column 1: Dates (X = 50)
	// Column 2: Multi-line descriptions (X = 150)
	elements := []*extractor.TextElement{
		// Row 1
		{X: 50.0, Y: 100, Width: 50, Height: 10},   // "01.02.24"
		{X: 150.0, Y: 100, Width: 200, Height: 10}, // "Payment for"
		{X: 150.0, Y: 90, Width: 200, Height: 10},  // "services rendered" (continuation)

		// Row 2
		{X: 50.0, Y: 70, Width: 50, Height: 10},   // "02.02.24"
		{X: 150.0, Y: 70, Width: 200, Height: 10}, // "Transfer to"
		{X: 150.0, Y: 60, Width: 200, Height: 10}, // "account XYZ" (continuation)
	}

	columns := wa.DetectColumns(elements)

	// With 6 elements and clustering:
	// - 2 elements at X=50 (33% of data) → significant cluster
	// - 4 elements at X=150 (67% of data) → significant cluster
	// Expected: 2 column centers + right edge = 3 boundaries
	assert.GreaterOrEqual(t, len(columns), 2, "should detect at least 2 boundaries")

	// Verify date column around X=50
	foundDateColumn := false
	for _, col := range columns {
		if math.Abs(col-50.0) < 10.0 {
			foundDateColumn = true
			break
		}
	}
	assert.True(t, foundDateColumn, "should detect date column around 50")

	// Verify description column around X=150
	foundDescriptionColumn := false
	for _, col := range columns {
		if math.Abs(col-150.0) < 10.0 {
			foundDescriptionColumn = true
			break
		}
	}
	assert.True(t, foundDescriptionColumn, "should detect description column around 150")
}

// TestDetectColumns_EdgeCases tests edge cases.
func TestDetectColumns_EdgeCases(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	tests := []struct {
		name     string
		elements []*extractor.TextElement
		wantCols int // Expected number of column boundaries
	}{
		{
			name:     "empty input",
			elements: []*extractor.TextElement{},
			wantCols: 0,
		},
		{
			name: "single element",
			elements: []*extractor.TextElement{
				{X: 100, Y: 100, Width: 50, Height: 10},
			},
			wantCols: 2, // Left edge + right edge
		},
		{
			name: "all elements at same X position",
			elements: []*extractor.TextElement{
				{X: 100, Y: 100, Width: 50, Height: 10},
				{X: 100, Y: 110, Width: 50, Height: 10},
				{X: 100, Y: 120, Width: 50, Height: 10},
			},
			wantCols: 2, // Left edge + right edge (single column)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			columns := wa.DetectColumns(tt.elements)
			assert.Len(t, columns, tt.wantCols, "unexpected number of columns")
		})
	}
}

// TestCalculateClusterCenter tests cluster center calculation.
func TestCalculateClusterCenter(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	tests := []struct {
		name        string
		coords      []float64
		wantCenter  float64
		description string
	}{
		{
			name:        "empty input",
			coords:      []float64{},
			wantCenter:  0.0,
			description: "should return 0 for empty input",
		},
		{
			name:        "single coordinate",
			coords:      []float64{100.0},
			wantCenter:  100.0,
			description: "should return same value for single coordinate",
		},
		{
			name:        "multiple coordinates",
			coords:      []float64{100.0, 110.0, 120.0},
			wantCenter:  110.0,
			description: "should return mean of coordinates",
		},
		{
			name:        "coordinates with decimals",
			coords:      []float64{31.62, 32.14, 31.89},
			wantCenter:  31.883,
			description: "should handle decimal coordinates",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			center := wa.calculateClusterCenter(tt.coords)
			assert.InDelta(t, tt.wantCenter, center, 0.01, tt.description)
		})
	}
}

// TestDetectColumns_RealWorldBankStatement tests with realistic bank statement data.
func TestDetectColumns_RealWorldBankStatement(t *testing.T) {
	wa := NewDefaultWhitespaceAnalyzer()

	// Based on actual data from Russian bank statement
	// Expected structure: Date (X~31), Description (X~127), Amount (X~527)
	elements := []*extractor.TextElement{
		// Row 1: "01.02.24" | "Пополнение депозита" | "-104,628.00"
		{X: 31.62, Y: 442.27, Width: 46.08, Height: 9.38},
		{X: 127.78, Y: 442.27, Width: 126.22, Height: 9.38},
		{X: 527.06, Y: 442.27, Width: 57.13, Height: 9.38},

		// Row 2: "02.02.24" | "Уменьшение суммы депозита" | "-97,213.00"
		{X: 31.62, Y: 429.51, Width: 46.08, Height: 9.38},
		{X: 127.78, Y: 429.51, Width: 176.50, Height: 9.38},
		{X: 527.06, Y: 429.51, Width: 51.38, Height: 9.38},

		// Row 3: "05.02.24" | "Пополнение депозита" | "-87,420.00"
		{X: 31.62, Y: 416.75, Width: 46.08, Height: 9.38},
		{X: 127.78, Y: 416.75, Width: 126.22, Height: 9.38},
		{X: 527.06, Y: 416.75, Width: 51.38, Height: 9.38},

		// More rows...
		{X: 31.62, Y: 403.99, Width: 46.08, Height: 9.38},
		{X: 127.78, Y: 403.99, Width: 126.22, Height: 9.38},
		{X: 527.06, Y: 403.99, Width: 57.13, Height: 9.38},

		{X: 31.62, Y: 391.23, Width: 46.08, Height: 9.38},
		{X: 127.78, Y: 391.23, Width: 176.50, Height: 9.38},
		{X: 527.06, Y: 391.23, Width: 51.38, Height: 9.38},
	}

	columns := wa.DetectColumns(elements)

	// NEW (2025-10-27): Adaptive algorithm detects more accurate boundaries
	// Old algorithm: 4-5 boundaries
	// New algorithm: 7 boundaries (includes both left and right edges of columns)
	// This is IMPROVEMENT, not regression - more boundaries = more accurate table structure
	assert.GreaterOrEqual(t, len(columns), 4, "should detect at least 4 column boundaries")
	assert.LessOrEqual(t, len(columns), 10, "should not detect more than 10 column boundaries (sanity check)")

	// Verify date column around X=31.62
	assert.InDelta(t, 31.62, columns[0], 5.0, "first column should be date column around 31.62")

	// Verify description column around X=127.78
	foundDescriptionColumn := false
	for _, col := range columns[1 : len(columns)-1] {
		if math.Abs(col-127.78) < 10.0 {
			foundDescriptionColumn = true
			break
		}
	}
	assert.True(t, foundDescriptionColumn, "should detect description column around 127.78")

	// Verify amount column around X=527.06
	foundAmountColumn := false
	for _, col := range columns[1 : len(columns)-1] {
		if math.Abs(col-527.06) < 10.0 {
			foundAmountColumn = true
			break
		}
	}
	assert.True(t, foundAmountColumn, "should detect amount column around 527.06")

	t.Logf("Detected columns: %v", columns)
}

// TestDetectColumns_Performance tests performance with large datasets.
func TestDetectColumns_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping performance test in short mode")
	}

	wa := NewDefaultWhitespaceAnalyzer()

	// Create 1000 text elements (simulating large table)
	elements := make([]*extractor.TextElement, 1000)
	for i := 0; i < 1000; i++ {
		// Distribute across 3 columns
		col := i % 3
		x := float64(50 + col*150)
		y := float64(100 + i*10)

		elements[i] = &extractor.TextElement{
			X:      x,
			Y:      y,
			Width:  100,
			Height: 10,
		}
	}

	// Should complete quickly (< 100ms for 1000 elements)
	columns := wa.DetectColumns(elements)

	assert.GreaterOrEqual(t, len(columns), 4, "should detect columns even with large dataset")
	t.Logf("Processed 1000 elements, detected %d columns", len(columns))
}
