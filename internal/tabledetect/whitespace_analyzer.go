// Package detector implements table detection algorithms.
package tabledetect

import (
	"fmt"
	"math"
	"sort"

	"github.com/coregx/gxpdf/internal/extractor"
)

// DefaultWhitespaceAnalyzer analyzes whitespace distribution to find table structure.
//
// This is the default implementation of the WhitespaceAnalyzer interface.
// It is used for stream mode table detection, where tables don't have
// visible ruling lines. Instead, we infer table structure from:
//   - Vertical whitespace (gaps between rows)
//   - Horizontal whitespace (gaps between columns)
//   - Text alignment patterns
//
// Algorithm inspired by tabula-java's BasicExtractionAlgorithm.
// Reference: tabula-java/technology/tabula/extractors/BasicExtractionAlgorithm.java
type DefaultWhitespaceAnalyzer struct {
	minGapWidth            float64 // Minimum gap width to consider (in points)
	alignmentTolerance     float64 // Tolerance for text alignment (in points)
	projectionAnalyzer     ProjectionAnalyzer
	columnBoundaryDetector *ColumnBoundaryDetector // Adaptive column detector (2025-10-27 VTB multi-line fix)
	useAdaptiveColumns     bool                    // Enable adaptive column detection (default: true)
	isLatticeMode          bool                    // Lattice mode (true) vs Stream mode (false) - affects row detection threshold
}

// NewDefaultWhitespaceAnalyzer creates a new DefaultWhitespaceAnalyzer with default settings.
func NewDefaultWhitespaceAnalyzer() *DefaultWhitespaceAnalyzer {
	return &DefaultWhitespaceAnalyzer{
		minGapWidth:            10.0, // Minimum 10 points (~3.5mm)
		alignmentTolerance:     2.0,  // 2 points tolerance
		projectionAnalyzer:     NewDefaultProjectionAnalyzer(),
		columnBoundaryDetector: NewColumnBoundaryDetector(), // NEW: Adaptive column detection
		useAdaptiveColumns:     true,                        // NEW: Enabled by default
		isLatticeMode:          false,                       // Stream mode by default (Alfa-Bank)
	}
}

// NewWhitespaceAnalyzer creates a new DefaultWhitespaceAnalyzer with default settings.
// Deprecated: Use NewDefaultWhitespaceAnalyzer instead. Kept for backward compatibility.
func NewWhitespaceAnalyzer() *DefaultWhitespaceAnalyzer {
	return NewDefaultWhitespaceAnalyzer()
}

// NewWhitespaceAnalyzerForLattice creates a WhitespaceAnalyzer for Lattice mode (VTB, Sberbank).
// Lattice mode: tables have visible ruling lines (Grid-based extraction).
// Uses larger gap threshold (2.0x fontSize) to ignore gaps within multi-line cells.
func NewWhitespaceAnalyzerForLattice() *DefaultWhitespaceAnalyzer {
	return &DefaultWhitespaceAnalyzer{
		minGapWidth:            10.0,
		alignmentTolerance:     2.0,
		projectionAnalyzer:     NewDefaultProjectionAnalyzer(),
		columnBoundaryDetector: NewColumnBoundaryDetector(),
		useAdaptiveColumns:     true,
		isLatticeMode:          true, // Lattice mode (VTB, Sberbank)
	}
}

// WithMinGapWidth sets the minimum gap width.
func (wa *DefaultWhitespaceAnalyzer) WithMinGapWidth(width float64) *DefaultWhitespaceAnalyzer {
	wa.minGapWidth = width
	return wa
}

// WithAlignmentTolerance sets the alignment tolerance.
func (wa *DefaultWhitespaceAnalyzer) WithAlignmentTolerance(tol float64) *DefaultWhitespaceAnalyzer {
	wa.alignmentTolerance = tol
	return wa
}

// WithProjectionAnalyzer sets a custom projection analyzer.
func (wa *DefaultWhitespaceAnalyzer) WithProjectionAnalyzer(analyzer ProjectionAnalyzer) *DefaultWhitespaceAnalyzer {
	wa.projectionAnalyzer = analyzer
	return wa
}

// DetectColumns finds vertical alignment patterns (column boundaries).
//
// Returns a slice of X coordinates representing column boundaries,
// sorted left to right.
//
// Algorithm:
//  1. Collect all left edges (X coordinates) from text elements
//  2. Cluster nearby coordinates using tolerance-based clustering
//  3. Filter clusters by minimum element count (significance threshold)
//  4. Return cluster centers as column boundaries
//
// This approach reduces false column boundaries caused by whitespace gaps
// within cell content, especially for tables with variable-width text.
func (wa *DefaultWhitespaceAnalyzer) DetectColumns(elements []*extractor.TextElement) []float64 {
	if len(elements) == 0 {
		return []float64{}
	}

	// NEW (2025-10-27): Use adaptive column detection if enabled
	// This solves VTB multi-line cell problem (see ANALYSIS_VTB_TABLE_MULTI_LINE_CELLS.md)
	if wa.useAdaptiveColumns && wa.columnBoundaryDetector != nil {
		return wa.columnBoundaryDetector.DetectBoundaries(elements)
	}

	// OLD ALGORITHM (fallback for backward compatibility)

	// Collect all left edges (X coordinates)
	leftEdges := make([]float64, 0, len(elements))
	for _, elem := range elements {
		leftEdges = append(leftEdges, elem.X)
	}

	// Cluster coordinates with tolerance of 5 points
	// This groups text elements that start at similar X positions
	clusterTolerance := 5.0
	clusters := wa.clusterCoordinates(leftEdges, clusterTolerance)

	// Filter clusters by minimum element count
	// Only keep clusters with at least 10% of total elements
	// This eliminates noise and keeps only significant column boundaries
	// Minimum 1 element for very small datasets, maximum 3 for small datasets
	minClusterSize := int(math.Max(1, math.Min(3, float64(len(elements))*0.1)))
	significantClusters := wa.filterClustersBySize(clusters, minClusterSize)

	// Extract column boundaries from cluster centers
	var columns []float64

	// Edge case: no significant clusters found
	if len(significantClusters) == 0 {
		// Fallback: just use min and max X coordinates
		minX := elements[0].X
		maxX := elements[0].Right()
		for _, elem := range elements {
			if elem.X < minX {
				minX = elem.X
			}
			if elem.Right() > maxX {
				maxX = elem.Right()
			}
		}
		return []float64{minX, maxX}
	}

	// Add left edge of content area (first cluster)
	columns = append(columns, significantClusters[0].center)

	// Add cluster centers as column boundaries (skip first, we already added it)
	for i := 1; i < len(significantClusters); i++ {
		columns = append(columns, significantClusters[i].center)
	}

	// Add right edge of content area
	// Find max X coordinate (rightmost text element)
	maxX := elements[0].Right()
	for _, elem := range elements {
		if elem.Right() > maxX {
			maxX = elem.Right()
		}
	}
	columns = append(columns, maxX)

	// Sort columns left to right
	sort.Float64s(columns)

	return columns
}

// DetectColumnsWithRulingLines uses HYBRID approach - combining text and graphics.
//
// NEW (2025-10-27): User requirement - "грамотно доработать" (properly refine)
//
// This method combines:
// - Text-based edge clustering (accurate but unstable)
// - Graphics-based ruling lines (stable but may have extra lines)
//
// Expected improvement: 66.7% → 90%+
func (wa *DefaultWhitespaceAnalyzer) DetectColumnsWithRulingLines(
	elements []*extractor.TextElement,
	rulingLineXPositions []float64,
) []float64 {
	if len(elements) == 0 {
		return []float64{}
	}

	// Use hybrid approach if adaptive detector is available
	if wa.useAdaptiveColumns && wa.columnBoundaryDetector != nil {
		return wa.columnBoundaryDetector.DetectBoundariesWithRulingLines(elements, rulingLineXPositions)
	}

	// Fallback: just use text-based detection
	return wa.DetectColumns(elements)
}

// DetectRows finds horizontal alignment patterns (row boundaries).
//
// Returns a slice of Y coordinates representing row boundaries,
// sorted bottom to top (in PDF coordinates).
//
// Uses hybrid approach (Gap + Overlap + Alignment) for maximum accuracy.
func (wa *DefaultWhitespaceAnalyzer) DetectRows(elements []*extractor.TextElement) []float64 {
	// Use hybrid approach for universal table detection
	return wa.DetectRowsHybrid(elements)
}

// findVerticalAlignments finds X coordinates where text elements align vertically.
//
// This helps detect columns in tables where text is aligned.
func (wa *DefaultWhitespaceAnalyzer) findVerticalAlignments(elements []*extractor.TextElement) []float64 {
	// Collect left edges, right edges, and centers
	type edgeCount struct {
		position float64
		count    int
	}

	leftEdges := make(map[int]float64)
	rightEdges := make(map[int]float64)

	for _, elem := range elements {
		// Left edge
		leftKey := int(math.Round(elem.X / wa.alignmentTolerance))
		leftEdges[leftKey] = elem.X

		// Right edge
		rightKey := int(math.Round(elem.Right() / wa.alignmentTolerance))
		rightEdges[rightKey] = elem.Right()
	}

	// Find edges with multiple elements
	var alignments []float64

	// Count left edges
	leftCounts := wa.countEdgeOccurrences(elements, true)
	for key, pos := range leftEdges {
		if leftCounts[key] >= 3 { // At least 3 elements aligned
			alignments = append(alignments, pos)
		}
	}

	// Count right edges
	rightCounts := wa.countEdgeOccurrences(elements, false)
	for key, pos := range rightEdges {
		if rightCounts[key] >= 3 { // At least 3 elements aligned
			alignments = append(alignments, pos)
		}
	}

	return alignments
}

// findHorizontalAlignments finds Y coordinates where text elements align horizontally.
//
// This helps detect rows in tables where text is aligned.
func (wa *DefaultWhitespaceAnalyzer) findHorizontalAlignments(elements []*extractor.TextElement) []float64 {
	// Collect bottom edges and top edges
	bottomEdges := make(map[int]float64)
	topEdges := make(map[int]float64)

	for _, elem := range elements {
		// Bottom edge
		bottomKey := int(math.Round(elem.Y / wa.alignmentTolerance))
		bottomEdges[bottomKey] = elem.Y

		// Top edge
		topKey := int(math.Round(elem.Top() / wa.alignmentTolerance))
		topEdges[topKey] = elem.Top()
	}

	// Find edges with multiple elements
	var alignments []float64

	// Count bottom edges
	bottomCounts := wa.countHorizontalOccurrences(elements, true)
	for key, pos := range bottomEdges {
		if bottomCounts[key] >= 3 { // At least 3 elements aligned
			alignments = append(alignments, pos)
		}
	}

	// Count top edges
	topCounts := wa.countHorizontalOccurrences(elements, false)
	for key, pos := range topEdges {
		if topCounts[key] >= 3 { // At least 3 elements aligned
			alignments = append(alignments, pos)
		}
	}

	return alignments
}

// countEdgeOccurrences counts how many elements have edges at each position.
func (wa *DefaultWhitespaceAnalyzer) countEdgeOccurrences(elements []*extractor.TextElement, left bool) map[int]int {
	counts := make(map[int]int)

	for _, elem := range elements {
		var key int
		if left {
			key = int(math.Round(elem.X / wa.alignmentTolerance))
		} else {
			key = int(math.Round(elem.Right() / wa.alignmentTolerance))
		}
		counts[key]++
	}

	return counts
}

// countHorizontalOccurrences counts how many elements have edges at each position.
func (wa *DefaultWhitespaceAnalyzer) countHorizontalOccurrences(elements []*extractor.TextElement, bottom bool) map[int]int {
	counts := make(map[int]int)

	for _, elem := range elements {
		var key int
		if bottom {
			key = int(math.Round(elem.Y / wa.alignmentTolerance))
		} else {
			key = int(math.Round(elem.Top() / wa.alignmentTolerance))
		}
		counts[key]++
	}

	return counts
}

// uniqueAndSort removes duplicate coordinates and sorts them.
func (wa *DefaultWhitespaceAnalyzer) uniqueAndSort(coords []float64, tolerance float64) []float64 {
	if len(coords) == 0 {
		return coords
	}

	// Sort first
	sort.Float64s(coords)

	// Remove duplicates within tolerance
	unique := []float64{coords[0]}

	for i := 1; i < len(coords); i++ {
		// Check if this coordinate is different from the last unique one
		if math.Abs(coords[i]-unique[len(unique)-1]) > tolerance {
			unique = append(unique, coords[i])
		}
	}

	return unique
}

// GroupIntoRows groups text elements into rows based on Y position.
//
// Returns a slice of rows, each containing text elements on that row.
func (wa *DefaultWhitespaceAnalyzer) GroupIntoRows(elements []*extractor.TextElement) [][]*extractor.TextElement {
	if len(elements) == 0 {
		return [][]*extractor.TextElement{}
	}

	// Sort elements by Y coordinate (top to bottom)
	sorted := make([]*extractor.TextElement, len(elements))
	copy(sorted, elements)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Y > sorted[j].Y // Higher Y first (top to bottom)
	})

	// Group into rows
	var rows [][]*extractor.TextElement
	currentRow := []*extractor.TextElement{sorted[0]}
	currentY := sorted[0].Y

	for i := 1; i < len(sorted); i++ {
		elem := sorted[i]

		// Check if element is on the same row
		if math.Abs(elem.Y-currentY) <= wa.alignmentTolerance {
			// Same row
			currentRow = append(currentRow, elem)
		} else {
			// New row
			rows = append(rows, currentRow)
			currentRow = []*extractor.TextElement{elem}
			currentY = elem.Y
		}
	}

	// Add last row
	if len(currentRow) > 0 {
		rows = append(rows, currentRow)
	}

	return rows
}

// GroupIntoColumns groups text elements into columns based on X position.
//
// Returns a slice of columns, each containing text elements in that column.
func (wa *DefaultWhitespaceAnalyzer) GroupIntoColumns(elements []*extractor.TextElement) [][]*extractor.TextElement {
	if len(elements) == 0 {
		return [][]*extractor.TextElement{}
	}

	// Sort elements by X coordinate (left to right)
	sorted := make([]*extractor.TextElement, len(elements))
	copy(sorted, elements)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].X < sorted[j].X
	})

	// Group into columns
	var columns [][]*extractor.TextElement
	currentColumn := []*extractor.TextElement{sorted[0]}
	currentX := sorted[0].X

	for i := 1; i < len(sorted); i++ {
		elem := sorted[i]

		// Check if element is in the same column
		if math.Abs(elem.X-currentX) <= wa.alignmentTolerance {
			// Same column
			currentColumn = append(currentColumn, elem)
		} else {
			// New column
			columns = append(columns, currentColumn)
			currentColumn = []*extractor.TextElement{elem}
			currentX = elem.X
		}
	}

	// Add last column
	if len(currentColumn) > 0 {
		columns = append(columns, currentColumn)
	}

	return columns
}

// DetectTableRegion detects a table region based on whitespace analysis.
//
// Returns the bounding rectangle of the detected table, or nil if no table.
func (wa *DefaultWhitespaceAnalyzer) DetectTableRegion(elements []*extractor.TextElement) *extractor.Rectangle {
	if len(elements) == 0 {
		return nil
	}

	// Detect rows and columns
	rows := wa.DetectRows(elements)
	columns := wa.DetectColumns(elements)

	// Need at least 2 rows and 2 columns for a table
	if len(rows) < 2 || len(columns) < 2 {
		return nil
	}

	// Create bounding rectangle
	minX := columns[0]
	maxX := columns[len(columns)-1]
	minY := rows[0]
	maxY := rows[len(rows)-1]

	rect := extractor.NewRectangle(minX, minY, maxX-minX, maxY-minY)
	return &rect
}

// String returns a string representation of the analyzer.
func (wa *DefaultWhitespaceAnalyzer) String() string {
	return fmt.Sprintf("WhitespaceAnalyzer{minGap=%.2f, tolerance=%.2f}",
		wa.minGapWidth, wa.alignmentTolerance)
}

// calculateAverageFontSize calculates the average font size of text elements.
//
// This is used for adaptive gap detection to distinguish:
//   - Gaps within multi-line cells (~1.2-1.5x font size)
//   - Gaps between table rows (~2-3x font size)
func (wa *DefaultWhitespaceAnalyzer) calculateAverageFontSize(elements []*extractor.TextElement) float64 {
	if len(elements) == 0 {
		return 10.0 // Default fallback
	}

	sum := 0.0
	count := 0
	for _, elem := range elements {
		if elem.FontSize > 0 {
			sum += elem.FontSize
			count++
		}
	}

	if count == 0 {
		return 10.0 // Default fallback
	}

	return sum / float64(count)
}

// coordinateCluster represents a cluster of nearby coordinates.
//
// Used for identifying significant column/row boundaries by grouping
// text elements with similar positions.
type coordinateCluster struct {
	coordinates []float64 // All coordinates in this cluster
	center      float64   // Center (mean) of cluster
	count       int       // Number of coordinates in cluster
}

// clusterCoordinates groups nearby coordinates into clusters.
//
// Algorithm:
//  1. Sort coordinates
//  2. Start new cluster with first coordinate
//  3. For each coordinate:
//     - If within tolerance of current cluster, add to cluster
//     - Otherwise, start new cluster
//  4. Calculate cluster centers (mean of coordinates)
//
// Parameters:
//   - coords: Coordinates to cluster
//   - tolerance: Maximum distance between coordinates in same cluster
//
// Returns slice of clusters, sorted by center position.
func (wa *DefaultWhitespaceAnalyzer) clusterCoordinates(coords []float64, tolerance float64) []coordinateCluster {
	if len(coords) == 0 {
		return []coordinateCluster{}
	}

	// Sort coordinates
	sorted := make([]float64, len(coords))
	copy(sorted, coords)
	sort.Float64s(sorted)

	// Build clusters
	var clusters []coordinateCluster
	currentCluster := coordinateCluster{
		coordinates: []float64{sorted[0]},
		count:       1,
	}

	for i := 1; i < len(sorted); i++ {
		coord := sorted[i]

		// Check if coordinate is within tolerance of current cluster
		// Use the last coordinate in cluster as reference
		lastInCluster := currentCluster.coordinates[len(currentCluster.coordinates)-1]

		if math.Abs(coord-lastInCluster) <= tolerance {
			// Add to current cluster
			currentCluster.coordinates = append(currentCluster.coordinates, coord)
			currentCluster.count++
		} else {
			// Calculate center of current cluster
			currentCluster.center = wa.calculateClusterCenter(currentCluster.coordinates)

			// Save current cluster and start new one
			clusters = append(clusters, currentCluster)
			currentCluster = coordinateCluster{
				coordinates: []float64{coord},
				count:       1,
			}
		}
	}

	// Add final cluster
	currentCluster.center = wa.calculateClusterCenter(currentCluster.coordinates)
	clusters = append(clusters, currentCluster)

	return clusters
}

// calculateClusterCenter calculates the center (mean) of a cluster.
func (wa *DefaultWhitespaceAnalyzer) calculateClusterCenter(coords []float64) float64 {
	if len(coords) == 0 {
		return 0
	}

	sum := 0.0
	for _, coord := range coords {
		sum += coord
	}
	return sum / float64(len(coords))
}

// filterClustersBySize filters clusters by minimum element count.
//
// Only clusters with at least minSize elements are kept.
// This eliminates noise and keeps only significant clusters.
func (wa *DefaultWhitespaceAnalyzer) filterClustersBySize(clusters []coordinateCluster, minSize int) []coordinateCluster {
	var filtered []coordinateCluster

	for _, cluster := range clusters {
		if cluster.count >= minSize {
			filtered = append(filtered, cluster)
		}
	}

	return filtered
}
