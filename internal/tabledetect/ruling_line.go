// Package detector implements table detection algorithms.
//
// This is the Application layer in DDD/Clean Architecture.
// It uses domain logic and extracted content to detect table regions.
package tabledetect

import (
	"fmt"
	"math"

	"github.com/coregx/gxpdf/internal/extractor"
)

// RulingLine represents a horizontal or vertical line in a PDF.
//
// Ruling lines are used in lattice mode table extraction to detect
// table boundaries and cell grids.
//
// Inspired by tabula-java's Ruling class.
// Reference: tabula-java/technology/tabula/Ruling.java
type RulingLine struct {
	Start        extractor.Point // Start point
	End          extractor.Point // End point
	IsHorizontal bool            // True if horizontal, false if vertical
}

// NewRulingLine creates a new RulingLine.
func NewRulingLine(start, end extractor.Point) *RulingLine {
	// Determine if horizontal or vertical
	isHorizontal := math.Abs(start.Y-end.Y) < math.Abs(start.X-end.X)

	return &RulingLine{
		Start:        start,
		End:          end,
		IsHorizontal: isHorizontal,
	}
}

// Length returns the length of the ruling line.
func (rl *RulingLine) Length() float64 {
	dx := rl.End.X - rl.Start.X
	dy := rl.End.Y - rl.Start.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Intersects checks if this ruling line intersects with another.
//
// Returns the intersection point, or nil if no intersection.
func (rl *RulingLine) Intersects(other *RulingLine) *extractor.Point {
	// Horizontal and vertical lines are most common case
	if rl.IsHorizontal && !other.IsHorizontal {
		// This is horizontal, other is vertical
		return rl.intersectHorizontalVertical(other)
	} else if !rl.IsHorizontal && other.IsHorizontal {
		// This is vertical, other is horizontal
		return other.intersectHorizontalVertical(rl)
	}

	// Parallel lines or both oblique - no intersection for our purposes
	return nil
}

// intersectHorizontalVertical finds intersection of horizontal and vertical lines.
func (rl *RulingLine) intersectHorizontalVertical(vertical *RulingLine) *extractor.Point {
	// rl is horizontal, vertical is vertical
	// Intersection at (vertical.X, rl.Y)
	x := vertical.Start.X // Vertical line has same X
	y := rl.Start.Y       // Horizontal line has same Y

	// Check if intersection point is within both line segments
	if x >= math.Min(rl.Start.X, rl.End.X) && x <= math.Max(rl.Start.X, rl.End.X) &&
		y >= math.Min(vertical.Start.Y, vertical.End.Y) && y <= math.Max(vertical.Start.Y, vertical.End.Y) {
		point := extractor.NewPoint(x, y)
		return &point
	}

	return nil
}

// String returns a string representation of the ruling line.
func (rl *RulingLine) String() string {
	orientation := "V"
	if rl.IsHorizontal {
		orientation = "H"
	}
	return fmt.Sprintf("RulingLine{%s, start=%s, end=%s, len=%.2f}",
		orientation, rl.Start.String(), rl.End.String(), rl.Length())
}

// DefaultRulingLineDetector detects ruling lines from graphics elements.
//
// This is the default implementation of the RulingLineDetector interface.
// It is used for lattice mode table detection, where tables have
// visible borders and grid lines.
//
// Algorithm inspired by tabula-java's SpreadsheetExtractionAlgorithm.
// Reference: tabula-java/technology/tabula/extractors/SpreadsheetExtractionAlgorithm.java
type DefaultRulingLineDetector struct {
	minLineLength float64 // Minimum line length to consider (in points)
	tolerance     float64 // Tolerance for alignment (in points)
}

// NewDefaultRulingLineDetector creates a new DefaultRulingLineDetector with default settings.
func NewDefaultRulingLineDetector() *DefaultRulingLineDetector {
	return &DefaultRulingLineDetector{
		minLineLength: 10.0, // Minimum 10 points (about 3.5mm)
		tolerance:     2.0,  // 2 points tolerance for alignment
	}
}

// NewRulingLineDetector creates a new DefaultRulingLineDetector with default settings.
// Deprecated: Use NewDefaultRulingLineDetector instead. Kept for backward compatibility.
func NewRulingLineDetector() *DefaultRulingLineDetector {
	return NewDefaultRulingLineDetector()
}

// WithMinLineLength sets the minimum line length.
func (d *DefaultRulingLineDetector) WithMinLineLength(length float64) *DefaultRulingLineDetector {
	d.minLineLength = length
	return d
}

// WithTolerance sets the alignment tolerance.
func (d *DefaultRulingLineDetector) WithTolerance(tol float64) *DefaultRulingLineDetector {
	d.tolerance = tol
	return d
}

// DetectRulingLines extracts horizontal and vertical lines from graphics.
//
// Returns a slice of RulingLines sorted by position.
func (d *DefaultRulingLineDetector) DetectRulingLines(graphics []*extractor.GraphicsElement) ([]*RulingLine, error) {
	var lines []*RulingLine

	// Extract lines from graphics elements
	for _, elem := range graphics {
		// Only process line elements
		if elem.Type != extractor.GraphicsTypeLine {
			continue
		}

		// Must have exactly 2 points
		if len(elem.Points) != 2 {
			continue
		}

		start := elem.Points[0]
		end := elem.Points[1]

		// Check if horizontal or vertical (within tolerance)
		isHorizontal := math.Abs(start.Y-end.Y) <= d.tolerance
		isVertical := math.Abs(start.X-end.X) <= d.tolerance

		if !isHorizontal && !isVertical {
			// Oblique line - skip
			continue
		}

		// Normalize line to make horizontal/vertical exact
		if isHorizontal {
			end.Y = start.Y // Make exactly horizontal
		} else if isVertical {
			end.X = start.X // Make exactly vertical
		}

		// Create ruling line
		line := NewRulingLine(start, end)

		// Check minimum length
		if line.Length() < d.minLineLength {
			continue
		}

		lines = append(lines, line)
	}

	// Merge collinear lines
	lines = d.mergeCollinearLines(lines)

	return lines, nil
}

// mergeCollinearLines merges lines that are on the same axis.
//
// This handles cases where a single logical line is represented as
// multiple line segments in the PDF.
//
// Algorithm inspired by tabula-java's Ruling.collapseOrientedRulings().
func (d *DefaultRulingLineDetector) mergeCollinearLines(lines []*RulingLine) []*RulingLine {
	if len(lines) == 0 {
		return lines
	}

	// Separate horizontal and vertical lines
	var horizontal, vertical []*RulingLine
	for _, line := range lines {
		if line.IsHorizontal {
			horizontal = append(horizontal, line)
		} else {
			vertical = append(vertical, line)
		}
	}

	// Merge each group
	horizontal = d.mergeGroup(horizontal, true)
	vertical = d.mergeGroup(vertical, false)

	// Combine and return
	result := make([]*RulingLine, 0, len(horizontal)+len(vertical))
	result = append(result, horizontal...)
	result = append(result, vertical...)

	return result
}

// mergeGroup merges lines in a group (all horizontal or all vertical).
func (d *DefaultRulingLineDetector) mergeGroup(lines []*RulingLine, isHorizontal bool) []*RulingLine {
	if len(lines) == 0 {
		return lines
	}

	// Group lines by their position on the perpendicular axis
	// For horizontal lines, group by Y coordinate
	// For vertical lines, group by X coordinate
	groups := make(map[int][]*RulingLine)

	for _, line := range lines {
		var key int
		if isHorizontal {
			key = int(math.Round(line.Start.Y / d.tolerance))
		} else {
			key = int(math.Round(line.Start.X / d.tolerance))
		}

		groups[key] = append(groups[key], line)
	}

	// Merge lines within each group
	var result []*RulingLine
	for _, group := range groups {
		if len(group) == 1 {
			result = append(result, group[0])
			continue
		}

		// Sort lines in group by position on main axis
		// For horizontal lines, sort by X
		// For vertical lines, sort by Y
		d.sortLines(group, isHorizontal)

		// Merge adjacent lines
		merged := d.mergeAdjacent(group, isHorizontal)
		result = append(result, merged...)
	}

	return result
}

// sortLines sorts lines by their position on the main axis.
func (d *DefaultRulingLineDetector) sortLines(lines []*RulingLine, isHorizontal bool) {
	// Simple bubble sort (small arrays)
	n := len(lines)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			var swap bool
			if isHorizontal {
				// Sort by X coordinate
				swap = lines[j].Start.X > lines[j+1].Start.X
			} else {
				// Sort by Y coordinate
				swap = lines[j].Start.Y > lines[j+1].Start.Y
			}

			if swap {
				lines[j], lines[j+1] = lines[j+1], lines[j]
			}
		}
	}
}

// mergeAdjacent merges adjacent overlapping lines.
func (d *DefaultRulingLineDetector) mergeAdjacent(lines []*RulingLine, isHorizontal bool) []*RulingLine {
	if len(lines) == 0 {
		return lines
	}

	var result []*RulingLine
	current := lines[0]

	for i := 1; i < len(lines); i++ {
		next := lines[i]

		// Check if lines are adjacent (overlapping or close together)
		if d.areAdjacent(current, next, isHorizontal) {
			// Merge lines
			current = d.mergeTwo(current, next, isHorizontal)
		} else {
			// Not adjacent - save current and start new
			result = append(result, current)
			current = next
		}
	}

	// Add last line
	result = append(result, current)

	return result
}

// areAdjacent checks if two lines are adjacent on the same axis.
func (d *DefaultRulingLineDetector) areAdjacent(line1, line2 *RulingLine, isHorizontal bool) bool {
	if isHorizontal {
		// Check if horizontal lines overlap or are close on X axis
		gap := line2.Start.X - line1.End.X
		return gap <= d.tolerance*2 // Allow small gap
	}

	// Check if vertical lines overlap or are close on Y axis
	gap := line2.Start.Y - line1.End.Y
	return gap <= d.tolerance*2 // Allow small gap
}

// mergeTwo merges two lines into one.
func (d *DefaultRulingLineDetector) mergeTwo(line1, line2 *RulingLine, isHorizontal bool) *RulingLine {
	if isHorizontal {
		// Merge horizontal lines - extend on X axis
		start := extractor.NewPoint(
			math.Min(line1.Start.X, line2.Start.X),
			(line1.Start.Y+line2.Start.Y)/2, // Average Y
		)
		end := extractor.NewPoint(
			math.Max(line1.End.X, line2.End.X),
			(line1.End.Y+line2.End.Y)/2, // Average Y
		)
		return NewRulingLine(start, end)
	}

	// Merge vertical lines - extend on Y axis
	start := extractor.NewPoint(
		(line1.Start.X+line2.Start.X)/2, // Average X
		math.Min(line1.Start.Y, line2.Start.Y),
	)
	end := extractor.NewPoint(
		(line1.End.X+line2.End.X)/2, // Average X
		math.Max(line1.End.Y, line2.End.Y),
	)
	return NewRulingLine(start, end)
}

// FindIntersections finds intersection points between ruling lines.
//
// Returns a slice of unique intersection points.
func (d *DefaultRulingLineDetector) FindIntersections(lines []*RulingLine) []extractor.Point {
	var intersections []extractor.Point

	// Check each pair of lines
	for i := 0; i < len(lines)-1; i++ {
		for j := i + 1; j < len(lines); j++ {
			point := lines[i].Intersects(lines[j])
			if point != nil {
				intersections = append(intersections, *point)
			}
		}
	}

	// Remove duplicates
	intersections = d.uniquePoints(intersections)

	return intersections
}

// uniquePoints removes duplicate points within tolerance.
func (d *DefaultRulingLineDetector) uniquePoints(points []extractor.Point) []extractor.Point {
	if len(points) <= 1 {
		return points
	}

	var unique []extractor.Point

	for _, p := range points {
		isDuplicate := false
		for _, u := range unique {
			if math.Abs(p.X-u.X) <= d.tolerance && math.Abs(p.Y-u.Y) <= d.tolerance {
				isDuplicate = true
				break
			}
		}

		if !isDuplicate {
			unique = append(unique, p)
		}
	}

	return unique
}
